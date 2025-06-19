package utils

import (
	"html"
	"regexp"
	"strings"
)

func SanitizeInput(input string) string {
	input = strings.TrimSpace(input)
	
	input = html.EscapeString(input)
	
	return input
}

func SanitizeEmail(email string) string {
	email = strings.TrimSpace(strings.ToLower(email))
	
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return ""
	}
	
	return email
}

func SanitizeName(name string) string {
	name = strings.TrimSpace(name)
	
	nameRegex := regexp.MustCompile(`[^a-zA-Z\s\-\']`)
	name = nameRegex.ReplaceAllString(name, "")
	
	spaceRegex := regexp.MustCompile(`\s+`)
	name = spaceRegex.ReplaceAllString(name, " ")
	
	return strings.TrimSpace(name)
}

func SanitizePhone(phone string) string {
	phoneRegex := regexp.MustCompile(`[^0-9\+\-\s]`)
	phone = phoneRegex.ReplaceAllString(phone, "")
	
	phone = strings.TrimSpace(phone)
	
	return phone
}

func SanitizeURL(url string) string {
	url = strings.TrimSpace(url)
	
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return ""
	}
	
	return url
}

func ValidatePasswordStrength(password string) []string {
	var errors []string
	
	if len(password) < 8 {
		errors = append(errors, "Password must be at least 8 characters long")
	}
	
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		errors = append(errors, "Password must contain at least one uppercase letter")
	}
	
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	if !hasLower {
		errors = append(errors, "Password must contain at least one lowercase letter")
	}
	
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumber {
		errors = append(errors, "Password must contain at least one number")
	}
	
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)
	if !hasSpecial {
		errors = append(errors, "Password must contain at least one special character")
	}
	
	return errors
}
