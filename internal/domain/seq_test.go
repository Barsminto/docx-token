package domain

import (
	"testing"
	"time"
)

func TestAssignIssueSeqListLetterAndConclusionOffset(t *testing.T) {
	doc := NewDocument("P", "A", mustTime(), []LineItemInput{
		{ConcernBrief: "b1", ConcernSummary: "s1"},
		{ConcernBrief: "b2"},
		{ConcernBrief: "b3"},
	})
	doc.AssignIssueSeq(3)

	if doc.Items[0].ListLetter != "a" || doc.Items[2].ListLetter != "c" {
		t.Fatalf("letters %+v", doc.Items)
	}
	if doc.Items[0].NumberedTitle != "3. INVESTIGATION: CONCERN 1" {
		t.Fatalf("title %q", doc.Items[0].NumberedTitle)
	}
	if doc.Items[1].StmtLabelA != "c" || doc.Items[1].StmtLabelB != "d" {
		t.Fatalf("stmt labels %s %s", doc.Items[1].StmtLabelA, doc.Items[1].StmtLabelB)
	}
	conclusionNum := 3 + len(doc.Items)
	if conclusionNum != 6 {
		t.Fatalf("expected conclusion section 6, got %d", conclusionNum)
	}
}

func mustTime() time.Time {
	return time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
}
