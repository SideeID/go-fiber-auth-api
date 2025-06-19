package services

import (
	"context"
	"errors"
	"log"
	"time"
	"ujikom-backend/internal/models"
	"ujikom-backend/internal/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthService struct {
	db   *mongo.Database
	ctx  context.Context
}

type AuthServiceInterface interface {
	RegisterUser(req *models.RegisterRequest) (*models.User, error)
	LoginUser(req *models.LoginRequest) (*models.User, error)
	ValidateUser(userID string) (*models.User, error)
	GenerateTokens(userID string) (string, error)
	RefreshUserToken(userID string) (string, error)
	DeactivateUser(userID string) error
	CheckEmailExists(email string) (bool, error)
	UpdateLastLogin(userID string) error
}

func NewAuthService(db *mongo.Database) AuthServiceInterface {
	return &AuthService{
		db:  db,
		ctx: context.Background(),
	}
}

func (s *AuthService) RegisterUser(req *models.RegisterRequest) (*models.User, error) {
	collection := s.db.Collection("users")

	exists, err := s.CheckEmailExists(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return nil, errors.New("failed to process password")
	}

	now := time.Now().UTC()
	user := models.User{
		Name:      utils.SanitizeInput(req.Name),
		Email:     utils.SanitizeInput(req.Email),
		Password:  hashedPassword,
		Phone:     utils.SanitizeInput(req.Phone),
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result, err := collection.InsertOne(s.ctx, user)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		return nil, errors.New("failed to create user")
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	
	log.Printf("✅ New user registered: %s (ID: %s)", user.Email, user.ID.Hex())
	return &user, nil
}

func (s *AuthService) LoginUser(req *models.LoginRequest) (*models.User, error) {
	collection := s.db.Collection("users")

	var user models.User
	filter := bson.M{
		"email":     utils.SanitizeInput(req.Email),
		"is_active": true,
	}

	err := collection.FindOne(s.ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("invalid email or password")
		}
		log.Printf("Error finding user: %v", err)
		return nil, errors.New("database error")
	}

	err = utils.ComparePasswords(user.Password, req.Password)
	if err != nil {
		log.Printf("Invalid password attempt for user: %s", user.Email)
		return nil, errors.New("invalid email or password")
	}

	err = s.UpdateLastLogin(user.ID.Hex())
	if err != nil {
		log.Printf("Warning: failed to update last login for user %s: %v", user.ID.Hex(), err)
	}

	log.Printf("✅ User logged in: %s (ID: %s)", user.Email, user.ID.Hex())
	return &user, nil
}

func (s *AuthService) ValidateUser(userID string) (*models.User, error) {
	collection := s.db.Collection("users")

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	var user models.User
	filter := bson.M{
		"_id":       objectID,
		"is_active": true,
	}

	err = collection.FindOne(s.ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found or inactive")
		}
		log.Printf("Error validating user: %v", err)
		return nil, errors.New("database error")
	}

	return &user, nil
}

func (s *AuthService) GenerateTokens(userID string) (string, error) {
	token, err := utils.GenerateJWT(userID)
	if err != nil {
		log.Printf("Error generating JWT for user %s: %v", userID, err)
		return "", errors.New("failed to generate authentication token")
	}

	return token, nil
}

func (s *AuthService) RefreshUserToken(userID string) (string, error) {
	user, err := s.ValidateUser(userID)
	if err != nil {
		return "", err
	}

	token, err := s.GenerateTokens(user.ID.Hex())
	if err != nil {
		return "", err
	}

	log.Printf("✅ Token refreshed for user: %s (ID: %s)", user.Email, user.ID.Hex())
	return token, nil
}

func (s *AuthService) DeactivateUser(userID string) error {
	collection := s.db.Collection("users")

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"is_active":  false,
			"updated_at": time.Now().UTC(),
		},
	}

	result, err := collection.UpdateOne(s.ctx, filter, update)
	if err != nil {
		log.Printf("Error deactivating user: %v", err)
		return errors.New("failed to deactivate user")
	}

	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}

	log.Printf("✅ User deactivated: %s", userID)
	return nil
}

func (s *AuthService) CheckEmailExists(email string) (bool, error) {
	collection := s.db.Collection("users")

	count, err := collection.CountDocuments(s.ctx, bson.M{"email": email})
	if err != nil {
		log.Printf("Error checking email existence: %v", err)
		return false, errors.New("database error")
	}

	return count > 0, nil
}

func (s *AuthService) UpdateLastLogin(userID string) error {
	collection := s.db.Collection("users")

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"updated_at": time.Now().UTC(),
		},
	}

	_, err = collection.UpdateOne(s.ctx, filter, update)
	if err != nil {
		log.Printf("Error updating last login: %v", err)
		return errors.New("failed to update last login")
	}

	return nil
}

func (s *AuthService) GetUserStats() (map[string]interface{}, error) {
	collection := s.db.Collection("users")

	totalUsers, err := collection.CountDocuments(s.ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	activeUsers, err := collection.CountDocuments(s.ctx, bson.M{"is_active": true})
	if err != nil {
		return nil, err
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	todayUsers, err := collection.CountDocuments(s.ctx, bson.M{
		"created_at": bson.M{"$gte": today},
	})
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_users":    totalUsers,
		"active_users":   activeUsers,
		"inactive_users": totalUsers - activeUsers,
		"today_users":    todayUsers,
		"timestamp":      time.Now().UTC().Format("2006-01-02 15:04:05"),
	}

	return stats, nil
}