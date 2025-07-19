package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	
	"github.com/flack/chaos-engineering-as-a-platform/pkg/api/handlers"
	"github.com/flack/chaos-engineering-as-a-platform/pkg/config"
	"github.com/flack/chaos-engineering-as-a-platform/pkg/monitoring"
	"github.com/flack/chaos-engineering-as-a-platform/pkg/storage"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := storage.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	metrics := monitoring.NewMetrics()
	
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Serve static files
	r.StaticFile("/", "./web/dashboard/api-index.html")

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy", 
			"timestamp": time.Now(),
		})
	})

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	v1 := r.Group("/api/v1")
	{
		experiments := handlers.NewExperimentHandler(db, metrics)
		v1.POST("/experiments", experiments.CreateExperiment)
		v1.GET("/experiments", experiments.ListExperiments)
		v1.GET("/experiments/:id", experiments.GetExperiment)
		v1.POST("/experiments/:id/execute", experiments.ExecuteExperiment)
		v1.DELETE("/experiments/:id", experiments.DeleteExperiment)

		targets := handlers.NewTargetHandler(db)
		v1.GET("/targets", targets.ListTargets)
		v1.POST("/targets", targets.CreateTarget)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", cfg.Port),
		Handler: r,
	}

	go func() {
		log.Printf("Starting server on port %d", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("Shutting down server...")
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	
	log.Println("Server exiting")
}