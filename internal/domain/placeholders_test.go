package domain

import (
	"testing"
	"time"
)

func TestDocumentPlaceholdersObjectKeys(t *testing.T) {
	doc := NewDocument("Q1 Report", "Jane Doe", time.Date(2026, 1, 2, 9, 0, 0, 0, time.UTC), []LineItemInput{
		{Title: "INC-1", Category: "Hardware"},
	})
	doc.AssignIssueSeq(2)

	m := doc.Placeholders()
	if m[Key("Record_0", "Category")] != "Hardware" {
		t.Fatalf("Record_0.Category: %q", m[Key("Record_0", "Category")])
	}
	if m[Key("Record_0", "Seq")] != "2" {
		t.Fatalf("Record_0.Seq: %q", m[Key("Record_0", "Seq")])
	}
}
