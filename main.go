package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/kelseyhightower/envconfig"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/xerrors"
)

type Config struct {
	Addr   string   `envconfig:"ADDR" default:":8000"`
	Tokens []string `envconfig:"TOKENS" required:"true"`
}

func main() {
	ctx := context.Background()

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	if err := run(ctx); err != nil {
		log.Printf("run error: %v", err)
		defer os.Exit(1)
	}
}

func run(ctx context.Context) error {
	var cfg Config

	err := envconfig.Process("exporter", &cfg)
	if err != nil {
		return xerrors.Errorf("parse config: %w", err)
	}

	registry := prometheus.NewRegistry()

	log.Printf("create exporter with %d bots", len(cfg.Tokens))

	exporter := NewExporter(ctx, cfg.Tokens)

	registry.MustRegister(exporter)

	log.Printf("start server at %s", cfg.Addr)

	return runServer(ctx, cfg.Addr, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
}
