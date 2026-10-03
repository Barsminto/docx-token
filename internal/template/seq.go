package template

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var (
	paragraphXMLRE   = regexp.MustCompile(`(?s)<w:p\b[^>]*>.*?</w:p>`)
	paragraphTextRE  = regexp.MustCompile(`<w:t[^>]*>([^<]*)</w:t>`)
	leadingSectionRE = regexp.MustCompile(`^(\d+)\.`)
)

// FirstIssueSeq scans section numbers (e.g. "1. Executive Summary") before {BEGIN_RECORD}
// and returns the next number for the first issue in the loop.
func FirstIssueSeq(docxBytes []byte) (int, error) {
	xml, err := readDocumentXML(docxBytes)
	if err != nil {
		return 0, err
	}
	begin := strings.Index(xml, BeginRecordMarker)
	if begin < 0 {
		return 0, fmt.Errorf("%s not found in template", BeginRecordMarker)
	}

	maxNum := 0
	for _, p := range paragraphXMLRE.FindAllString(xml[:begin], -1) {
		text := paragraphPlainText(p)
		text = strings.TrimSpace(text)
		m := leadingSectionRE.FindStringSubmatch(text)
		if m == nil {
			continue
		}
		n, err := strconv.Atoi(m[1])
		if err != nil {
			continue
		}
		if n > maxNum {
			maxNum = n
		}
	}
	if maxNum == 0 {
		return 1, nil
	}
	return maxNum + 1, nil
}

func paragraphPlainText(paragraph string) string {
	var b strings.Builder
	for _, m := range paragraphTextRE.FindAllStringSubmatch(paragraph, -1) {
		b.WriteString(m[1])
	}
	return b.String()
}

func readDocumentXML(docxBytes []byte) (string, error) {
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
