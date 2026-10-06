package docxtoken

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

func TestFillHeaderPlaceholders(t *testing.T) {
	dir := t.TempDir()
	tpl := filepath.Join(dir, "template.docx")
	out := filepath.Join(dir, "out.docx")
	if err := template.WriteReportTemplate(tpl); err != nil {
		t.Fatal(err)
	}
	data := sample.FillMap(time.Date(2026, 2, 1, 9, 0, 0, 0, time.UTC))
	if err := Fill(tpl, out, data); err != nil {
		t.Fatal(err)
	}
	hdr, err := readZipPart(out, "word/header1.xml")
	if err != nil {
		t.Fatal(err)
	}
	plain := strings.ReplaceAll(hdr, " ", "")
	if strings.Contains(plain, "{projectName}") {
		t.Fatalf("header projectName not replaced: %s", hdr)
	}
	if !strings.Contains(hdr, "XXX") || !strings.Contains(hdr, "Sample investigation case") {
		t.Fatalf("expected header project/case text: %s", hdr)
	}
}

func readZipPart(docxPath, name string) (string, error) {
	raw, err := os.ReadFile(docxPath)
	if err != nil {
		return "", err
	}
	zr, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		return "", err
	}
	for _, f := range zr.File {
		if f.Name != name {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return "", err
		}
		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
	return "", os.ErrNotExist
}
