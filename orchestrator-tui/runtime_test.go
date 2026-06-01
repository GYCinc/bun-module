package main

import (
	"testing"
)

func TestRoleBlueprint(t *testing.T) {
	bp := NewRoleBlueprint("researcher", 3)
	specs := bp.BuildSpecs("test task")

	if len(specs) != 3 {
		t.Errorf("Expected 3 specs, got %d", len(specs))
	}

	if specs[0].Name != "context" {
		t.Errorf("Expected first spec name to be 'context', got '%s'", specs[0].Name)
	}

	bp2 := NewRoleBlueprint("custom", 10) // should cap at 5
	specs2 := bp2.BuildSpecs("test task")

	if len(specs2) != 5 {
		t.Errorf("Expected max 5 specs, got %d", len(specs2))
	}
}
