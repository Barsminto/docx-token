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

	m := map[string]string{}
	BindLoopPrefix(m, "issue", []map[string]interface{}{
		{"brief": "Brief", "title": "CONCERN 1"},
	}, doc.Items)
	if m[Key("issue_0", "brief")] != "Brief" {
		t.Fatalf("brief: %q", m[Key("issue_0", "brief")])
	}
	if m[Key("issue_0", "numberedTitle")] != "3. INVESTIGATION: CONCERN 1" {
		t.Fatalf("numberedTitle: %q", m[Key("issue_0", "numberedTitle")])
	}
	if m[Key("issue_0", "listIndex")] != "1" {
		t.Fatalf("listIndex: %q", m[Key("issue_0", "listIndex")])
	}
}
