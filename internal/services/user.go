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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserService struct {
	db  *mongo.Database
	ctx context.Context
}

type UserServiceInterface interface {
	GetUserByID(userID string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUserProfile(userID string, req *models.UpdateProfileRequest) (*models.User, error)
	ChangeUserPassword(userID string, req *models.ChangePasswordRequest) error
	DeleteUser(userID string) error
	GetAllUsers(limit, offset int) ([]models.User, int64, error)
	SearchUsers(query string, limit, offset int) ([]models.User, int64, error)
	GetUserProfile(userID string) (*models.User, error)
	UpdateUserAvatar(userID, avatarURL string) error
	GetUserActivity(userID string) (map[string]interface{}, error)
}

func NewUserService(db *mongo.Database) UserServiceInterface {
	return &UserService{
		db:  db,
		ctx: context.Background(),
	}
}

func (s *UserService) GetUserByID(userID string) (*models.User, error) {
	collection := s.db.Collection("users")

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	var user models.User
	err = collection.FindOne(s.ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		log.Printf("Error finding user by ID: %v", err)
		return nil, errors.New("database error")
	}

	return &user, nil
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	collection := s.db.Collection("users")

	var user models.User
	err := collection.FindOne(s.ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		log.Printf("Error finding user by email: %v", err)
		return nil, errors.New("database error")
	}

	return &user, nil
}

func (s *UserService) GetUserProfile(userID string) (*models.User, error) {
	user, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("user account is inactive")
	}

	return user, nil
}

func (s *UserService) UpdateUserProfile(userID string, req *models.UpdateProfileRequest) (*models.User, error) {
	collection := s.db.Collection("users")

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	updateDoc := bson.M{
		"updated_at": time.Now().UTC(),
	}

	if req.Name != "" {
		updateDoc["name"] = utils.SanitizeInput(req.Name)
	}
	if req.Phone != "" {
		updateDoc["phone"] = utils.SanitizeInput(req.Phone)
	}
	if req.Avatar != "" {
		updateDoc["avatar"] = utils.SanitizeInput(req.Avatar)
	}

	filter := bson.M{"_id": objectID, "is_active": true}
	update := bson.M{"$set": updateDoc}

	result, err := collection.UpdateOne(s.ctx, filter, update)
	if err != nil {
		log.Printf("Error updating user profile: %v", err)
		return nil, errors.New("failed to update profile")
	}

	if result.MatchedCount == 0 {
		return nil, errors.New("user not found or inactive")
	}

	updatedUser, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	log.Printf("User profile updated: %s (ID: %s)", updatedUser.Email, userID)
	return updatedUser, nil
}

func (s *UserService) ChangeUserPassword(userID string, req *models.ChangePasswordRequest) error {
	collection := s.db.Collection("users")

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	var user models.User
	err = collection.FindOne(s.ctx, bson.M{"_id": objectID, "is_active": true}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("user not found or inactive")
		}
		return errors.New("database error")
	}

	err = utils.ComparePasswords(user.Password, req.CurrentPassword)
	if err != nil {
		return errors.New("current password is incorrect")
	}

	if passwordErrors := utils.ValidatePasswordStrength(req.NewPassword); len(passwordErrors) > 0 {
		return errors.New("new password does not meet security requirements")
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		log.Printf("Error hashing new password: %v", err)
		return errors.New("failed to process new password")
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"password":   hashedPassword,
			"updated_at": time.Now().UTC(),
		},
	}

	_, err = collection.UpdateOne(s.ctx, filter, update)
	if err != nil {
		log.Printf("Error updating password: %v", err)
		return errors.New("failed to update password")
	}

	log.Printf("Password changed for user: %s (ID: %s)", user.Email, userID)
	return nil
}

func (s *UserService) UpdateUserAvatar(userID, avatarURL string) error {
	collection := s.db.Collection("users")

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	filter := bson.M{"_id": objectID, "is_active": true}
	update := bson.M{
		"$set": bson.M{
			"avatar":     utils.SanitizeInput(avatarURL),
			"updated_at": time.Now().UTC(),
		},
	}

	result, err := collection.UpdateOne(s.ctx, filter, update)
	if err != nil {
		log.Printf("Error updating user avatar: %v", err)
		return errors.New("failed to update avatar")
	}

	if result.MatchedCount == 0 {
		return errors.New("user not found or inactive")
	}

	log.Printf("Avatar updated for user: %s", userID)
	return nil
}

func (s *UserService) DeleteUser(userID string) error {
	collection := s.db.Collection("users")

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	// Soft delete
	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"is_active":  false,
			"updated_at": time.Now().UTC(),
		},
	}

	result, err := collection.UpdateOne(s.ctx, filter, update)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		return errors.New("failed to delete user")
	}

	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}

	log.Printf("User deleted (soft): %s", userID)
	return nil
}

func (s *UserService) GetAllUsers(limit, offset int) ([]models.User, int64, error) {
	collection := s.db.Collection("users")

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 
	}

	total, err := collection.CountDocuments(s.ctx, bson.M{"is_active": true})
	if err != nil {
		return nil, 0, errors.New("failed to count users")
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))
	findOptions.SetSort(bson.M{"created_at": -1}) 

	cursor, err := collection.Find(s.ctx, bson.M{"is_active": true}, findOptions)
	if err != nil {
		log.Printf("Error finding users: %v", err)
		return nil, 0, errors.New("failed to retrieve users")
	}
	defer cursor.Close(s.ctx)

	var users []models.User
	if err = cursor.All(s.ctx, &users); err != nil {
		log.Printf("Error decoding users: %v", err)
		return nil, 0, errors.New("failed to decode users")
	}

	return users, total, nil
}

func (s *UserService) SearchUsers(query string, limit, offset int) ([]models.User, int64, error) {
	collection := s.db.Collection("users")

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	searchFilter := bson.M{
		"is_active": true,
		"$or": []bson.M{
			{"name": bson.M{"$regex": query, "$options": "i"}},
			{"email": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	total, err := collection.CountDocuments(s.ctx, searchFilter)
	if err != nil {
		return nil, 0, errors.New("failed to count search results")
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))
	findOptions.SetSort(bson.M{"created_at": -1})

	cursor, err := collection.Find(s.ctx, searchFilter, findOptions)
	if err != nil {
		log.Printf("Error searching users: %v", err)
		return nil, 0, errors.New("failed to search users")
	}
	defer cursor.Close(s.ctx)

	var users []models.User
	if err = cursor.All(s.ctx, &users); err != nil {
		log.Printf("Error decoding search results: %v", err)
		return nil, 0, errors.New("failed to decode search results")
	}

	return users, total, nil
}

func (s *UserService) GetUserActivity(userID string) (map[string]interface{}, error) {
	user, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	accountAge := time.Since(user.CreatedAt)
	
	activity := map[string]interface{}{
		"user_id":      user.ID.Hex(),
		"email":        user.Email,
		"name":         user.Name,
		"is_active":    user.IsActive,
		"account_age":  accountAge.Hours() / 24, 
		"created_at":   user.CreatedAt.Format("2006-01-02 15:04:05"),
		"updated_at":   user.UpdatedAt.Format("2006-01-02 15:04:05"),
		"last_seen":    user.UpdatedAt.Format("2006-01-02 15:04:05"),
		"profile_complete": s.isProfileComplete(user),
	}

	return activity, nil
}

func (s *UserService) isProfileComplete(user *models.User) bool {
	return user.Name != "" && user.Email != "" && user.Phone != ""
}