package beads

import "time"

// Issue represents a beads task/issue
type Issue struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`       // "open", "closed", "in_progress"
	Priority    int       `json:"priority"`     // 0-3 (P0-P3)
	Type        string    `json:"issue_type"`   // "task", "bug", "feature"
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Labels      []string  `json:"labels,omitempty"`
}

// IssueSummary provides a summary of issues for a branch
type IssueSummary struct {
	OpenCount   int
	ClosedCount int
	TotalCount  int
}
