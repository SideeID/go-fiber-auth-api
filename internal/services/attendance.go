package services

import (
	"context"
	"errors"
	"log"
	"time"
	"ujikom-backend/internal/config"
	"ujikom-backend/internal/models"
	"ujikom-backend/internal/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AttendanceService struct {
	db     *mongo.Database
	ctx    context.Context
	config *config.Config
}

type AttendanceServiceInterface interface {
	CheckIn(userID string, req *models.AttendanceRequest) (*models.Attendance, error)
	CheckOut(userID string, req *models.AttendanceRequest) (*models.Attendance, error)
	GetTodayAttendance(userID string) (*models.Attendance, error)
	GetAttendanceHistory(userID string, limit, offset int) ([]models.Attendance, int64, error)
	GetAttendanceStats(userID string) (*models.AttendanceStats, error)
	GetAttendanceByDate(userID string, date time.Time) (*models.Attendance, error)
	IsValidLocation(lat, lng float64) bool
	DetermineStatus(checkInTime time.Time) string
}

func NewAttendanceService(db *mongo.Database, cfg *config.Config) AttendanceServiceInterface {
	return &AttendanceService{
		db:     db,
		ctx:    context.Background(),
		config: cfg,
	}
}

func (s *AttendanceService) CheckIn(userID string, req *models.AttendanceRequest) (*models.Attendance, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	if !s.IsValidLocation(req.Latitude, req.Longitude) {
		return nil, errors.New("location is outside school area")
	}

	today := time.Now().UTC()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	collection := s.db.Collection("attendances")
	var existingAttendance models.Attendance
	err = collection.FindOne(s.ctx, bson.M{
		"user_id": objectID,
		"date": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		},
	}).Decode(&existingAttendance)

	if err == nil {
		return nil, errors.New("already checked in today")
	}

	now := time.Now().UTC()
	status := s.DetermineStatus(now)

	attendance := models.Attendance{
		UserID:  objectID,
		Date:    startOfDay,
		CheckIn: &now,
		Status:  status,
		Location: models.Location{
			Latitude:  req.Latitude,
			Longitude: req.Longitude,
			Address:   req.Address,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	result, err := collection.InsertOne(s.ctx, attendance)
	if err != nil {
		log.Printf("Error creating attendance: %v", err)
		return nil, errors.New("failed to create attendance")
	}

	attendance.ID = result.InsertedID.(primitive.ObjectID)
	log.Printf("Check in successful for user %s at %s", userID, now.Format("15:04:05"))
	
	return &attendance, nil
}

func (s *AttendanceService) CheckOut(userID string, req *models.AttendanceRequest) (*models.Attendance, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	if !s.IsValidLocation(req.Latitude, req.Longitude) {
		return nil, errors.New("location is outside school area")
	}

	today := time.Now().UTC()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	collection := s.db.Collection("attendances")
	var attendance models.Attendance
	err = collection.FindOne(s.ctx, bson.M{
		"user_id": objectID,
		"date": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		},
	}).Decode(&attendance)

	if err != nil {
		return nil, errors.New("no check in record found for today")
	}

	if attendance.CheckOut != nil {
		return nil, errors.New("already checked out today")
	}

	now := time.Now().UTC()
	_, err = collection.UpdateOne(
		s.ctx,
		bson.M{"_id": attendance.ID},
		bson.M{
			"$set": bson.M{
				"check_out":  now,
				"updated_at": now,
			},
		},
	)

	if err != nil {
		log.Printf("Error updating attendance: %v", err)
		return nil, errors.New("failed to update attendance")
	}

	attendance.CheckOut = &now
	attendance.UpdatedAt = now

	log.Printf("Check out successful for user %s at %s", userID, now.Format("15:04:05"))
	return &attendance, nil
}

