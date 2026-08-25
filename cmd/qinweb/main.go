package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"qin-culture-site/internal/catalog"
	"qin-culture-site/internal/config"
	"qin-culture-site/internal/service"
	"qin-culture-site/internal/store"
	"qin-culture-site/internal/web"
)

func main() {
	cfg := config.Load()
	if err := run(cfg); err != nil {
		log.Fatal(err)
	}
}

func run(cfg config.Config) error {
	if err := cfg.Validate(); err != nil {
		return err
	}
	data, err := store.Open(cfg.Database)
	if err != nil {
		return err
	}
	defer data.Close()
	svc := service.New(catalog.New(), data)
	server, err := web.NewServer(svc, cfg.MaxBody)
	if err != nil {
		return err
	}
	httpServer := &http.Server{Addr: cfg.Address, Handler: server.Handler(), ReadHeaderTimeout: 5 * time.Second}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-stop
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = httpServer.Shutdown(ctx)
	}()
	log.Printf("%s", cfg.Summary())
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
