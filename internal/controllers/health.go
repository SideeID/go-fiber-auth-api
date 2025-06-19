package controllers

import (
	"ujikom-backend/internal/utils"
	"ujikom-backend/pkg/database"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type HealthController struct {
	db *mongo.Database
}

func NewHealthController(db *mongo.Database) *HealthController {
	return &HealthController{db: db}
}

func (h *HealthController) GetHealth(c *fiber.Ctx) error {
	atlasStatus := "connected"
	var atlasError string
	
	if err := database.HealthCheck(h.db); err != nil {
		atlasStatus = "disconnected"
		atlasError = err.Error()
	}

	connectionStats, _ := database.GetConnectionStats(h.db)

	healthData := fiber.Map{
		"service":    "ujikom-backend",
		"version":    "1.0.0",
		"status":     "healthy",
		"database": fiber.Map{
			"type":             "MongoDB Atlas",
			"status":           atlasStatus,
			"error":            atlasError,
			"connection_stats": connectionStats,
		},
		"uptime":     "running",
		"created_by": "SideeID",
		"timestamp":  "2025-06-19 04:06:32",
	}

	if atlasStatus == "disconnected" {
		return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "Service unhealthy - Atlas disconnected")
	}

	return utils.SuccessResponse(c, "Service is healthy", healthData)
}