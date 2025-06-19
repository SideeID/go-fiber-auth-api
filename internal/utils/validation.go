package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidatorErrors(err error) []string {
	var errors []string
	
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			switch e.Tag() {
			case "required":
				errors = append(errors, e.Field()+" is required")
			case "email":
				errors = append(errors, e.Field()+" must be a valid email")
			case "min":
				errors = append(errors, e.Field()+" must be at least "+e.Param()+" characters")
			case "max":
				errors = append(errors, e.Field()+" must be at most "+e.Param()+" characters")
			case "url":
				errors = append(errors, e.Field()+" must be a valid URL")
			default:
				errors = append(errors, e.Field()+" is invalid")
			}
		}
	} else {
		errors = append(errors, err.Error())
	}
	
	return errors
}

func IsValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func IsValidPassword(password string) bool {
	return len(password) >= 6
}