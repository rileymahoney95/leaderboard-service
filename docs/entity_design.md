# Leaderboard Service - Entity Structure Design

## Core Entities

### 1. Leaderboard
- `id`: UUID
- `name`: String
- `description`: String (optional)
- `category`: String (services, sports, competitions, etc.)
- `type`: String (individual, team, group)
- `timeFrame`: Enum (all-time, daily, weekly, monthly, custom)
- `startDate`: DateTime (optional, for custom timeframes)
- `endDate`: DateTime (optional, for custom timeframes)
- `sortOrder`: Enum (ascending, descending)
- `visibilityScope`: Enum (public, private, restricted)
- `maxEntries`: Integer (optional)
- `isActive`: Boolean

### 2. Participant
- `id`: UUID
- `externalId`: String (reference to user/team/entity in external system)
- `name`: String
- `type`: Enum (individual, team, group)
- `metadata`: JSON (for additional properties)

### 3. Metric
- `id`: UUID
- `name`: String
- `description`: String
- `dataType`: Enum (integer, decimal, boolean, timestamp, duration)
- `unit`: String (optional)
- `aggregationType`: Enum (sum, average, count, max, min, latest)
- `resetPeriod`: Enum (none, daily, weekly, monthly)
- `isHigherBetter`: Boolean (determines if higher values rank better)

### 4. MetricValue
- `id`: UUID
- `metricId`: UUID (reference to Metric)
- `participantId`: UUID (reference to Participant)
- `value`: Variant (based on dataType)
- `timestamp`: DateTime
- `source`: String (identifying where/how this value was recorded)
- `context`: JSON (additional context data)

### 5. LeaderboardEntry
- `id`: UUID
- `leaderboardId`: UUID (reference to Leaderboard)
- `participantId`: UUID (reference to Participant)
- `rank`: Integer
- `score`: Decimal
- `lastUpdated`: DateTime

### 6. LeaderboardMetric
- `id`: UUID
- `leaderboardId`: UUID (reference to Leaderboard)
- `metricId`: UUID (reference to Metric)
- `weight`: Decimal (for composite scoring)
- `displayPriority`: Integer

## Supporting Entities

### 7. Event
- `id`: UUID
- `name`: String
- `description`: String (optional)
- `startDate`: DateTime
- `endDate`: DateTime (optional)
- `category`: String
- `metadata`: JSON

### 8. Achievement
- `id`: UUID
- `name`: String
- `description`: String
- `thresholds`: JSON (conditions for unlocking)
- `category`: String
- `pointValue`: Integer (optional)
- `badge`: String (reference to badge image)

### 9. ParticipantAchievement
- `id`: UUID
- `participantId`: UUID (reference to Participant)
- `achievementId`: UUID (reference to Achievement)
- `dateEarned`: DateTime
- `context`: JSON

### 10. Notification
- `id`: UUID
- `participantId`: UUID (reference to Participant)
- `type`: String
- `message`: String
- `read`: Boolean
- `createdAt`: DateTime

## Design Principles

### 1. Decoupling
The core entities are completely decoupled from specific domains. There's no hard-coded reference to tasks, sports, or predictions, making the service domain-agnostic.

### 2. Flexibility
The `Metric` entity allows you to define any type of measurable value, and the `MetricValue` entity stores the actual measurements, enabling any quantifiable attribute to be tracked.

### 3. Composability
Each leaderboard can combine multiple metrics with different weights for complex scoring systems, allowing for sophisticated ranking algorithms.

### 4. Extensibility
The JSON metadata fields let you store domain-specific attributes without changing the schema, future-proofing the design.

## Use Case Implementation Examples

### Services App
- **Metrics**: Define metrics like "tasks_completed", "task_completion_rate", "services_delivered"
- **Leaderboards**: Create time-bound leaderboards (daily, weekly, monthly)
- **Participants**: Track individual service providers
- **Example**: A weekly leaderboard showing top performers by task completion rate

### Pickup League App
- **Metrics**: Define sport-specific metrics ("points_scored", "assists", "rebounds", etc.)
- **Participants**: Support both individual athletes and teams
- **Leaderboards**: Create separate leaderboards for different sports and statistics
- **Events**: Use the Event entity to track games/matches
- **Example**: Basketball season leaderboard showing teams ranked by win percentage

### Competition/Prediction App
- **Metrics**: Define metrics for "prediction_accuracy", "correct_predictions", "streak_length"
- **Context**: Use the MetricValue's context field to store prediction details
- **Leaderboards**: Create leaderboards by competition category or timeframe
- **Example**: Monthly leaderboard showing users with highest prediction accuracy percentages

## Implementation Considerations

### Database Design
- Consider using a relational database for core entities and relationships
- JSON fields can be implemented as JSONB in PostgreSQL or equivalent in other databases
- Index frequently queried fields, especially those used in leaderboard calculations

### API Design
- Create RESTful endpoints for each entity
- Implement specialized endpoints for leaderboard operations:
  - `/leaderboards/{id}/rankings` - Get current rankings
  - `/leaderboards/{id}/calculate` - Trigger recalculation
  - `/participants/{id}/metrics` - Get all metrics for a participant

### Performance Optimization
- Consider caching frequently accessed leaderboards
- Implement background jobs for regular leaderboard recalculations
- Use materialized views for complex aggregations

### Security Considerations
- Implement proper authentication and authorization
- Consider rate limiting for leaderboard calculations
- Ensure participant privacy based on visibility settings