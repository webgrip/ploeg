package config

import "testing"

func TestScopeTeams_PinnedProjectsOnly(t *testing.T) {
	f, err := Load(write(t, `
trackers:
  vikunja:
    projects:
      - id: "2"
        repo: webgrip/ploeg
  clickup:
    projects:
      - id: "901525751875"
        repo: code14nl/internal/poc-silk
        team: docs
      - id: "901525751881"
        repo: code14nl/internal/poc-silk
        team: infra
teams:
  docs:
    assignees: []
  infra:
    assignees: []
`))
	if err != nil {
		t.Fatal(err)
	}
	pins := f.ScopeTeams()
	if len(pins) != 2 {
		t.Fatalf("pins = %v, want exactly the two pinned containers", pins)
	}
	if pins["901525751875"] != "docs" || pins["901525751881"] != "infra" {
		t.Errorf("pins = %v", pins)
	}
}
