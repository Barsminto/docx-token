package render

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/liboyang/docx-token/internal/domain"
	"github.com/liboyang/docx-token/internal/template"
)

const recordSubsectionNumIDBase = 100

var (
	recordHeadingRE            = regexp.MustCompile(`\{[a-zA-Z_][a-zA-Z0-9_]*(?:_\d+)?\.numberedTitle\}`)
	leadingChapterNumRE        = regexp.MustCompile(`^\d+\.\s*`)
	recordSubsectionLeadingRE  = regexp.MustCompile(`^\d+\.\d+\.\s*`)
	legacyRecordSubsecRE       = regexp.MustCompile(`^\{issue(?:_\d+)?\.sectionNumber\}\.\d+\s*`)
	numAbstractFromNumIDRE     = regexp.MustCompile(`<w:num w:numId="(\d+)"[^>]*>\s*<w:abstractNumId w:val="(\d+)"`)
	recordPrototypeListNumID   = 6
	execSummaryListNumID       = 3
)

// patchRecordSubsections: one Word numId per concern — ilvl 0 = "3.", ilvl 1 = 3.1, 3.2, …
func patchRecordSubsections(documentXML string, items []domain.LineItem) string {
	if len(items) == 0 {
		return documentXML
	}
	paragraphs := paragraphRE.FindAllString(documentXML, -1)
	if len(paragraphs) == 0 {
		return documentXML
	}

	itemIdx := -1
	currentNumID := recordSubsectionNumIDBase
	for i, p := range paragraphs {
		text := paragraphPlainText(p)
		if isInvestigationChapterTitle(p, text) {
			itemIdx++
			if itemIdx >= len(items) {
				break
			}
			currentNumID = recordSubsectionNumIDBase + itemIdx
			text = leadingChapterNumRE.ReplaceAllString(text, "")
			text = strings.TrimSpace(text)
			text = strings.ReplaceAll(text, template.GiBoldMarker, "")
			p = replaceParagraphTextPreservingRuns(p, text)
			p = strings.ReplaceAll(p, template.GiBoldMarker, "")
			paragraphs[i] = injectNumPrIntoParagraph(p, currentNumID, 0)
			continue
		}
		if itemIdx < 0 {
			continue
		}
		newP, ok := applyRecordSubsectionNumPr(p, text, currentNumID)
		if ok {
			paragraphs[i] = newP
		}
	}

	return replaceParagraphsInOrder(documentXML, paragraphs)
}

func isInvestigationChapterTitle(paragraph, plain string) bool {
	if strings.Contains(paragraph, template.DotMarker) || strings.Contains(paragraph, template.SubsecRecordMarker) {
		return false
	}
	if isConcernListLine(plain) {
		return false
	}
	if recordHeadingRE.MatchString(paragraph) {
		return true
	}
	if strings.Contains(plain, "INVESTIGATION REPORT") || strings.Contains(plain, "INVESTIGATION METHODOLOGY") {
		return false
	}
	return strings.Contains(plain, "INVESTIGATION:")
}

func applyRecordSubsectionNumPr(paragraph, plain string, numID int) (string, bool) {
	if !shouldApplyRecordSubsectionNumPr(paragraph, plain) {
		return paragraph, false
	}
	p := stripMarkersInParagraph(paragraph, template.DotMarker, template.SubsecRecordMarker)
	plain = paragraphPlainText(p)
	plain = legacyRecordSubsecRE.ReplaceAllString(plain, "")
	plain = recordSubsectionLeadingRE.ReplaceAllString(plain, "")
	plain = strings.TrimSpace(plain)
	if plain != strings.TrimSpace(paragraphPlainText(paragraph)) {
		p = replaceParagraphTextPreservingRuns(p, plain)
	}
	return injectNumPrIntoParagraph(p, numID, 1), true
}

func shouldApplyRecordSubsectionNumPr(paragraph, plain string) bool {
	if strings.Contains(paragraph, template.DotMarker) || strings.Contains(paragraph, template.SubsecRecordMarker) {
		return true
	}
	if legacyRecordSubsecRE.MatchString(plain) || recordSubsectionLeadingRE.MatchString(plain) {
		return true
	}
	// User templates often omit {.} but keep Word list ilvl 1 on record body paragraphs.
	if numPrBlockRE.MatchString(paragraph) && strings.Contains(paragraph, `<w:ilvl w:val="1"`) {
		return true
	}
	return false
}

