package enums

import (
	"database/sql/driver"
	"errors"
)

// AggregationType represents how metric values should be aggregated
type AggregationType string

const (
	Sum     AggregationType = "sum"
	Average AggregationType = "average"
	Count   AggregationType = "count"
	Min     AggregationType = "min"
	Max     AggregationType = "max"
	Last    AggregationType = "last"
)

// Scan implements the sql.Scanner interface for AggregationType
func (at *AggregationType) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("invalid data for AggregationType")
	}

	switch str {
	case string(Sum), string(Average), string(Count), string(Min), string(Max), string(Last):
		*at = AggregationType(str)
		return nil
	default:
		return errors.New("invalid value for AggregationType")
	}
}

// Value implements the driver.Valuer interface for AggregationType
func (at AggregationType) Value() (driver.Value, error) {
	switch at {
	case Sum, Average, Count, Min, Max, Last:
		return string(at), nil
	default:
		return nil, errors.New("invalid AggregationType")
	}
}

// Valid checks if the enum value is valid
func (at AggregationType) Valid() bool {
	switch at {
	case Sum, Average, Count, Min, Max, Last:
		return true
	}
	return false
}

// GetValidAggregationTypes returns all valid aggregation types
func GetValidAggregationTypes() []string {
	return []string{
		string(Sum),
		string(Average),
		string(Count),
		string(Min),
		string(Max),
		string(Last),
	}
}
