package render

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"

	"github.com/liboyang/docx-token/internal/domain"
)

func patchDocumentXML(docxBytes []byte, items []domain.LineItem) ([]byte, error) {
	zr, err := zip.NewReader(bytes.NewReader(docxBytes), int64(len(docxBytes)))
	if err != nil {
		return nil, err
	}

	fileOrder := make([]string, 0, len(zr.File))
	files := make(map[string][]byte, len(zr.File))
	for _, f := range zr.File {
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return nil, err
		}
		fileOrder = append(fileOrder, f.Name)
		files[f.Name] = data
	}

	raw, ok := files["word/document.xml"]
	if !ok {
		return nil, fmt.Errorf("word/document.xml not found")
	}
	files["word/document.xml"] = []byte(patchIssueHeadings(string(raw), items))

	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)
	for _, name := range fileOrder {
		w, err := zw.Create(name)
		if err != nil {
			return nil, err
		}
		if _, err := w.Write(files[name]); err != nil {
			return nil, err
		}
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
