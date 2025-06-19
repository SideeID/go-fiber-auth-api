package controllers

import (
	"context"
	"log"
	"time"
	"ujikom-backend/internal/models"
	"ujikom-backend/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	db        *mongo.Database
	validator *validator.Validate
}

func NewAuthController(db *mongo.Database) *AuthController {
	return &AuthController{
		db:        db,
		validator: validator.New(),
	}
}

func (ac *AuthController) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := ac.validator.Struct(req); err != nil {
		errors := utils.ValidatorErrors(err)
		return utils.ValidationErrorResponse(c, errors)
	}

	collection := ac.db.Collection("users")
	
	var existingUser models.User
	err := collection.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&existingUser)
	if err == nil {
		return utils.ErrorResponse(c, fiber.StatusConflict, "Email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to hash password")
	}

	now := time.Now().UTC()
	user := models.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Phone:     req.Phone,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create user")
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	
	token, err := utils.GenerateJWT(user.ID.Hex())
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate token")
	}

	response := models.LoginResponse{
		Token: token,
		User:  user.UserPublic(),
	}

	log.Printf("User registered successfully: %s", user.Email)
	return utils.SuccessResponse(c, "User registered successfully", response)
}

func (ac *AuthController) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := ac.validator.Struct(req); err != nil {
		errors := utils.ValidatorErrors(err)
		return utils.ValidationErrorResponse(c, errors)
	}

	collection := ac.db.Collection("users")
	
	var user models.User
	err := collection.FindOne(context.Background(), bson.M{
		"email":     req.Email,
		"is_active": true,
	}).Decode(&user)
	
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid email or password")
		}
		log.Printf("Error finding user: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid email or password")
	}

	_, err = collection.UpdateOne(
		context.Background(),
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{"updated_at": time.Now().UTC()}},
	)
	if err != nil {
		log.Printf("Error updating last login: %v", err)
	}

	token, err := utils.GenerateJWT(user.ID.Hex())
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate token")
	}

	response := models.LoginResponse{
		Token: token,
		User:  user.UserPublic(),
	}

	log.Printf("User logged in successfully: %s", user.Email)
	return utils.SuccessResponse(c, "Login successful", response)
}

func (ac *AuthController) Logout(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if ok {
		log.Printf("User logged out: %s", user.Email)
	}

	return utils.SuccessResponse(c, "Logout successful", nil)
}

func (ac *AuthController) RefreshToken(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not found")
	}

	token, err := utils.GenerateJWT(user.ID.Hex())
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate token")
	}

	response := fiber.Map{
		"token": token,
		"user":  user.UserPublic(),
	}

	return utils.SuccessResponse(c, "Token refreshed successfully", response)
}