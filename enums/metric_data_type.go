package enums

import (
	"database/sql/driver"
	"errors"
)

// MetricDataType represents the data type of a metric
type MetricDataType string

const (
	Integer MetricDataType = "integer"
	Decimal MetricDataType = "decimal"
	Boolean MetricDataType = "boolean"
	String  MetricDataType = "string"
)

// Scan implements the sql.Scanner interface for MetricDataType
func (dt *MetricDataType) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("invalid data for MetricDataType")
	}

	switch str {
	case string(Integer), string(Decimal), string(Boolean), string(String):
		*dt = MetricDataType(str)
		return nil
	default:
		return errors.New("invalid value for MetricDataType")
	}
}

// Value implements the driver.Valuer interface for MetricDataType
func (dt MetricDataType) Value() (driver.Value, error) {
	switch dt {
	case Integer, Decimal, Boolean, String:
		return string(dt), nil
	default:
		return nil, errors.New("invalid MetricDataType")
	}
}

// Valid checks if the enum value is valid
func (dt MetricDataType) Valid() bool {
	switch dt {
	case Integer, Decimal, Boolean, String:
		return true
	}
	return false
}

// GetValidMetricDataTypes returns all valid metric data types
func GetValidMetricDataTypes() []string {
	return []string{
		string(Integer),
		string(Decimal),
		string(Boolean),
		string(String),
	}
}
