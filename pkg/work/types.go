// Package work defines Ploeg's core semantics: work items, leases,
// checkpoints, and outcomes. See docs/design.md §3.
package work

import "time"

// State is the lifecycle position of a WorkItem inside Ploeg.
// The tracker item it mirrors remains the source of truth for content.
type State string

const (
	StateIngested State = "ingested"
	StateQueued   State = "queued"
	StateLeased   State = "leased"
	StateStale    State = "stale"
	StateDone     State = "done"
)

// Outcome is the terminal report of a run. Stuck carries a mandatory
// reason and routes to a human queue.
type Outcome string

const (
	OutcomePROpened        Outcome = "pr_opened"
	OutcomePRUpdated       Outcome = "pr_updated"
	OutcomeIssueUpdated    Outcome = "issue_updated"
	OutcomeFollowUpCreated Outcome = "follow_up_created"
	OutcomeStuck           Outcome = "stuck"
	OutcomeFailed          Outcome = "failed"
	OutcomeNoChangeNeeded  Outcome = "no_change_needed"
)

// WorkItem mirrors one tracker item (provider + external id + revision).
type WorkItem struct {
	ID         string
	Provider   string // tracker provider name, e.g. "vikunja"
	ExternalID string // provider-scoped id of the tracker item
	Revision   string // provider revision/etag for staleness detection
	Team       string // team the item is queued for; empty until assigned
	State      State
	Title      string
	URL        string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Lease is a team's crash-safe claim on a work item. Renewed by the
// running Job; expiry releases the item mechanically.
type Lease struct {
	WorkItemID string
	Team       string
	ExpiresAt  time.Time
	RenewedAt  time.Time
}

// Checkpoint is the small durable progress record enabling resume.
// Everything else is re-derived from git/forge state.
type Checkpoint struct {
	WorkItemID string
	Phase      string // e.g. "branch_created", "changes_made", "pr_opened"
	Branch     string
	PRURL      string
	At         time.Time
}
