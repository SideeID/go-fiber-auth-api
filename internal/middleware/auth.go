package middleware

import (
	"context"
	"strings"
	"ujikom-backend/internal/models"
	"ujikom-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AuthMiddleware(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Authorization header is required")
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid token format. Use 'Bearer <token>'")
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Token is required")
		}

		claims, err := utils.ValidateJWT(token)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid or expired token")
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid token claims")
		}

		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid user ID in token")
		}

		collection := db.Collection("users")
		var user models.User
		err = collection.FindOne(context.Background(), bson.M{
			"_id":       userID,
			"is_active": true,
		}).Decode(&user)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not found or inactive")
			}
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
		}

		c.Locals("user", user)
		c.Locals("user_id", userID)

		return c.Next()
	}
}

// OptionalAuthMiddleware - middleware for optional authentication
func OptionalAuthMiddleware(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		
		if authHeader == "" {
			return c.Next()
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Next() 
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			return c.Next()
		}

		claims, err := utils.ValidateJWT(token)
		if err != nil {
			return c.Next() 
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return c.Next()
		}

		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			return c.Next()
		}

		collection := db.Collection("users")
		var user models.User
		err = collection.FindOne(context.Background(), bson.M{
			"_id":       userID,
			"is_active": true,
		}).Decode(&user)

		if err == nil {
			c.Locals("user", user)
			c.Locals("user_id", userID)
		}

		return c.Next()
	}
}