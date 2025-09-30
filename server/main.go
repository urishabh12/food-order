package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"order-food/internal/api"
	"order-food/internal/config"
	"order-food/internal/coupons"
)

func main() {
	cfg := config.FromEnv()

	couponStore, err := coupons.Load(cfg.CouponsDir)
	if err != nil {
		log.Fatalf("load coupons: %v", err)
	}

	h := &api.Handlers{Coupons: couponStore}
	r := api.NewRouter(cfg, h)

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// start
	go func() {
		log.Printf("listening on %s (basePath %s)", cfg.Addr, cfg.BasePath)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Println("server stopped")
}
