package controllers

import (
	"log"
	"math"
	"ujikom-backend/internal/config"
	"ujikom-backend/internal/models"
	"ujikom-backend/internal/services"
	"ujikom-backend/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type AttendanceController struct {
	db                *mongo.Database
	validator         *validator.Validate
	attendanceService services.AttendanceServiceInterface
	config            *config.Config
}

func NewAttendanceController(db *mongo.Database, cfg *config.Config) *AttendanceController {
	return &AttendanceController{
		db:                db,
		validator:         validator.New(),
		attendanceService: services.NewAttendanceService(db, cfg),
		config:            cfg,
	}
}

func (ac *AttendanceController) CheckIn(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not found")
	}

	var req models.AttendanceRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := ac.validator.Struct(req); err != nil {
		return utils.ValidationErrorResponse(c, utils.ValidatorErrors(err))
	}

	attendance, err := ac.attendanceService.CheckIn(user.ID.Hex(), &req)
	if err != nil {
		log.Printf("CheckIn error for user %s: %v", user.Name, err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	response := models.AttendanceResponse{
		ID:        attendance.ID,
		UserID:    attendance.UserID,
		Date:      attendance.Date,
		CheckIn:   attendance.CheckIn,
		CheckOut:  attendance.CheckOut,
		Status:    attendance.Status,
		Location:  attendance.Location,
		User:      ac.createUserPublic(user),
		CreatedAt: attendance.CreatedAt,
		UpdatedAt: attendance.UpdatedAt,
	}

	log.Printf("User %s checked in successfully at %s", user.Name, attendance.CheckIn.Format("15:04:05"))
	return utils.SuccessResponse(c, "Check in successful", response)
}

func (ac *AttendanceController) CheckOut(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not found")
	}

	var req models.AttendanceRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := ac.validator.Struct(req); err != nil {
		return utils.ValidationErrorResponse(c, utils.ValidatorErrors(err))
	}

	attendance, err := ac.attendanceService.CheckOut(user.ID.Hex(), &req)
	if err != nil {
		log.Printf("CheckOut error for user %s: %v", user.Name, err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	response := models.AttendanceResponse{
		ID:        attendance.ID,
		UserID:    attendance.UserID,
		Date:      attendance.Date,
		CheckIn:   attendance.CheckIn,
		CheckOut:  attendance.CheckOut,
		Status:    attendance.Status,
		Location:  attendance.Location,
		User:      ac.createUserPublic(user),
		CreatedAt: attendance.CreatedAt,
		UpdatedAt: attendance.UpdatedAt,
	}

	log.Printf("User %s checked out successfully at %s", user.Name, attendance.CheckOut.Format("15:04:05"))
	return utils.SuccessResponse(c, "Check out successful", response)
}

func (ac *AttendanceController) GetTodayAttendance(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not found")
	}

	attendance, err := ac.attendanceService.GetTodayAttendance(user.ID.Hex())
	if err != nil {
		log.Printf("GetTodayAttendance error for user %s: %v", user.Name, err)
		return utils.SuccessResponse(c, "No attendance record for today", nil)
	}

	response := models.AttendanceResponse{
		ID:        attendance.ID,
		UserID:    attendance.UserID,
		Date:      attendance.Date,
		CheckIn:   attendance.CheckIn,
		CheckOut:  attendance.CheckOut,
		Status:    attendance.Status,
		Location:  attendance.Location,
		User:      ac.createUserPublic(user),
		CreatedAt: attendance.CreatedAt,
		UpdatedAt: attendance.UpdatedAt,
	}

	return utils.SuccessResponse(c, "Today's attendance retrieved", response)
}

func (ac *AttendanceController) GetAttendanceHistory(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not found")
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	attendances, total, err := ac.attendanceService.GetAttendanceHistory(user.ID.Hex(), limit, offset)
	if err != nil {
		log.Printf("GetAttendanceHistory error for user %s: %v", user.Name, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch attendance history")
	}

	var history []models.AttendanceHistory
	for _, attendance := range attendances {
		history = append(history, models.AttendanceHistory{
			Date:     attendance.Date,
			CheckIn:  attendance.CheckIn,
			CheckOut: attendance.CheckOut,
			Status:   attendance.Status,
			Location: attendance.Location,
		})
	}

	response := fiber.Map{
		"history": history,
		"pagination": fiber.Map{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"total_pages": int(math.Ceil(float64(total) / float64(limit))),
		},
	}

	return utils.SuccessResponse(c, "Attendance history retrieved", response)
}

func (ac *AttendanceController) GetAttendanceStats(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not found")
	}

	stats, err := ac.attendanceService.GetAttendanceStats(user.ID.Hex())
	if err != nil {
		log.Printf("GetAttendanceStats error for user %s: %v", user.Name, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to calculate stats")
	}

	return utils.SuccessResponse(c, "Attendance statistics retrieved", stats)
}

func (ac *AttendanceController) createUserPublic(user models.User) models.UserPublic {
	return models.UserPublic{
		ID:      user.ID,
		NIS:     user.NIS,
		Name:    user.Name,
		Kelas:   user.Kelas,
		Jurusan: user.Jurusan,
		Email:   user.Email,
		Phone:   user.Phone,
	}
}
