package template

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFirstIssueSeqFromBuiltTemplate(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "t.docx")
	if err := WriteReportTemplate(path); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	seq, err := FirstIssueSeq(raw)
	if err != nil {
		t.Fatal(err)
	}
	// Sections 1. BACKGROUND and 2. INVESTIGATION METHODOLOGY precede {BEGIN_RECORD}.
	if seq != 3 {
		t.Fatalf("expected first concern section seq 3, got %d", seq)
	}
}
