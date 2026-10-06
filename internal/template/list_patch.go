package template

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"regexp"
	"strings"
)

const numberingRID = "rId10"

var (
	paragraphBlockRE  = regexp.MustCompile(`(?s)(<w:p\b[^>]*>)(.*?)(</w:p>)`)
	concernListLineRE = regexp.MustCompile(`\{[a-zA-Z_][a-zA-Z0-9_]*\.listIndex\}.*\{[a-zA-Z_][a-zA-Z0-9_]*\.brief\}`)
)

// NumberingAttachOptions controls whether structural markers stay visible in the saved template.
type NumberingAttachOptions struct {
	PreserveMarkers bool
}

// attachDocumentNumbering adds Word list styles: a./b./c. concerns and decimal 1.1 / 2.1 subsections.
func attachDocumentNumbering(docxPath string, opts NumberingAttachOptions) error {
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

	doc, ok := files["word/document.xml"]
	if !ok {
		return nil
	}
	patched := patchConcernListParagraphs(string(doc), BeginIssueListMarker, EndIssueListMarker, 1)
	patched = patchConcernListParagraphs(patched, BeginBackgroundListMarker, EndBackgroundListMarker, 2)
	patched = patchSubsectionParagraphs(patched, opts.PreserveMarkers)
	patched = patchRecordSubsectionParagraphs(patched, 6, 3)
	patched = patchClosingParagraphs(patched, opts.PreserveMarkers)
	if !opts.PreserveMarkers {
		patched = stripGiBoldMarkers(patched)
	}
	files["word/document.xml"] = []byte(patched)

	files["word/numbering.xml"] = []byte(numberingXML)
	if !contains(order, "word/numbering.xml") {
		order = append(order, "word/numbering.xml")
	}

	rels := string(files["word/_rels/document.xml.rels"])
	if !strings.Contains(rels, "numbering.xml") {
		insert := `<Relationship Id="` + numberingRID + `" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/numbering" Target="numbering.xml"></Relationship>`
		rels = strings.Replace(rels, "</Relationships>", insert+"</Relationships>", 1)
		files["word/_rels/document.xml.rels"] = []byte(rels)
	}

	ct := string(files["[Content_Types].xml"])
	if !strings.Contains(ct, "/word/numbering") {
		ct = contentTypeOverrideRE.ReplaceAllString(ct,
			`<Override PartName="/word/numbering.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.numbering+xml"/>$1`)
		files["[Content_Types].xml"] = []byte(ct)
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

// patchRecordSubsectionParagraphs applies decimal x.1 numbering inside {BEGIN_RECORD} prototype (template only).
func patchRecordSubsectionParagraphs(xml string, numID, chapterStart int) string {
	beginIdx := strings.Index(xml, BeginRecordMarker)
	endIdx := strings.Index(xml, EndRecordMarker)
	if beginIdx < 0 || endIdx < 0 || endIdx <= beginIdx {
		return xml
	}
	head := xml[:beginIdx]
	mid := xml[beginIdx:endIdx]
	tail := xml[endIdx:]
	mid = paragraphBlockRE.ReplaceAllStringFunc(mid, func(p string) string {
		if strings.Contains(p, SubsecRecordMarker) {
			return injectNumPr(p, numID, 1)
		}
		if strings.Contains(p, GiBoldMarker) && strings.Contains(p, "INVESTIGATION:") {
			return injectNumPr(p, numID, 0)
		}
		return p
	})
	_ = chapterStart
	return head + mid + tail
}

func patchClosingParagraphs(xml string, keepMarkers bool) string {
	return paragraphBlockRE.ReplaceAllStringFunc(xml, func(p string) string {
		switch {
		case strings.Contains(p, GiBoldMarker) && strings.Contains(p, "CONCLUSION"):
			return injectNumPr(p, 7, 0)
		case strings.Contains(p, SubsecConclusionMarker):
			if !keepMarkers {
				p = strings.ReplaceAll(p, SubsecConclusionMarker, "")
			}
			return injectNumPr(p, 7, 1)
		case strings.Contains(p, SubsecLetterMarker):
			if !keepMarkers {
				p = strings.ReplaceAll(p, SubsecLetterMarker, "")
			}
			return injectNumPr(p, 9, 0)
		case strings.Contains(p, GiBoldMarker) && strings.Contains(p, "RECOMMENDATIONS"):
			if !keepMarkers {
				p = strings.ReplaceAll(p, GiBoldMarker, "")
			}
			return ensureRunBold(injectNumPr(p, 8, 0))
		case strings.Contains(p, SubsecRecommendMarker):
			if !keepMarkers {
				p = strings.ReplaceAll(p, SubsecRecommendMarker, "")
			}
			return injectNumPr(p, 8, 1)
		default:
			return p
		}
	})
}

func patchSubsectionParagraphs(xml string, keepMarkers bool) string {
	paragraphs := paragraphBlockRE.FindAllString(xml, -1)
	if len(paragraphs) == 0 {
		return xml
	}
	indices := paragraphBlockRE.FindAllStringIndex(xml, -1)
	if len(indices) != len(paragraphs) {
		return xml
	}
	section := "exec"
	for i, p := range paragraphs {
		plain := strings.TrimSpace(stripXMLText(p))
		if strings.Contains(plain, "1. BACKGROUND") {
			section = "bg"
		}
		if strings.Contains(plain, "2. INVESTIGATION METHODOLOGY") {
			section = "meth"
		}
		if !strings.Contains(p, DotMarker) {
			continue
		}
		if !keepMarkers {
			p = strings.ReplaceAll(p, DotMarker, "")
		}
		switch section {
		case "bg":
			paragraphs[i] = injectNumPr(p, 4, 1)
		case "meth":
			paragraphs[i] = injectNumPr(p, 5, 1)
		default:
			paragraphs[i] = injectNumPr(p, 3, 1)
		}
	}
	var b strings.Builder
	last := 0
	for i, idx := range indices {
		b.WriteString(xml[last:idx[0]])
		b.WriteString(paragraphs[i])
		last = idx[1]
	}
	b.WriteString(xml[last:])
	return b.String()
}

func stripXMLText(paragraph string) string {
	return regexp.MustCompile(`<[^>]+>`).ReplaceAllString(paragraph, "")
}

func patchConcernListParagraphs(xml, begin, end string, numID int) string {
	beginIdx := strings.Index(xml, begin)
	endIdx := strings.Index(xml, end)
	if beginIdx < 0 || endIdx < 0 || endIdx <= beginIdx {
		return xml
	}
	head := xml[:beginIdx]
	mid := xml[beginIdx:endIdx]
	tail := xml[endIdx:]

	out := paragraphBlockRE.ReplaceAllStringFunc(mid, func(p string) string {
		if !concernListLineRE.MatchString(p) {
			return p
		}
		return injectNumPr(p, numID, 0)
	})
	return head + out + tail
}

func injectNumPr(paragraph string, numID, ilvl int) string {
	numPr := `<w:numPr><w:ilvl w:val="` + itoa(ilvl) + `"/><w:numId w:val="` + itoa(numID) + `"/></w:numPr>`
	if strings.Contains(paragraph, "<w:pPr>") {
		return strings.Replace(paragraph, "<w:pPr>", "<w:pPr>"+numPr, 1)
	}
	return strings.Replace(paragraph, "<w:p>", "<w:p><w:pPr>"+numPr+"</w:pPr>", 1)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [12]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}

const numberingXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:numbering xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:abstractNum w:abstractNumId="0">
    <w:lvl w:ilvl="0">
      <w:start w:val="1"/>
      <w:numFmt w:val="lowerLetter"/>
      <w:lvlText w:val="%1."/>
      <w:lvlJc w:val="left"/>
      <w:pPr><w:ind w:left="720" w:hanging="360"/></w:pPr>
    </w:lvl>
  </w:abstractNum>
  <w:abstractNum w:abstractNumId="1">
    <w:lvl w:ilvl="0">
      <w:start w:val="1"/>
      <w:numFmt w:val="decimal"/>
      <w:lvlText w:val="%1."/>
      <w:lvlJc w:val="left"/>
      <w:pPr><w:ind w:left="0" w:hanging="0"/></w:pPr>
    </w:lvl>
    <w:lvl w:ilvl="1">
      <w:start w:val="1"/>
      <w:numFmt w:val="decimal"/>
      <w:lvlText w:val="%1.%2."/>
      <w:lvlJc w:val="left"/>
      <w:pPr><w:ind w:left="720" w:hanging="360"/></w:pPr>
    </w:lvl>
  </w:abstractNum>
  <w:num w:numId="1"><w:abstractNumId w:val="0"/></w:num>
  <w:num w:numId="2">
    <w:abstractNumId w:val="0"/>
    <w:lvlOverride w:ilvl="0"><w:startOverride w:val="1"/></w:lvlOverride>
  </w:num>
  <w:num w:numId="3">
    <w:abstractNumId w:val="1"/>
    <w:lvlOverride w:ilvl="0"><w:startOverride w:val="1"/></w:lvlOverride>
  </w:num>
  <w:num w:numId="4">
    <w:abstractNumId w:val="1"/>
    <w:lvlOverride w:ilvl="0"><w:startOverride w:val="1"/></w:lvlOverride>
    <w:lvlOverride w:ilvl="1"><w:startOverride w:val="1"/></w:lvlOverride>
  </w:num>
  <w:num w:numId="5">
    <w:abstractNumId w:val="1"/>
    <w:lvlOverride w:ilvl="0"><w:startOverride w:val="2"/></w:lvlOverride>
    <w:lvlOverride w:ilvl="1"><w:startOverride w:val="1"/></w:lvlOverride>
  </w:num>
  <w:num w:numId="6">
    <w:abstractNumId w:val="1"/>
    <w:lvlOverride w:ilvl="0"><w:startOverride w:val="3"/></w:lvlOverride>
    <w:lvlOverride w:ilvl="1"><w:startOverride w:val="1"/></w:lvlOverride>
  </w:num>
  <w:num w:numId="7">
    <w:abstractNumId w:val="1"/>
    <w:lvlOverride w:ilvl="0"><w:startOverride w:val="8"/></w:lvlOverride>
    <w:lvlOverride w:ilvl="1"><w:startOverride w:val="1"/></w:lvlOverride>
  </w:num>
  <w:num w:numId="8">
    <w:abstractNumId w:val="1"/>
    <w:lvlOverride w:ilvl="0"><w:startOverride w:val="9"/></w:lvlOverride>
    <w:lvlOverride w:ilvl="1"><w:startOverride w:val="1"/></w:lvlOverride>
  </w:num>
  <w:num w:numId="9">
    <w:abstractNumId w:val="0"/>
    <w:lvlOverride w:ilvl="0"><w:startOverride w:val="1"/></w:lvlOverride>
  </w:num>
</w:numbering>`
