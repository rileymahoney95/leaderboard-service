package validation

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

// Test structure for validation
type TestStruct struct {
	Name      string  `json:"name" validate:"required"`
	TimeFrame string  `json:"time_frame" validate:"required,oneof=daily weekly monthly yearly all_time custom,custom_timeframe"`
	StartDate *string `json:"start_date,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z"`
	EndDate   *string `json:"end_date,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z"`
}

func TestCustomTimeframeValidation(t *testing.T) {
	// Helper function to create string pointer
	stringPtr := func(s string) *string {
		return &s
	}

	testCases := []struct {
		name          string
		input         TestStruct
		expectedValid bool
	}{
		{
			name: "Valid non-custom timeframe",
			input: TestStruct{
				Name:      "Test Leaderboard",
				TimeFrame: "weekly",
			},
			expectedValid: true,
		},
		{
			name: "Custom timeframe with both dates",
			input: TestStruct{
				Name:      "Test Custom Leaderboard",
				TimeFrame: "custom",
				StartDate: stringPtr("2023-01-01T00:00:00Z"),
				EndDate:   stringPtr("2023-01-31T23:59:59Z"),
			},
			expectedValid: true,
		},
		{
			name: "Custom timeframe with missing dates",
			input: TestStruct{
				Name:      "Test Custom Leaderboard",
				TimeFrame: "custom",
			},
			expectedValid: false,
		},
		{
			name: "Custom timeframe with only start date",
			input: TestStruct{
				Name:      "Test Custom Leaderboard",
				TimeFrame: "custom",
				StartDate: stringPtr("2023-01-01T00:00:00Z"),
			},
			expectedValid: false,
		},
		{
			name: "Custom timeframe with only end date",
			input: TestStruct{
				Name:      "Test Custom Leaderboard",
				TimeFrame: "custom",
				EndDate:   stringPtr("2023-01-31T23:59:59Z"),
			},
			expectedValid: false,
		},
		{
			name: "Invalid datetime format",
			input: TestStruct{
				Name:      "Test Custom Leaderboard",
				TimeFrame: "custom",
				StartDate: stringPtr("2023-01-01"), // Missing time part
				EndDate:   stringPtr("2023-01-31T23:59:59Z"),
			},
			expectedValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := Validate.Struct(tc.input)
			isValid := err == nil

			if isValid != tc.expectedValid {
				if tc.expectedValid {
					t.Errorf("Expected validation to pass but got error: %v", err)
				} else {
					t.Errorf("Expected validation to fail but it passed")
				}
			}

			// If we expect validation to fail, let's verify it's for the right reason
			if !tc.expectedValid && err != nil {
				validationErrors, ok := err.(validator.ValidationErrors)
				if !ok {
					t.Errorf("Expected validator.ValidationErrors but got %T", err)
					return
				}

				// Check if we have at least one error with the custom_timeframe tag
				if tc.input.TimeFrame == "custom" && (tc.input.StartDate == nil || tc.input.EndDate == nil) {
					foundCustomTimeFrameError := false
					for _, e := range validationErrors {
						if e.Tag() == "custom_timeframe" {
							foundCustomTimeFrameError = true
							break
						}
					}

					if !foundCustomTimeFrameError {
						t.Errorf("Expected custom_timeframe validation error but got: %v", validationErrors)
					}
				}
			}
		})
	}
}
