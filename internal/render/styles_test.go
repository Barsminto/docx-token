package render

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/liboyang/docx-token/internal/sample"
	"github.com/liboyang/docx-token/internal/template"
)

func readDocumentXML(docxPath string) (string, error) {
	raw, err := os.ReadFile(docxPath)
	if err != nil {
		return "", err
	}
	zr, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		return "", err
	}
	for _, f := range zr.File {
		if f.Name != "word/document.xml" {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return "", err
		}
		b, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	return "", io.EOF
}

func TestGenerateFillsIssueTokens(t *testing.T) {
	dir := t.TempDir()
	tpl := filepath.Join(dir, "report.docx")
	out := filepath.Join(dir, "out.docx")
	if err := template.WriteReportTemplate(tpl); err != nil {
		t.Fatal(err)
	}

	doc := sample.Document(time.Date(2026, 10, 3, 15, 0, 0, 0, time.UTC))
	if err := Generate(tpl, out, doc, nil); err != nil {
		t.Fatal(err)
	}

	docXML, err := readDocumentXML(out)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(docXML, "{Record_") || strings.Contains(docXML, "{BEGIN_RECORD}") {
		t.Fatal("placeholders should be replaced and loop markers removed")
	}
	if !strings.Contains(docXML, "INC-1001") {
		t.Fatal("expected issue title in output")
	}
}
