package domain

import (
	"testing"
	"time"
)

func TestDocumentPlaceholdersObjectKeys(t *testing.T) {
	doc := NewDocument("PROJECT XXX", "Inv", time.Date(2026, 2, 1, 9, 0, 0, 0, time.UTC), []LineItemInput{
		{Title: "CONCERN 1", ConcernBrief: "Brief", FindingOutcome: "not substantiated"},
	})
	doc.AssignIssueSeq(3)

	m := doc.Placeholders()
	if m[Key("issue_0", "brief")] != "Brief" {
		t.Fatalf("brief: %q", m[Key("issue_0", "brief")])
	}
	if m[Key("issue_0", "numberedTitle")] != "3. INVESTIGATION: CONCERN 1" {
		t.Fatalf("numberedTitle: %q", m[Key("issue_0", "numberedTitle")])
	}
}
