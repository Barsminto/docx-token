package docxtoken

import (
	"testing"
	"time"
)

func TestBuildFromMapFlatAndRecords(t *testing.T) {
	data := map[string]interface{}{
		"Report.Title":        "Q1",
		"Summary.Author":      "Alice",
		"Report.GeneratedAt":  "2026-01-02 09:00",
		RecordsKey: []map[string]interface{}{
			{"Title": "INC-1", "Category": "Hardware"},
		},
	}
	doc, err := BuildFromMap(data)
	if err != nil {
		t.Fatal(err)
	}
	if doc.Report.Title != "Q1" { // Report.Title when Project.Name absent
		t.Fatalf("title %q", doc.Report.Title)
	}
	if len(doc.Items) != 1 || doc.Items[0].Category != "Hardware" {
		t.Fatalf("items %+v", doc.Items)
	}
	doc.AssignIssueSeq(2)
	if doc.Items[0].NumberedTitle != "2. INVESTIGATION: INC-1" {
		t.Fatalf("numbered %q", doc.Items[0].NumberedTitle)
	}
}

func TestBuildFromMapNested(t *testing.T) {
	data := map[string]interface{}{
		"Report": map[string]interface{}{
			"Title": "Nested",
		},
		RecordsKey: []interface{}{
			map[string]interface{}{"Title": "A", "Category": "C"},
		},
	}
	doc, err := BuildFromMap(data)
	if err != nil {
		t.Fatal(err)
	}
	if doc.Report.Title != "Nested" {
		t.Fatalf("title %q", doc.Report.Title)
	}
}

func TestBuildFromMapDefaultTime(t *testing.T) {
	doc, err := BuildFromMap(map[string]interface{}{RecordsKey: []map[string]interface{}{
		{"Title": "X", "Category": "Y"},
	}})
	if err != nil {
		t.Fatal(err)
	}
	if doc.Report.GeneratedAt == "" {
		t.Fatal("expected generated at")
	}
	_, err = time.Parse("2006-01-02 15:04", doc.Report.GeneratedAt)
	if err != nil {
		t.Fatal(err)
	}
}
