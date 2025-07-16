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
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	db        *mongo.Database
	validator *validator.Validate
}

func NewUserController(db *mongo.Database) *UserController {
	return &UserController{
		db:        db,
		validator: validator.New(),
	}
}

func (uc *UserController) GetProfile(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not found")
	}

	return utils.SuccessResponse(c, "Profile retrieved successfully", user.UserPublic())
}

func (uc *UserController) UpdateProfile(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not found")
	}

	var req models.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := uc.validator.Struct(req); err != nil {
		errors := utils.ValidatorErrors(err)
		return utils.ValidationErrorResponse(c, errors)
	}

	updateDoc := bson.M{
		"updated_at": time.Now().UTC(),
	}

	if req.Name != "" {
		updateDoc["name"] = req.Name
	}
	if req.Kelas != "" {
		updateDoc["kelas"] = req.Kelas
	}
	if req.Jurusan != "" {
		updateDoc["jurusan"] = req.Jurusan
	}
	if req.Phone != "" {
		updateDoc["phone"] = req.Phone
	}
	if req.Avatar != "" {
		updateDoc["avatar"] = req.Avatar
	}

	collection := uc.db.Collection("users")
	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": user.ID},
		bson.M{"$set": updateDoc},
	)

	if err != nil {
		log.Printf("Error updating user profile: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update profile")
	}

	var updatedUser models.User
	err = collection.FindOne(context.Background(), bson.M{"_id": user.ID}).Decode(&updatedUser)
	if err != nil {
		log.Printf("Error fetching updated user: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch updated profile")
	}

	log.Printf("User profile updated: %s", updatedUser.Email)
	return utils.SuccessResponse(c, "Profile updated successfully", updatedUser.UserPublic())
}

func (uc *UserController) ChangePassword(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not found")
	}

	var req models.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := uc.validator.Struct(req); err != nil {
		errors := utils.ValidatorErrors(err)
		return utils.ValidationErrorResponse(c, errors)
	}

	collection := uc.db.Collection("users")
	var currentUser models.User
	err := collection.FindOne(context.Background(), bson.M{"_id": user.ID}).Decode(&currentUser)
	if err != nil {
		log.Printf("Error finding user: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}

	err = bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(req.CurrentPassword))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Current password is incorrect")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to hash password")
	}

	_, err = collection.UpdateOne(
		context.Background(),
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{
			"password":   string(hashedPassword),
			"updated_at": time.Now().UTC(),
		}},
	)

	if err != nil {
		log.Printf("Error updating password: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update password")
	}

	log.Printf("Password changed for user: %s", user.Email)
	return utils.SuccessResponse(c, "Password changed successfully", nil)
}

func (uc *UserController) DeactivateAccount(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not found")
	}

	collection := uc.db.Collection("users")
	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{
			"is_active":  false,
			"updated_at": time.Now().UTC(),
		}},
	)

	if err != nil {
		log.Printf("Error deactivating account: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to deactivate account")
	}

	log.Printf("Account deactivated: %s", user.Email)
	return utils.SuccessResponse(c, "Account deactivated successfully", nil)
}

func (uc *UserController) GetAllUsers(c *fiber.Ctx) error {
	
	collection := uc.db.Collection("users")
	
	cursor, err := collection.Find(context.Background(), bson.M{"is_active": true})
	if err != nil {
		log.Printf("Error finding users: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}
	defer cursor.Close(context.Background())

	var users []models.User
	if err = cursor.All(context.Background(), &users); err != nil {
		log.Printf("Error decoding users: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}

	var publicUsers []models.User
	for _, user := range users {
		publicUsers = append(publicUsers, user.UserPublic())
	}

	return utils.SuccessResponse(c, "Users retrieved successfully", fiber.Map{
		"users": publicUsers,
		"count": len(publicUsers),
	})
}