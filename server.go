package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

func runServer(ctx context.Context, addr string, handler http.Handler) error {
	server := http.Server{
		Addr:    ":8000",
		Handler: handler,
	}

	go func() {
		<-ctx.Done()
		const timeout = time.Second

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		log.Printf("shutdown server timeout=%s", timeout)
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("shutdown fail: %+v", err)
		}
	}()

	err := server.ListenAndServe()

	if err == http.ErrServerClosed {
		return nil
	}

	return err
}
