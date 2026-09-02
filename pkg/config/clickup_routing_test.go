package config

import (
	"context"
	"strings"
	"testing"
)

func TestTargetSpec_ClickupListsRouteByPinnedID(t *testing.T) {
	f, err := Load(write(t, `
trackers:
  clickup:
    projects:
      - name: "Ploeg App"
        id: "901525751877"
        repo: code14nl/internal/poc-silk
        branch: main
        forge: gitlab
        team: app
      - name: "Ploeg Infra"
        id: "901525751881"
        repo: code14nl/internal/poc-silk
        branch: main
        forge: gitlab
        team: infra
teams:
  app:
    assignees: []
  infra:
    assignees: []
`))
	if err != nil {
		t.Fatal(err)
	}
	spec, err := f.TargetSpec(context.Background(), nil, discard())
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"901525751877/app=code14nl/internal/poc-silk@main;forge=gitlab",
		"901525751881/infra=code14nl/internal/poc-silk@main;forge=gitlab",
	} {
		if !strings.Contains(spec, want) {
			t.Errorf("spec = %q, want %q", spec, want)
		}
	}
}

func TestValidate_ClickupProjectNeedsAPinnedID(t *testing.T) {
	_, err := Load(write(t, `
trackers:
  clickup:
    projects:
      - name: "Ploeg App"
        repo: code14nl/internal/poc-silk
`))
	if err == nil || !strings.Contains(err.Error(), "trackers.clickup.projects[0]") {
		t.Fatalf("err = %v, want it to name the clickup entry that lacks an id", err)
	}
}

func TestValidate_ScopeIDsShareOneNamespaceAcrossTrackers(t *testing.T) {
	_, err := Load(write(t, `
trackers:
  vikunja:
    projects:
      - id: "7"
        repo: webgrip/ploeg
  clickup:
    projects:
      - id: "7"
        repo: code14nl/internal/poc-silk
`))
	if err == nil || !strings.Contains(err.Error(), "routed twice") {
		t.Fatalf("err = %v, want the cross-tracker id collision refused at boot", err)
	}
}

func TestTargetSpec_MixedTrackersJoinOneSpec(t *testing.T) {
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
teams:
  docs:
    assignees: []
`))
	if err != nil {
		t.Fatal(err)
	}
	spec, err := f.TargetSpec(context.Background(), nil, discard())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(spec, "2=webgrip/ploeg") || !strings.Contains(spec, "901525751875/docs=code14nl/internal/poc-silk") {
		t.Errorf("spec = %q, want entries from both trackers", spec)
	}
}
