package enums

import (
	"database/sql/driver"
	"errors"
)

// VisibilityScope represents who can view the leaderboard
type VisibilityScope string

const (
	Public  VisibilityScope = "public"
	Private VisibilityScope = "private"
)

// Scan implements the sql.Scanner interface for VisibilityScope
func (vs *VisibilityScope) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("invalid data for VisibilityScope")
	}

	switch str {
	case string(Public), string(Private):
		*vs = VisibilityScope(str)
		return nil
	default:
		return errors.New("invalid value for VisibilityScope")
	}
}

// Value implements the driver.Valuer interface for VisibilityScope
func (vs VisibilityScope) Value() (driver.Value, error) {
	switch vs {
	case Public, Private:
		return string(vs), nil
	default:
		return nil, errors.New("invalid VisibilityScope")
	}
}

// Valid checks if the enum value is valid
func (vs VisibilityScope) Valid() bool {
	switch vs {
	case Public, Private:
		return true
	}
	return false
}

func GetValidVisibilityScopes() []string {
	return []string{string(Public), string(Private)}
}
