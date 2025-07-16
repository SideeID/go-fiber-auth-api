package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Attendance struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	Date      time.Time          `json:"date" bson:"date"`
	CheckIn   *time.Time         `json:"check_in,omitempty" bson:"check_in,omitempty"`
	CheckOut  *time.Time         `json:"check_out,omitempty" bson:"check_out,omitempty"`
	Status    string             `json:"status" bson:"status"` // present, late, absent
	Location  Location           `json:"location" bson:"location"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type Location struct {
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
	Address   string  `json:"address,omitempty" bson:"address,omitempty"`
}

type AttendanceRequest struct {
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
	Address   string  `json:"address,omitempty"`
}

func (ar AttendanceRequest) ToLocation() Location {
	return Location(ar)
}

type AttendanceResponse struct {
	ID        primitive.ObjectID `json:"id"`
	UserID    primitive.ObjectID `json:"user_id"`
	Date      time.Time          `json:"date"`
	CheckIn   *time.Time         `json:"check_in"`
	CheckOut  *time.Time         `json:"check_out"`
	Status    string             `json:"status"`
	Location  Location           `json:"location"`
	User      UserPublic         `json:"user"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type UserPublic struct {
	ID      primitive.ObjectID `json:"id"`
	NIS     string             `json:"nis"`
	Name    string             `json:"name"`
	Kelas   string             `json:"kelas"`
	Jurusan string             `json:"jurusan"`
	Email   string             `json:"email"`
	Phone   string             `json:"phone"`
}

type AttendanceStats struct {
	TotalPresent int     `json:"total_present"`
	TotalLate    int     `json:"total_late"`
	TotalAbsent  int     `json:"total_absent"`
	Percentage   float64 `json:"percentage"`
}

type AttendanceHistory struct {
	Date     time.Time `json:"date"`
	CheckIn  *time.Time `json:"check_in"`
	CheckOut *time.Time `json:"check_out"`
	Status   string    `json:"status"`
	Location Location  `json:"location"`
}
