package enums

import (
	"database/sql/driver"
	"errors"
)

// SortOrder represents the sorting direction for leaderboard entries
type SortOrder string

const (
	Ascending  SortOrder = "ascending"
	Descending SortOrder = "descending"
)

// Scan implements the sql.Scanner interface for SortOrder
func (so *SortOrder) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("invalid data for SortOrder")
	}

	switch str {
	case string(Ascending), string(Descending):
		*so = SortOrder(str)
		return nil
	default:
		return errors.New("invalid value for SortOrder")
	}
}

// Value implements the driver.Valuer interface for SortOrder
func (so SortOrder) Value() (driver.Value, error) {
	switch so {
	case Ascending, Descending:
		return string(so), nil
	default:
		return nil, errors.New("invalid SortOrder")
	}
}

// Valid checks if the enum value is valid
func (so SortOrder) Valid() bool {
	switch so {
	case Ascending, Descending:
		return true
	}
	return false
}

func GetValidSortOrders() []string {
	return []string{string(Ascending), string(Descending)}
}
