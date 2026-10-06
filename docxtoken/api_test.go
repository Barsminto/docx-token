package docxtoken

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/liboyang/docx-token/internal/sample"
	"github.com/liboyang/docx-token/internal/template"
)

func TestGenerateRequiresPaths(t *testing.T) {
	if err := Generate(Options{Data: map[string]interface{}{}}); err == nil {
		t.Fatal("expected error for missing paths")
	}
}

func TestGenerateWritesDocx(t *testing.T) {
	dir := t.TempDir()
	tpl := filepath.Join(dir, "template.docx")
	out := filepath.Join(dir, "out.docx")
	if err := template.WriteReportTemplate(tpl); err != nil {
		t.Fatal(err)
	}
	err := Generate(Options{
		TemplatePath: tpl,
		OutputPath:   out,
		Data:         sample.FillMap(time.Date(2026, 2, 1, 9, 0, 0, 0, time.UTC)),
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(out); err != nil {
		t.Fatal(err)
	}
}