// decimalMultiLevelAbstractNumID finds the abstract list used for "3." / "3.1" style numbering.
// Word-edited templates use abstractNumId 0 for decimal and 1 for a./b./c.; built-in templates invert that.
func decimalMultiLevelAbstractNumID(numberingXML string) int {
	for _, numID := range []int{recordPrototypeListNumID, execSummaryListNumID} {
		if id := abstractNumIDForNum(numberingXML, numID); id >= 0 {
			return id
		}
	}
	return 1
}

func letterListAbstractNumID(numberingXML string) int {
	for _, numID := range []int{1, 2, 9} {
		if id := abstractNumIDForNum(numberingXML, numID); id >= 0 {
			return id
		}
	}
	return 0
}

func abstractNumIDForNum(numberingXML string, numID int) int {
	needle := fmt.Sprintf(`<w:num w:numId="%d"`, numID)
	idx := strings.Index(numberingXML, needle)
	if idx < 0 {
		return -1
	}
	end := idx + 280
	if end > len(numberingXML) {
		end = len(numberingXML)
	}
	m := numAbstractFromNumIDRE.FindStringSubmatch(numberingXML[idx:end])
	if len(m) < 3 || m[1] != itoa(numID) {
		return -1
	}
	n, _ := strconv.Atoi(m[2])
	return n
}

func injectNumPrIntoParagraph(paragraph string, numID, ilvl int) string {
	numPr := fmt.Sprintf(`<w:numPr><w:ilvl w:val="%d"/><w:numId w:val="%d"/></w:numPr>`, ilvl, numID)
	if numPrBlockRE.MatchString(paragraph) {
		return numPrBlockRE.ReplaceAllString(paragraph, numPr)
	}
	if strings.Contains(paragraph, "<w:pPr>") {
		return strings.Replace(paragraph, "<w:pPr>", "<w:pPr>"+numPr, 1)
	}
	return strings.Replace(paragraph, "<w:p>", "<w:p><w:pPr>"+numPr+"</w:pPr>", 1)
}

func paragraphPlainText(paragraph string) string {
	var b strings.Builder
	for _, m := range paragraphTextRE.FindAllStringSubmatch(paragraph, -1) {
		b.WriteString(m[1])
	}
	return b.String()
}

var paragraphTextRE = regexp.MustCompile(`<w:t[^>]*>([^<]*)</w:t>`)

func setParagraphPlainText(paragraph, text string) string {
	if strings.TrimSpace(paragraphPlainText(paragraph)) == strings.TrimSpace(text) {
		return paragraph
	}
	if paragraphTextRE.MatchString(paragraph) {
		return replaceParagraphTextPreservingRuns(paragraph, text)
	}
	return rebuildParagraphPlain(paragraph, text)
}

func xmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	return s
}

func replaceParagraphsInOrder(xml string, paragraphs []string) string {
	indices := paragraphRE.FindAllStringIndex(xml, -1)
	if len(indices) != len(paragraphs) {
		return xml
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

func mergeRecordNumberingDefs(numberingXML string, items []domain.LineItem) string {
	if !strings.Contains(numberingXML, "<w:numbering") {
		return numberingXML
	}
	abstractID := decimalMultiLevelAbstractNumID(numberingXML)
	var extra strings.Builder
	for i, item := range items {
		chapter, _ := strconv.Atoi(item.SectionNumber)
		if chapter == 0 {
			chapter = 3
		}
		numID := recordSubsectionNumIDBase + i
		if strings.Contains(numberingXML, `w:numId="`+itoa(numID)+`"`) {
			continue
		}
		extra.WriteString(fmt.Sprintf(
			`<w:num w:numId="%d"><w:abstractNumId w:val="%d"/>`+
				`<w:lvlOverride w:ilvl="0"><w:startOverride w:val="%d"/></w:lvlOverride>`+
				`<w:lvlOverride w:ilvl="1"><w:startOverride w:val="1"/></w:lvlOverride></w:num>`,
			numID, abstractID, chapter))
	}
	if extra.Len() == 0 {
		return numberingXML
	}
	return strings.Replace(numberingXML, "</w:numbering>", extra.String()+"</w:numbering>", 1)
}

func itoa(n int) string {
	return fmt.Sprintf("%d", n)
}
