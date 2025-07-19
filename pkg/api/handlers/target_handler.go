package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	
	"github.com/flack/chaos-engineering-as-a-platform/pkg/storage"
)

// TargetHandler handles target-related API requests
type TargetHandler struct {
	db *storage.Database
}

// NewTargetHandler creates a new target handler
func NewTargetHandler(db *storage.Database) *TargetHandler {
	return &TargetHandler{
		db: db,
	}
}

// CreateTargetRequest represents a request to create a new target
type CreateTargetRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Type        string `json:"type" binding:"required"`
	Namespace   string `json:"namespace" binding:"required"`
	Selector    string `json:"selector" binding:"required"`
}

// CreateTarget handles the creation of a new target
func (h *TargetHandler) CreateTarget(c *gin.Context) {
	var req CreateTargetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a new target
	now := time.Now()
	target := &storage.Target{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Type:        storage.TargetType(req.Type),
		Namespace:   req.Namespace,
		Selector:    req.Selector,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Save the target to the database
	if err := h.db.CreateTarget(target); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, target)
}

// ListTargets handles listing all targets
func (h *TargetHandler) ListTargets(c *gin.Context) {
	targets, err := h.db.ListTargets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, targets)
}