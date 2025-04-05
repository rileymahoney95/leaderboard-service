package enums

import (
	"database/sql/driver"
	"errors"
)

// TimeFrame represents the time period for a leaderboard
type TimeFrame string

const (
	Daily   TimeFrame = "daily"
	Weekly  TimeFrame = "weekly"
	Monthly TimeFrame = "monthly"
	Yearly  TimeFrame = "yearly"
	AllTime TimeFrame = "all-time"
)

// Scan implements the sql.Scanner interface for TimeFrame
func (tf *TimeFrame) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("invalid data for TimeFrame")
	}

	switch str {
	case string(Daily), string(Weekly), string(Monthly), string(Yearly), string(AllTime):
		*tf = TimeFrame(str)
		return nil
	default:
		return errors.New("invalid value for TimeFrame")
	}
}

// Value implements the driver.Valuer interface for TimeFrame
func (tf TimeFrame) Value() (driver.Value, error) {
	switch tf {
	case Daily, Weekly, Monthly, Yearly, AllTime:
		return string(tf), nil
	default:
		return nil, errors.New("invalid TimeFrame")
	}
}

// Valid checks if the enum value is valid
func (tf TimeFrame) Valid() bool {
	switch tf {
	case Daily, Weekly, Monthly, Yearly, AllTime:
		return true
	}
	return false
}

// GetValidTimeFrames returns all valid time frame values as strings
func GetValidTimeFrames() []string {
	timeFrames := []TimeFrame{Daily, Weekly, Monthly, Yearly, AllTime}
	result := make([]string, len(timeFrames))

	for i, tf := range timeFrames {
		result[i] = string(tf)
	}

	return result
}
