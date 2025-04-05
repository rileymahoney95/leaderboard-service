package enums

import (
	"database/sql/driver"
	"errors"
)

// LeaderboardType represents the type of leaderboard (individual or team)
type LeaderboardType string

const (
	Individual LeaderboardType = "individual"
	Team       LeaderboardType = "team"
)

// Scan implements the sql.Scanner interface for LeaderboardType
func (lt *LeaderboardType) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("invalid data for LeaderboardType")
	}

	switch str {
	case string(Individual), string(Team):
		*lt = LeaderboardType(str)
		return nil
	default:
		return errors.New("invalid value for LeaderboardType")
	}
}

// Value implements the driver.Valuer interface for LeaderboardType
func (lt LeaderboardType) Value() (driver.Value, error) {
	switch lt {
	case Individual, Team:
		return string(lt), nil
	default:
		return nil, errors.New("invalid LeaderboardType")
	}
}

// Valid checks if the enum value is valid
func (lt LeaderboardType) Valid() bool {
	switch lt {
	case Individual, Team:
		return true
	}
	return false
}

func GetValidLeaderboardTypes() []string {
	return []string{string(Individual), string(Team)}
}
