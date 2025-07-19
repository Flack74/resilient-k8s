package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/flack/chaos-engineering-as-a-platform/pkg/config"
	"github.com/flack/chaos-engineering-as-a-platform/pkg/k8s/operator"
	"github.com/flack/chaos-engineering-as-a-platform/pkg/monitoring"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	metrics := monitoring.NewMetrics()

	chaosOperator, err := operator.NewChaosOperator(cfg, metrics)
	if err != nil {
		log.Fatalf("Failed to create chaos operator: %v", err)
	}

	go func() {
		log.Println("Starting chaos operator...")
		if err := chaosOperator.Start(); err != nil {
			log.Fatalf("Failed to start chaos operator: %v", err)
		}
	}()

	// Set up signal handling for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down chaos operator...")

	if err := chaosOperator.Stop(); err != nil {
		log.Fatalf("Error during operator shutdown: %v", err)
	}

	log.Println("Chaos operator exited")
}
