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
	if seq != 2 {
		t.Fatalf("expected first issue seq 2, got %d", seq)
	}
}
