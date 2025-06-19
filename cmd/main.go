package main

import (
	"log"
	"os"
	"strings"
	"ujikom-backend/internal/config"
	"ujikom-backend/internal/routes"
	"ujikom-backend/pkg/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

var (
	Version   = "1.0.0"
	BuildTime = "2025-06-19 04:06:32"
	GitCommit = "main"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	log.Printf("Ujikom Backend API")
	log.Printf("Version: %s", Version)
	log.Printf("Build Time: %s", BuildTime)
	log.Printf("Git Commit: %s", GitCommit)
	log.Printf("Created by: SideeID")
	log.Printf("Database: MongoDB Atlas")

	cfg := config.Load()

	if err := cfg.ValidateAtlasConnection(); err != nil {
		log.Fatal("Atlas configuration error:", err)
	}

	atlasInfo := extractAtlasInfo(cfg.MongoURI)
	log.Printf("Atlas Cluster: %s", atlasInfo)

	db, err := database.Connect(cfg.MongoURI, cfg.DBName)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB Atlas:", err)
	}

	app := fiber.New(fiber.Config{
		AppName:      "Ujikom API v" + Version + " (Atlas)",
		ServerHeader: "Ujikom-Backend-Atlas",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return c.Status(code).JSON(fiber.Map{
				"success":   false,
				"message":   err.Error(),
				"data":      nil,
				"version":   Version,
				"database":  "MongoDB Atlas",
				"timestamp": cfg.GetCurrentTime(),
			})
		},
	})

	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} - ${latency} | ${ip} | Atlas\n",
	}))
	app.Use(recover.New())

	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin,Content-Type,Accept,Authorization,X-Requested-With")
		c.Set("Access-Control-Allow-Credentials", "true")
		c.Set("X-API-Version", Version)
		c.Set("X-Build-Time", BuildTime)
		c.Set("X-Database", "MongoDB-Atlas")
		c.Set("X-Created-By", "SideeID")
		
		if c.Method() == "OPTIONS" {
			return c.SendStatus(200)
		}
		
		return c.Next()
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Ujikom API is running with MongoDB Atlas!",
			"data": fiber.Map{
				"version":     Version,
				"build_time":  BuildTime,
				"git_commit":  GitCommit,
				"server":      "Go Fiber",
				"database":    "MongoDB Atlas",
				"cluster":     atlasInfo,
				"public_api":  true,
				"cors":        "disabled",
				"created_by":  "SideeID",
				"timestamp":   cfg.GetCurrentTime(),
			},
		})
	})

	routes.Setup(app, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Public API: http://localhost:%s", port)
	log.Printf("Docs: http://localhost:%s/api/v1/docs", port)
	log.Printf("Health: http://localhost:%s/api/v1/health", port)
	log.Printf("Database: MongoDB Atlas Connected")
	
	log.Fatal(app.Listen(":" + port))
}

func extractAtlasInfo(uri string) string {
	if strings.Contains(uri, "@") && strings.Contains(uri, ".mongodb.net") {
		parts := strings.Split(uri, "@")
		if len(parts) > 1 {
			clusterPart := strings.Split(parts[1], "/")[0]
			return strings.Split(clusterPart, ".mongodb.net")[0]
		}
	}
	return "Unknown Atlas Cluster"
}