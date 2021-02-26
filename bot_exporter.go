package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type botExporter func(ctx context.Context, ch chan<- prometheus.Metric)

func newBotExporter(token string) botExporter {
	var user struct {
		Username string
	}

	splitedToken := strings.Split(token, ":")

	id := "?"
	if len(splitedToken) > 1 {
		id = splitedToken[0]
	}

	isUsernameNotSet := true

	wrapCall := func(ctx context.Context, ch chan<- prometheus.Metric, fn func(ctx context.Context) error) error {
		if err := fn(ctx); err != nil {

			if !isUsernameNotSet {
				log.Printf("report down")
				ch <- prometheus.MustNewConstMetric(
					metricaUp, prometheus.GaugeValue,
					0,
					id,
				)
			}

			return err
		}

		if !isUsernameNotSet {
			ch <- prometheus.MustNewConstMetric(
				metricaUp, prometheus.GaugeValue,
				1,
				id,
			)
		}

		return nil
	}

	return func(ctx context.Context, ch chan<- prometheus.Metric) {
		if isUsernameNotSet {
			if err := wrapCall(ctx, ch, func(ctx context.Context) error {
				return call(ctx, token, "getMe", &user)
			}); err != nil {
				log.Printf("call getMe, for get username: %+v", err)
				return
			}
			id = user.Username
			isUsernameNotSet = false
		}

		var webhook struct {
			PendingUpdateCount int    `json:"pending_update_count"`
			LastErrorDate      int64  `json:"last_error_date"`
			LastErrorMessage   string `json:"last_error_message"`
			MaxConnections     int    `json:"max_connections"`
		}

		if err := wrapCall(ctx, ch, func(ctx context.Context) error {
			started := time.Now()

			if err := call(ctx, token, "getWebhookInfo", &webhook); err != nil {
				return err
			}

			ch <- prometheus.MustNewConstMetric(
				metricaResponseTime,
				prometheus.GaugeValue,
				float64(time.Since(started).Milliseconds()),
				id,
			)

			ch <- prometheus.MustNewConstMetric(
				metricaPendingUpdates,
				prometheus.GaugeValue,
				float64(webhook.PendingUpdateCount),
				id,
			)

			ch <- prometheus.MustNewConstMetric(
				metricaMaxConns,
				prometheus.GaugeValue,
				float64(webhook.MaxConnections),
				id,
			)

			ch <- prometheus.MustNewConstMetric(
				metricaLastDeliveryError,
				prometheus.GaugeValue,
				float64(webhook.LastErrorDate),
				id, webhook.LastErrorMessage,
			)

			return nil
		}); err != nil {
			log.Printf("check webhook: %+v", err)
			return
		}

	}
}
