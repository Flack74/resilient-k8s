package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	
	"github.com/flack/chaos-engineering-as-a-platform/pkg/api/handlers"
	"github.com/flack/chaos-engineering-as-a-platform/pkg/api/middleware"
	"github.com/flack/chaos-engineering-as-a-platform/pkg/config"
	"github.com/flack/chaos-engineering-as-a-platform/pkg/k8s/operator"
	"github.com/flack/chaos-engineering-as-a-platform/pkg/monitoring"
	"github.com/flack/chaos-engineering-as-a-platform/pkg/storage"
)

// SetupRouter sets up the API routes
func SetupRouter(db *storage.Database, metrics *monitoring.Metrics, chaosOperator *operator.ChaosOperator, cfg *config.Config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.MetricsMiddleware(metrics))

	// Public endpoints
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API endpoints (authentication disabled for development)
	api := r.Group("/api")
	// Uncomment the following line to enable authentication
	// api.Use(middleware.AuthMiddleware(getAPIKey(cfg)))
	
	v1 := api.Group("/v1")
	{
		// Experiment endpoints
		experimentHandler := handlers.NewExperimentHandler(db, metrics)
		experimentHandler.SetOperator(chaosOperator)
		
		v1.POST("/experiments", experimentHandler.CreateExperiment)
		v1.GET("/experiments", experimentHandler.ListExperiments)
		v1.GET("/experiments/:id", experimentHandler.GetExperiment)
		v1.POST("/experiments/:id/execute", experimentHandler.ExecuteExperiment)
		v1.DELETE("/experiments/:id", experimentHandler.DeleteExperiment)

		// Target endpoints
		targetHandler := handlers.NewTargetHandler(db)
		v1.GET("/targets", targetHandler.ListTargets)
		v1.POST("/targets", targetHandler.CreateTarget)
	}

	return r
}

// getAPIKey gets the API key from the configuration
func getAPIKey(cfg *config.Config) string {
	// In a real application, this would be a secure API key
	// For development, we'll use a simple default
	return "chaos-platform-api-key"
}