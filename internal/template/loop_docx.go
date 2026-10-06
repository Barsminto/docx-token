package template

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"strings"
)

// DiscoverLoopSpecsFromDocx reads word/document.xml and returns loop regions.
func DiscoverLoopSpecsFromDocx(docxBytes []byte) ([]LoopSpec, error) {
	xml, err := readDocumentXMLBytes(docxBytes)
	if err != nil {
		return nil, err
	}
	return DiscoverLoopSpecs(xml)
}

func readDocumentXMLBytes(docxBytes []byte) (string, error) {
	zr, err := zip.NewReader(bytes.NewReader(docxBytes), int64(len(docxBytes)))
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
		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
	return "", fmt.Errorf("word/document.xml not found")
}

// PrimaryRecordDataPath returns the slice path used for investigation chapter numbering.
func PrimaryRecordDataPath(specs []LoopSpec) string {
	for _, s := range specs {
		if s.End != "" && containsRecordRegionTag(s.Begin) {
			return s.DataPath
		}
	}
	if len(specs) > 0 {
		return specs[len(specs)-1].DataPath
	}
	return "case.issues"
}

func containsRecordRegionTag(begin string) bool {
	return strings.Contains(begin, ":record")
}
