package middleware

import (
	"math"
	"strconv"
	"time"
	"ujikom-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func LocationValidationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Path() != "/api/v1/attendance/checkin" && c.Path() != "/api/v1/attendance/checkout" {
			return c.Next()
		}

		var requestBody struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		}

		if err := c.BodyParser(&requestBody); err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
		}

		if !isValidCoordinates(requestBody.Latitude, requestBody.Longitude) {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid GPS coordinates")
		}

		if !isInIndonesia(requestBody.Latitude, requestBody.Longitude) {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Location must be within Indonesia")
		}

		return c.Next()
	}
}

func SecurityHeadersMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Security headers untuk mobile app
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Set("Content-Security-Policy", "default-src 'self'")
		
		// Mobile-specific headers
		c.Set("X-Mobile-API", "true")
		c.Set("X-GPS-Required", "true")
		c.Set("X-Location-Based", "true")
		
		return c.Next()
	}
}

// middleware untuk validasi device mobile
func DeviceValidationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userAgent := c.Get("User-Agent")
		
		if userAgent == "" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "User-Agent header is required")
		}

		// Set device info ke context
		c.Locals("user_agent", userAgent)
		c.Locals("device_type", detectDeviceType(userAgent))
		
		return c.Next()
	}
}

func RateLimitMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		clientIP := c.IP()
		c.Set("X-Client-IP", clientIP)
		
		c.Set("X-Request-ID", generateRequestID())
		
		return c.Next()
	}
}

func isValidCoordinates(lat, lng float64) bool {
	if lat < -90 || lat > 90 {
		return false
	}
	if lng < -180 || lng > 180 {
		return false
	}
	return true
}

func isInIndonesia(lat, lng float64) bool {
	return lat >= -11.0 && lat <= 6.0 && lng >= 95.0 && lng <= 141.0
}

func detectDeviceType(userAgent string) string {
	if userAgent == "" {
		return "unknown"
	}
	
	return "mobile"
}

func generateRequestID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}

func calculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadius = 6371 // km

	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}