func (s *AttendanceService) GetTodayAttendance(userID string) (*models.Attendance, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	today := time.Now().UTC()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	collection := s.db.Collection("attendances")
	var attendance models.Attendance
	err = collection.FindOne(s.ctx, bson.M{
		"user_id": objectID,
		"date": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		},
	}).Decode(&attendance)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.New("failed to fetch attendance")
	}

	return &attendance, nil
}

func (s *AttendanceService) GetAttendanceHistory(userID string, limit, offset int) ([]models.Attendance, int64, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, 0, errors.New("invalid user ID")
	}

	collection := s.db.Collection("attendances")
	
	total, err := collection.CountDocuments(s.ctx, bson.M{"user_id": objectID})
	if err != nil {
		return nil, 0, errors.New("failed to count attendance records")
	}

	cursor, err := collection.Find(
		s.ctx,
		bson.M{"user_id": objectID},
		options.Find().SetSort(bson.D{{Key: "date", Value: -1}}).SetSkip(int64(offset)).SetLimit(int64(limit)),
	)
	if err != nil {
		return nil, 0, errors.New("failed to fetch attendance history")
	}
	defer cursor.Close(s.ctx)

	var attendances []models.Attendance
	if err = cursor.All(s.ctx, &attendances); err != nil {
		return nil, 0, errors.New("failed to decode attendance records")
	}

	return attendances, total, nil
}

func (s *AttendanceService) GetAttendanceStats(userID string) (*models.AttendanceStats, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	collection := s.db.Collection("attendances")
	
	pipeline := []bson.M{
		{"$match": bson.M{"user_id": objectID}},
		{"$group": bson.M{
			"_id": "$status",
			"count": bson.M{"$sum": 1},
		}},
	}

	cursor, err := collection.Aggregate(s.ctx, pipeline)
	if err != nil {
		return nil, errors.New("failed to calculate stats")
	}
	defer cursor.Close(s.ctx)

	var results []bson.M
	if err = cursor.All(s.ctx, &results); err != nil {
		return nil, errors.New("failed to decode stats")
	}

	stats := &models.AttendanceStats{}
	var totalDays int

	for _, result := range results {
		status := result["_id"].(string)
		count := int(result["count"].(int32))
		totalDays += count

		switch status {
		case "present":
			stats.TotalPresent = count
		case "late":
			stats.TotalLate = count
		case "absent":
			stats.TotalAbsent = count
		}
	}

	if totalDays > 0 {
		stats.Percentage = float64(stats.TotalPresent) / float64(totalDays) * 100
	}

	return stats, nil
}

func (s *AttendanceService) GetAttendanceByDate(userID string, date time.Time) (*models.Attendance, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	collection := s.db.Collection("attendances")
	var attendance models.Attendance
	err = collection.FindOne(s.ctx, bson.M{
		"user_id": objectID,
		"date": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		},
	}).Decode(&attendance)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.New("failed to fetch attendance")
	}

	return &attendance, nil
}

func (s *AttendanceService) IsValidLocation(lat, lng float64) bool {
	schoolLat, schoolLng, maxDistance := s.config.GetSchoolLocation()
	distance := utils.CalculateDistance(lat, lng, schoolLat, schoolLng)
	return distance <= maxDistance
}

func (s *AttendanceService) calculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	return utils.CalculateDistance(lat1, lng1, lat2, lng2)
}

func (s *AttendanceService) DetermineStatus(checkInTime time.Time) string {
	startHour, startMinute, _, _ := s.config.GetSchoolHours()
	lateThreshold := s.config.GetLateThreshold()
	
	loc, _ := time.LoadLocation("Asia/Jakarta")
	localTime := checkInTime.In(loc)
	
	schoolStartTime := time.Date(localTime.Year(), localTime.Month(), localTime.Day(), startHour, startMinute, 0, 0, loc)
	
	if localTime.Before(schoolStartTime) || localTime.Equal(schoolStartTime) {
		return "present"
	} else if localTime.Before(schoolStartTime.Add(time.Duration(lateThreshold) * time.Minute)) {
		return "late"
	} else {
		return "absent"
	}
}
