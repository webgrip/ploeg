// Package harness defines the contract between Ploeg and an agent
// container: TaskSpec in, OutcomeReport out. Adapters wrap concrete
// harnesses (Claude Code, opencode, …) behind this boundary.
// See docs/design.md §5.
package harness

import "github.com/webgrip/ploeg/pkg/work"

// TaskSpec is injected into the agent Job (file mount or env).
type TaskSpec struct {
	WorkItem   work.WorkItem    `json:"workItem"`
	Role       string           `json:"role"` // specialist role within the team
	Checkpoint *work.Checkpoint `json:"checkpoint,omitempty"`
	RepoURL    string           `json:"repoUrl"`
	// Credentials are delivered out-of-band (mounted secrets), never here.
}

// OutcomeReport must be written by the container before exit.
// Exit-without-report is recorded as OutcomeFailed by the watcher.
type OutcomeReport struct {
	Outcome     work.Outcome     `json:"outcome"`
	Summary     string           `json:"summary"`
	Links       []string         `json:"links,omitempty"` // PRs, commits, created follow-ups
	Checkpoint  *work.Checkpoint `json:"checkpoint,omitempty"`
	StuckReason string           `json:"stuckReason,omitempty"` // mandatory when Outcome == stuck
}
