package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"-" bson:"password"`
	Phone     string             `json:"phone,omitempty" bson:"phone,omitempty"`
	Avatar    string             `json:"avatar,omitempty" bson:"avatar,omitempty"`
	IsActive  bool               `json:"is_active" bson:"is_active"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type LoginRequest struct {
	Email 		string `json:"email" validate:"required,email"`
	Password	string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	Name 			string `json:"name" validate:"required,min=2,max=100"`
	Email 		string `json:"email" validate:"required,email"`
	Password	string `json:"password" validate:"required,min=6"`
	Phone   	string `json:"phone,omitempty" validate:"omitempty,min=10,max=15"`
}

type UpdateProfileRequest struct {
	Name 		string `json:"name" validate:"required,min=2,max=100"`
	Phone		string `json:"phone,omitempty" validate:"omitempty,min=10,max=15"`
	Avatar	string `json:"avatar,omitempty" validate:"omitempty,url"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required,min=6"`
	NewPassword     string `json:"new_password" validate:"required,min=6,max=50"`
}

func (u *User) UserPublic() User {
	return User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		Avatar:    u.Avatar,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}