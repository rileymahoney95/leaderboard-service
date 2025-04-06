package enums

import (
	"database/sql/driver"
	"errors"
)

// ResetPeriod represents when metric values should be reset
type ResetPeriod string

const (
	NoReset      ResetPeriod = "none"
	DailyReset   ResetPeriod = "daily"
	WeeklyReset  ResetPeriod = "weekly"
	MonthlyReset ResetPeriod = "monthly"
	YearlyReset  ResetPeriod = "yearly"
)

// Scan implements the sql.Scanner interface for ResetPeriod
func (rp *ResetPeriod) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("invalid data for ResetPeriod")
	}

	switch str {
	case string(NoReset), string(DailyReset), string(WeeklyReset), string(MonthlyReset), string(YearlyReset):
		*rp = ResetPeriod(str)
		return nil
	default:
		return errors.New("invalid value for ResetPeriod")
	}
}

// Value implements the driver.Valuer interface for ResetPeriod
func (rp ResetPeriod) Value() (driver.Value, error) {
	switch rp {
	case NoReset, DailyReset, WeeklyReset, MonthlyReset, YearlyReset:
		return string(rp), nil
	default:
		return nil, errors.New("invalid ResetPeriod")
	}
}

// Valid checks if the enum value is valid
func (rp ResetPeriod) Valid() bool {
	switch rp {
	case NoReset, DailyReset, WeeklyReset, MonthlyReset, YearlyReset:
		return true
	}
	return false
}

// GetValidResetPeriods returns all valid reset periods
func GetValidResetPeriods() []string {
	return []string{
		string(NoReset),
		string(DailyReset),
		string(WeeklyReset),
		string(MonthlyReset),
		string(YearlyReset),
	}
}
