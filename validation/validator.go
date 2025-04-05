package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Global Validator instance
var Validate *validator.Validate

func init() {
	Validate = validator.New()

	// Register custom validations
	Validate.RegisterValidation("custom_timeframe", validateCustomTimeframe)

	// Use JSON tag names in error messages
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// FormatValidationErrors converts validation errors into a user-friendly error message
func FormatValidationErrors(validationErrors validator.ValidationErrors) error {
	var errMsgs []string
	for _, err := range validationErrors {
		switch err.Tag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("%s is required", err.Field()))
		case "min":
			errMsgs = append(errMsgs, fmt.Sprintf("%s must be at least %s", err.Field(), err.Param()))
		case "max":
			errMsgs = append(errMsgs, fmt.Sprintf("%s must be at most %s", err.Field(), err.Param()))
		case "oneof":
			errMsgs = append(errMsgs, fmt.Sprintf("%s must be one of: %s", err.Field(), err.Param()))
		case "datetime":
			errMsgs = append(errMsgs, fmt.Sprintf("%s must be a valid date-time in format %s", err.Field(), err.Param()))
		case "custom_timeframe":
			errMsgs = append(errMsgs, "When time_frame is 'custom', both start_date and end_date must be provided")
		case "email":
			errMsgs = append(errMsgs, fmt.Sprintf("%s must be a valid email address", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("%s must be a valid URL", err.Field()))
		case "uuid":
			errMsgs = append(errMsgs, fmt.Sprintf("%s must be a valid UUID", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("%s failed validation: %v", err.Field(), err.Tag()))
		}
	}
	return fmt.Errorf("%s", strings.Join(errMsgs, "; "))
}

// Custom validation function to check that when TimeFrame is 'custom', both StartDate and EndDate are provided
func validateCustomTimeframe(fl validator.FieldLevel) bool {
	// Since we're working with structs that are in a different package,
	// we need to use reflection to check the values
	val := fl.Parent()
	timeFrameField := val.FieldByName("TimeFrame")

	// If TimeFrame doesn't exist or is not string type, skip validation
	if !timeFrameField.IsValid() || timeFrameField.Kind() != reflect.String {
		// Try pointer to string
		if timeFrameField.Kind() == reflect.Ptr && !timeFrameField.IsNil() {
			// Get the value the pointer points to
			timeFrameField = timeFrameField.Elem()
			if timeFrameField.Kind() != reflect.String {
				return true
			}
		} else {
			return true
		}
	}

	// Get TimeFrame value
	var timeFrameValue string
	if timeFrameField.Kind() == reflect.String {
		timeFrameValue = timeFrameField.String()
	} else {
		return true
	}

	// If TimeFrame is not 'custom', skip validation
	if timeFrameValue != "custom" {
		return true
	}

	// Check for StartDate and EndDate fields
	startDateField := val.FieldByName("StartDate")
	endDateField := val.FieldByName("EndDate")

	// Both fields must exist and not be nil
	if !startDateField.IsValid() || !endDateField.IsValid() {
		return false
	}

	// For pointer types, they should not be nil
	if startDateField.Kind() == reflect.Ptr && endDateField.Kind() == reflect.Ptr {
		return !startDateField.IsNil() && !endDateField.IsNil()
	}

	return true
}
