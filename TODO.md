# TODO

## 1. Add Metric entities

### Metric
Purpose:
Define the types of measurable values (e.g., "monthly_calls_completed," "monthly_texts_answered," "response_rate," etc.) that will be used in your leaderboards.

```golang
type Metric struct {
	BaseModel
	Name            string `gorm:"not null"`
	Description     string `gorm:"type:text"`
	DataType        string `gorm:"not null"` // e.g., "integer", "decimal", "boolean"
	Unit            string // e.g., "calls", "texts", "%"
	AggregationType string `gorm:"not null"` // e.g., "sum", "average", "count"
	ResetPeriod     string `gorm:"not null"` // e.g., "none", "daily", "monthly"
	IsHigherBetter  bool   `gorm:"not null"`


	// Association to MetricValues (optional, for preloading)
	Values []MetricValue `gorm:"foreignKey:MetricID;references:ID"`
}
```

### MetricValue
Purpose:
Store the actual recorded values for a given metric for each participant. Every time a "minnect" (call or text) is completed, a new MetricValue is created.

```golang
type MetricValue struct {
	BaseModel
	MetricID      uuid.UUID   `gorm:"type:uuid;not null"`
	ParticipantID uuid.UUID   `gorm:"type:uuid;not null"`
	// Value is recorded as float64, even if it's an integer. You can convert as needed.
	Value     float64   `gorm:"not null"`
	Timestamp time.Time `gorm:"not null"`
	Source    string    // Identifies where/how this value was recorded
	Context   interface{} `gorm:"type:jsonb"` // For any additional data (e.g., distinguishing call vs. text)

  // Optional associations
	Metric      Metric      `gorm:"foreignKey:MetricID"`
	Participant Participant `gorm:"foreignKey:ParticipantID"`

}
```


### Associate with other models

// Participant Association to MetricValues (optional)
	MetricValues []MetricValue `gorm:"foreignKey:ParticipantID;references:ID"`