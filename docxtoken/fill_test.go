package docxtoken

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/liboyang/docx-token/internal/sample"
	"github.com/liboyang/docx-token/internal/template"
)

func TestFillGITemplate(t *testing.T) {
	dir := t.TempDir()
	tpl := filepath.Join(dir, "template.docx")
	out := filepath.Join(dir, "out.docx")
	if err := template.WriteReportTemplate(tpl); err != nil {
		t.Fatal(err)
	}
	if err := Fill(tpl, out, sample.FillMap(time.Date(2026, 2, 1, 9, 0, 0, 0, time.UTC))); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), "{BEGIN_RECORD}") {
		t.Fatal("loop markers should be removed")
	}
}
