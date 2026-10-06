package template

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"regexp"
	"strings"
)

const headerRID = "rId9"

var contentTypeOverrideRE = regexp.MustCompile(`(</Types>)`)

func attachHeader(docxPath string, headerBody string) error {
	raw, err := os.ReadFile(docxPath)
	if err != nil {
		return err
	}
	zr, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		return err
	}

	files := map[string][]byte{}
	order := make([]string, 0, len(zr.File)+1)
	for _, f := range zr.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return err
		}
		files[f.Name] = data
		order = append(order, f.Name)
	}

	hdr := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<w:hdr xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" ` +
		`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">` +
		headerBody + `</w:hdr>`

	files["word/header1.xml"] = []byte(hdr)
	if !contains(order, "word/header1.xml") {
		order = append(order, "word/header1.xml")
	}

	rels := string(files["word/_rels/document.xml.rels"])
	if !strings.Contains(rels, "header1.xml") {
		insert := `<Relationship Id="` + headerRID + `" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/header" Target="header1.xml"></Relationship>`
		rels = strings.Replace(rels, "</Relationships>", insert+"</Relationships>", 1)
		files["word/_rels/document.xml.rels"] = []byte(rels)
	}

	ct := string(files["[Content_Types].xml"])
	if !strings.Contains(ct, "/word/header") {
		ct = contentTypeOverrideRE.ReplaceAllString(ct,
			`<Override PartName="/word/header1.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.header+xml"/>$1`)
		files["[Content_Types].xml"] = []byte(ct)
	}

	doc := string(files["word/document.xml"])
	ref := `<w:headerReference w:type="default" r:id="` + headerRID + `"></w:headerReference>`
	if !strings.Contains(doc, ref) {
		doc = strings.Replace(doc, "<w:sectPr>", "<w:sectPr>"+ref, 1)
		files["word/document.xml"] = []byte(doc)
	}

	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)
	for _, name := range order {
		w, err := zw.Create(name)
		if err != nil {
			return err
		}
		if _, err := w.Write(files[name]); err != nil {
			return err
		}
	}
	if err := zw.Close(); err != nil {
		return err
	}
	return os.WriteFile(docxPath, buf.Bytes(), 0o644)
}

// attachDefaultHeader adds word/header1.xml with Report tokens and links it in sectPr.
func attachDefaultHeader(docxPath string) error {
	body := `<w:p><w:r><w:rPr><w:color w:val="FF0000"/></w:rPr><w:t>CONFIDENTIAL</w:t></w:r></w:p>` +
		`<w:p><w:r><w:t>{projectName} | {case.name} | DRAFT: {draftDate}</w:t></w:r></w:p>`
	return attachHeader(docxPath, body)
}

// attachReferenceHeader includes all common header placeholders for the reference template.
func attachReferenceHeader(docxPath string) error {
	body := `<w:p><w:r><w:rPr><w:color w:val="FF0000"/></w:rPr><w:t>CONFIDENTIAL</w:t></w:r></w:p>` +
		`<w:p><w:r><w:t>{projectName} | {case.name} | DRAFT: {draftDate}</w:t></w:r></w:p>` +
		`<w:p><w:r><w:t>{Report.Title} | {Report.Date} | generated {Report.GeneratedAt}</w:t></w:r></w:p>`
	return attachHeader(docxPath, body)
}

func contains(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}
