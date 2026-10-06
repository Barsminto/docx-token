package template

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteReferenceTemplateKeepsMarkers(t *testing.T) {
	path := filepath.Join(t.TempDir(), "reference-template.docx")
	if err := WriteReferenceTemplate(path); err != nil {
		t.Fatal(err)
	}
	docXML, err := readZipPart(path, "word/document.xml")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{DotMarker, TildeMarker, BeginCaseIssuesRecordMarker, BeginCaseIssuesMarker} {
		if !strings.Contains(docXML, want) {
			t.Fatalf("reference template should contain %q", want)
		}
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
