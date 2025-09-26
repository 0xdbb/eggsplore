package server

import (
	"errors"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Custom Validation Error Messages
func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " field is required."
	case "email":
		return fe.Field() + " must be a valid email address."
	case "gte":
		return fe.Field() + " must be greater than or equal to " + fe.Param() + "."
	case "lte":
		return fe.Field() + " must be less than or equal to " + fe.Param() + "."
	case "min":
		return fe.Field() + " must be at least " + fe.Param() + " characters long."
	case "max":
		return fe.Field() + " must be at most " + fe.Param() + " characters long."
	case "len":
		return fe.Field() + " must be exactly " + fe.Param() + " characters long."
	case "oneof":
		return fe.Field() + " must be one of: " + fe.Param() + "."
	case "alphanum":
		return fe.Field() + " must contain only alphanumeric characters."
	case "StrongPassword":
		return "Password must be at least 8 characters long, include uppercase and lowercase letters, a number, and a special character."
	case "ValidUsername":
		return fe.Field() + " may only contain letters, numbers, and underscores."
	default:
		return "Validation failed on field " + fe.Field() + "."
	}
}

// strongPassword is a custom validation function for strong passwords
var StrongPassword validator.Func = func(fieldLevel validator.FieldLevel) bool {
	password := fieldLevel.Field().String()
	var (
		minLength    = 8
		hasUppercase = regexp.MustCompile(`[A-Z]`).MatchString
		hasLowercase = regexp.MustCompile(`[a-z]`).MatchString
		hasNumber    = regexp.MustCompile(`[0-9]`).MatchString
		hasSpecial   = regexp.MustCompile(`[!@#{}\[\]\$%\^&\*\(\)]`).MatchString
	)

	return len(password) >= minLength &&
		hasUppercase(password) &&
		hasLowercase(password) &&
		hasNumber(password) &&
		hasSpecial(password)
}

// ValidUsername checks that the username contains only letters, numbers, and underscores.
var ValidUsername validator.Func = func(fieldLevel validator.FieldLevel) bool {
	username := fieldLevel.Field().String()
	// ^[a-zA-Z0-9_]+$ → one or more letters, numbers, or underscores
	matched, err := regexp.MatchString(`^[a-zA-Z0-9_]+$`, username)
	return err == nil && matched
}

// HandleValidationError returns the first validation error only
func HandleValidationError(err error) *ErrorResponse {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) && len(ve) > 0 {
		msgs := make([]string, len(ve))
		for i, fe := range ve {
			msgs[i] = msgForTag(fe)
		}
		return &ErrorResponse{
			Status:  "error",
			Message: strings.Join(msgs, ", "),
			Code:    400,
		}
	}
	return nil
}
