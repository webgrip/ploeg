package httpapi

import (
	"log/slog"
	"testing"

	"github.com/webgrip/ploeg/pkg/provider"
	"github.com/webgrip/ploeg/pkg/target"
	"github.com/webgrip/ploeg/pkg/work"
)

func TestPinTeam_ContainerPinBeatsAssigneeMapping(t *testing.T) {
	s := &Server{
		ScopeTeams: map[string]string{"901525751877": "app"},
		Log:        slog.New(slog.DiscardHandler),
	}
	item := work.WorkItem{Team: "docs", ExternalScope: "901525751877"}
	s.pinTeam(&item)
	if item.Team != "app" {
		t.Errorf("team = %q, want the container's pinned team to win", item.Team)
	}
}

func TestPinTeam_UnpinnedContainerLeavesAssigneeDecision(t *testing.T) {
	s := &Server{
		ScopeTeams: map[string]string{"901525751877": "app"},
		Log:        slog.New(slog.DiscardHandler),
	}
	item := work.WorkItem{Team: "infra", ExternalScope: "somewhere-else"}
	s.pinTeam(&item)
	if item.Team != "infra" {
		t.Errorf("team = %q, want the assignee's team untouched", item.Team)
	}
}

func TestPinTeam_NoScopeNoOpinion(t *testing.T) {
	s := &Server{
		ScopeTeams: map[string]string{"901525751877": "app"},
		Log:        slog.New(slog.DiscardHandler),
	}
	item := work.WorkItem{Team: "docs"}
	s.pinTeam(&item)
	if item.Team != "docs" {
		t.Errorf("team = %q, want unchanged when the scope is unknown", item.Team)
	}
}

func TestPinTeam_RunsBeforeTargetResolution(t *testing.T) {
	r, err := target.NewMapResolver("901525751877/app=code14nl/internal/poc-silk@main", "gitlab")
	if err != nil {
		t.Fatal(err)
	}
	s := &Server{
		Targets:    r,
		ScopeTeams: map[string]string{"901525751877": "app"},
		Log:        slog.New(slog.DiscardHandler),
	}
	item := work.WorkItem{Team: "docs", ExternalScope: "901525751877"}
	s.pinTeam(&item)
	s.resolveTarget(&item, provider.TrackerEvent{ExternalID: "9001"})
	if item.Target == nil || item.Target.Repo != "internal/poc-silk" {
		t.Fatalf("target = %+v, want the team-qualified rule matched via the pinned team", item.Target)
	}
	if item.RouteRule != "901525751877/app" {
		t.Errorf("rule = %q", item.RouteRule)
	}
}
