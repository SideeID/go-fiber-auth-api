package routes

import (
	"ujikom-backend/internal/controllers"
	"ujikom-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func Setup(app *fiber.App, db *mongo.Database) {
	authController := controllers.NewAuthController(db)
	userController := controllers.NewUserController(db)

	api := app.Group("/api/v1")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Ujikom API is healthy",
			"service": "ujikom-backend",
			"version": "1.0.0",
			"timestamp": "2025-06-17 16:18:07",
		})
	})

	auth := api.Group("/auth")
	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
	auth.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Auth endpoint is working!",
			"public":  true,
			"cors":    "disabled",
			"usage":   "free for everyone",
		})
	})

	protected := api.Group("/user")
	protected.Use(middleware.AuthMiddleware(db))
	
	protected.Get("/profile", userController.GetProfile)
	protected.Put("/profile", userController.UpdateProfile)
	protected.Post("/change-password", userController.ChangePassword)
	protected.Post("/deactivate", userController.DeactivateAccount)
	
	protected.Post("/logout", authController.Logout)
	protected.Post("/refresh-token", authController.RefreshToken)

	testing := api.Group("/testing")
	testing.Use(middleware.OptionalAuthMiddleware(db))
	testing.Get("/users", userController.GetAllUsers)

	api.Get("/docs", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Ujikom API Documentation",
			"endpoints": fiber.Map{
				"public": []string{
					"GET /api/v1/health",
					"POST /api/v1/auth/register",
					"POST /api/v1/auth/login",
					"GET /api/v1/auth/test",
				},
				"protected": []string{
					"GET /api/v1/user/profile",
					"PUT /api/v1/user/profile",
					"POST /api/v1/user/change-password",
					"POST /api/v1/user/deactivate",
					"POST /api/v1/user/logout",
					"POST /api/v1/user/refresh-token",
				},
				"testing": []string{
					"GET /api/v1/testing/users",
				},
			},
			"authentication": "Bearer token required for protected routes",
			"cors": "disabled - public API",
			"created_by": "SideeID",
			"created_at": "2025-06-17 16:18:07",
		})
	})
}