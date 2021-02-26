package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Exporter of webhook metrics.
// Build
type Exporter struct {
	bots []botExporter
}

// NewExporter creates exporter for provided bots
func NewExporter(ctx context.Context, tokens []string) *Exporter {
	exporter := &Exporter{
		bots: make([]botExporter, len(tokens)),
	}

	for i, token := range tokens {
		exporter.bots[i] = newBotExporter(token)
	}

	return exporter
}

// Describe implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- metricaUp
	ch <- metricaPendingUpdates
	ch <- metricaMaxConns
	ch <- metricaLastDeliveryError
}

// Collect implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	wg := &sync.WaitGroup{}

	ctx := context.Background()

	for i, bot := range e.bots {
		bot := bot
		i := i

		wg.Add(1)
		go func() {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(ctx, time.Second*30)
			defer cancel()

			started := time.Now()

			bot(ctx, ch)

			log.Printf("stats of bot #%d collect in %s", i, time.Since(started))
		}()
	}

	wg.Wait()
}
