package utils

import (
	"crypto/rand"
	"encoding/hex"
	"math"
	"strconv"
	"strings"
	"time"
)

type GPSCoordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Accuracy  float64 `json:"accuracy,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type LocationBounds struct {
	North float64 `json:"north"`
	South float64 `json:"south"`
	East  float64 `json:"east"`
	West  float64 `json:"west"`
}

// menghitung jarak antara dua titik GPS menggunakan Haversine formula
func CalculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadius = 6371 // km

	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}

// menghitung jarak dalam meter
func CalculateDistanceInMeters(lat1, lng1, lat2, lng2 float64) float64 {
	return CalculateDistance(lat1, lng1, lat2, lng2) * 1000
}

// mengecek apakah koordinat berada dalam radius tertentu
func IsWithinRadius(centerLat, centerLng, pointLat, pointLng, radiusKm float64) bool {
	distance := CalculateDistance(centerLat, centerLng, pointLat, pointLng)
	return distance <= radiusKm
}

// validasi koordinat GPS
func IsValidGPSCoordinate(lat, lng float64) bool {
	return lat >= -90 && lat <= 90 && lng >= -180 && lng <= 180
}

// mengecek apakah koordinat berada di Indonesia
func IsInIndonesia(lat, lng float64) bool {
	// Bounding box Indonesia (aproximate)
	bounds := LocationBounds{
		North: 6.0,
		South: -11.0,
		East:  141.0,
		West:  95.0,
	}
	
	return lat >= bounds.South && lat <= bounds.North && 
		   lng >= bounds.West && lng <= bounds.East
}

func GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GenerateRequestID() string {
	timestamp := time.Now().UnixNano()
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	return hex.EncodeToString(randomBytes) + "-" + strconv.FormatInt(timestamp, 16)
}

type DeviceInfo struct {
	Type      string `json:"type"`
	OS        string `json:"os"`
	Browser   string `json:"browser"`
	UserAgent string `json:"user_agent"`
}

// deteksi device dari user agent
func DetectDevice(userAgent string) DeviceInfo {
	device := DeviceInfo{
		Type:      "unknown",
		OS:        "unknown",
		Browser:   "unknown",
		UserAgent: userAgent,
	}

	userAgent = strings.ToLower(userAgent)

	if strings.Contains(userAgent, "mobile") || 
	   strings.Contains(userAgent, "android") ||
	   strings.Contains(userAgent, "iphone") ||
	   strings.Contains(userAgent, "ipad") {
		device.Type = "mobile"
	}

	if strings.Contains(userAgent, "android") {
		device.OS = "android"
	} else if strings.Contains(userAgent, "iphone") || strings.Contains(userAgent, "ipad") {
		device.OS = "ios"
	} else if strings.Contains(userAgent, "windows") {
		device.OS = "windows"
	}

	if strings.Contains(userAgent, "chrome") {
		device.Browser = "chrome"
	} else if strings.Contains(userAgent, "firefox") {
		device.Browser = "firefox"
	} else if strings.Contains(userAgent, "safari") {
		device.Browser = "safari"
	}

	return device
}

func IsMobileDevice(userAgent string) bool {
	device := DetectDevice(userAgent)
	return device.Type == "mobile"
}

func GetIndonesianTime() time.Time {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	return time.Now().In(loc)
}

func GetStartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func GetEndOfDay(t time.Time) time.Time {
	return GetStartOfDay(t).Add(24 * time.Hour).Add(-1 * time.Nanosecond)
}

func IsWeekend(t time.Time) bool {
	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

func IsHoliday(t time.Time) bool {
	// libur kemerdekaan
	if t.Month() == time.August && t.Day() == 17 {
		return true
	}
	
	// libur tahun baru
	if t.Month() == time.January && t.Day() == 1 {
		return true
	}
	
	return false
}

func IsSchoolDay(t time.Time) bool {
	return !IsWeekend(t) && !IsHoliday(t)
}

type AttendanceStatus string

const (
	StatusPresent AttendanceStatus = "present"
	StatusLate    AttendanceStatus = "late"
	StatusAbsent  AttendanceStatus = "absent"
)

// menentukan status kehadiran berdasarkan waktu dan config
func DetermineAttendanceStatus(checkInTime time.Time, startHour, startMinute, lateThreshold int) AttendanceStatus {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	localTime := checkInTime.In(loc)
	
	schoolStartTime := time.Date(localTime.Year(), localTime.Month(), localTime.Day(), startHour, startMinute, 0, 0, loc)
	
	lateThresholdDuration := time.Duration(lateThreshold) * time.Minute
	
	if localTime.Before(schoolStartTime) || localTime.Equal(schoolStartTime) {
		return StatusPresent
	} else if localTime.Before(schoolStartTime.Add(lateThresholdDuration)) {
		return StatusLate
	} else {
		return StatusAbsent
	}
}

func GetSchoolHours(startHour, startMinute, endHour, endMinute int) (start, end time.Time) {
	now := GetIndonesianTime()
	loc, _ := time.LoadLocation("Asia/Jakarta")
	
	start = time.Date(now.Year(), now.Month(), now.Day(), startHour, startMinute, 0, 0, loc)
	end = time.Date(now.Year(), now.Month(), now.Day(), endHour, endMinute, 0, 0, loc)
	
	return start, end
}

func IsWithinSchoolHours(t time.Time, startHour, startMinute, endHour, endMinute int) bool {
	start, end := GetSchoolHours(startHour, startMinute, endHour, endMinute)
	return t.After(start) && t.Before(end)
}
