Minnect Use Case – Leaderboard Implementation
This document describes how the generic leaderboard service entity structure can be used for Minnect.

Mapping the Entities
1. Leaderboard
A Leaderboard represents a ranking view. In Minnect, you might define several leaderboard instances:

User Leaderboard – New Users:

Category: "Minnect"

Type: "individual"

TimeFrame: "monthly"

Metadata (or a naming convention): Indicates that only new users are included

User Leaderboard – Existing Users:

Similar to the new user leaderboard but filtered by registration date (stored in the Participant’s metadata)

Expert Leaderboard:

Category: "Minnect"

Type: "individual" (or you might define a separate type such as "expert")

TimeFrame: "monthly"

Key fields in the Leaderboard entity (such as timeFrame, category, and visibilityScope) allow you to tailor each leaderboard to the intended audience and period.

2. Participant
Each Minnect user is stored as a Participant. Use:

externalId: to reference the Minnect user record

name: for display

type: to distinguish between regular users and experts (or use a metadata field to mark experts)

3. Metric & MetricValue
Define a set of Metric entities that capture the performance measures:

For User Leaderboards:

Metric: "monthly_calls_completed"

Aggregation: sum

ResetPeriod: monthly

Metric: "monthly_texts_answered"

Aggregation: sum

ResetPeriod: monthly

Metric: "user_score"

Calculated using a custom algorithm based on the above values (and perhaps other factors)

For Expert Leaderboards:

Metric: "monthly_calls_completed"

Metric: "monthly_texts_answered"

Metric: "response_rate"

Could be defined as a percentage (e.g., calls/texts answered within a target time)

Metric: "expert_score"

A composite score calculated using a custom algorithm that factors in response rate along with call and text volume

When a minnect (call or text) is completed, the system records a MetricValue for the appropriate Metric. The MetricValue includes:

The value (for instance, an increment of 1)

The timestamp (to enable monthly aggregation)

A context (which could include whether the minnect was a call or text)

4. LeaderboardMetric & LeaderboardEntry
LeaderboardMetric defines which metrics contribute to a leaderboard and their weights. For example, on the expert leaderboard:

You might assign a higher weight to response rate versus volume.

LeaderboardEntry stores the computed rank and composite score for each Participant on a given leaderboard. After a monthly aggregation and application of your custom scoring algorithm, each expert or user will have a corresponding LeaderboardEntry record.

Example Workflow
Recording Events:
Every time a call or text (minnect) is completed, your system creates a MetricValue record. For instance:

A completed call by an expert increases the count for "monthly_calls_completed."

A completed text by a user increases the count for "monthly_texts_answered."

Aggregation & Scoring:
A scheduled job (or a real-time aggregation process) periodically aggregates the monthly MetricValue data for each Participant:

The custom algorithm computes a score by combining the aggregated metrics.

For experts, the response rate (e.g., percentage of texts answered within a target time) is factored in.

LeaderboardMetric definitions provide the weights for each metric.

Leaderboard Update:
The system creates or updates LeaderboardEntry records with:

The rank (position in the leaderboard)

The computed score

A reference to the associated Leaderboard (user new, user existing, or expert)

Retrieving the Leaderboard:
An API endpoint (e.g., /leaderboards/{id}/rankings) fetches the LeaderboardEntry records, allowing you to display the leaderboard on the Minnect app.

Diagram
Below is a Mermaid.js diagram that shows how the entities interrelate for this use case:

mermaid
Copy
erDiagram
    PARTICIPANT {
        UUID id PK
        string externalId
        string name
        string type
        json metadata
    }
    LEADERBOARD {
        UUID id PK
        string name
        string description
        string category
        string type
        string timeFrame
        timestamp startDate
        timestamp endDate
        string sortOrder
        string visibilityScope
        int maxEntries
        boolean isActive
    }
    METRIC {
        UUID id PK
        string name
        string description
        string dataType
        string unit
        string aggregationType
        string resetPeriod
        boolean isHigherBetter
    }
    METRIC_VALUE {
        UUID id PK
        UUID metricId FK
        UUID participantId FK
        float value
        timestamp timestamp
        string source
        json context
    }
    LEADERBOARD_ENTRY {
        UUID id PK
        UUID leaderboardId FK
        UUID participantId FK
        int rank
        float score
        timestamp lastUpdated
    }
    LEADERBOARD_METRIC {
        UUID id PK
        UUID leaderboardId FK
        UUID metricId FK
        float weight
        int displayPriority
    }

    PARTICIPANT ||--o{ METRIC_VALUE : "records"
    LEADERBOARD ||--o{ LEADERBOARD_ENTRY : "has"
    LEADERBOARD ||--o{ LEADERBOARD_METRIC : "defines"
    METRIC ||--o{ METRIC_VALUE : "captures"
Conclusion
By mapping Minnect's requirements to your entity structure, you achieve the following:

Flexibility: The same abstract entities (Leaderboard, Participant, Metric, etc.) are used for both user and expert perspectives.

Scalability: New metrics can be added easily (e.g., response rate, custom scores) without changing the core design.

Decoupling: The leaderboard service remains domain-agnostic, meaning it can be reused in other apps or use cases by simply defining different metrics and weights.

This design gives you a robust foundation to implement leaderboards for Minnect that can evolve as your business logic changes. Let me know if you need further examples or sample code for any part of this workflow!