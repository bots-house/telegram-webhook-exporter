package main

import "github.com/prometheus/client_golang/prometheus"

const (
	namespace = "telegram"
	subsystem = "bot_api"
)

var (
	metricaUp = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "up"),
		"true, if api is up",
		[]string{"bot"}, nil,
	)

	metricaPendingUpdates = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "pending_updates"),
		"number of updates awaiting delivery",
		[]string{"bot"}, nil,
	)

	metricaMaxConns = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "max_conns"),
		"Maximum allowed number of simultaneous HTTPS connections to the webhook for update delivery",
		[]string{"bot"}, nil,
	)

	metricaLastDeliveryError = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "last_delivery_error"),
		"time and message for the most recent error that happened when trying to deliver an update via webhook",
		[]string{"bot", "msg"}, nil,
	)

	metricaResponseTime = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "response_time"),
		"gauge of bot api response time",
		[]string{"bot"}, nil,
	)
)
