package render

import (
	"fmt"
	"strings"

	"github.com/liboyang/docx-token/internal/template"
)

const (
	closingConclusionNumID      = 200
	closingRecommendationsNumID = 201
	closingLetterNumID          = 202
)

// patchClosingSections applies Word numbering to CONCLUSION / RECOMMENDATIONS (8., 8.1, 9., 9.1…).
func patchClosingSections(documentXML string, conclusionNum, recommendationsNum int) string {
	paragraphs := paragraphRE.FindAllString(documentXML, -1)
	if len(paragraphs) == 0 {
		return documentXML
	}

	phase := ""
	for i, p := range paragraphs {
		plain := paragraphPlainText(p)
		if isClosingChapterTitle(p, plain, "RECOMMENDATIONS") {
			phase = "recommendations"
			p = strings.ReplaceAll(p, template.GiBoldMarker, "")
			p = setParagraphPlainText(p, "RECOMMENDATIONS")
			paragraphs[i] = ensureRunBold(injectNumPrIntoParagraph(p, closingRecommendationsNumID, 0))
			continue
		}
		if isClosingChapterTitle(p, plain, "CONCLUSION") {
			phase = "conclusion"
			p = strings.ReplaceAll(p, template.GiBoldMarker, "")
			p = setParagraphPlainText(p, "CONCLUSION")
			paragraphs[i] = ensureRunBold(injectNumPrIntoParagraph(p, closingConclusionNumID, 0))
			continue
		}
		if strings.Contains(p, template.TildeMarker) && phase == "conclusion" {
			p = stripMarkersInParagraph(p, template.TildeMarker)
			p = trimParagraphTextRuns(p)
			paragraphs[i] = injectNumPrIntoParagraph(p, closingLetterNumID, 0)
			continue
		}
		if strings.Contains(p, template.DotMarker) && phase == "recommendations" {
			p = stripMarkersInParagraph(p, template.DotMarker)
			p = trimParagraphTextRuns(p)
			paragraphs[i] = injectNumPrIntoParagraph(p, closingRecommendationsNumID, 1)
			continue
		}
		if strings.Contains(p, template.DotMarker) && phase == "conclusion" {
			p = stripMarkersInParagraph(p, template.DotMarker)
			p = trimParagraphTextRuns(p)
			paragraphs[i] = injectNumPrIntoParagraph(p, closingConclusionNumID, 1)
		}
	}

	out := replaceParagraphsInOrder(documentXML, paragraphs)
	out = strings.ReplaceAll(out, `w:numId w:val="7"`, `w:numId w:val="200"`)
	out = strings.ReplaceAll(out, `w:numId w:val="8"`, `w:numId w:val="201"`)
	out = strings.ReplaceAll(out, `w:numId w:val="9"`, `w:numId w:val="202"`)
	return out
}

func ensureRunBold(paragraph string) string {
	if strings.Contains(paragraph, "<w:b") {
		return paragraph
	}
	if strings.Contains(paragraph, "<w:rPr>") {
		return strings.Replace(paragraph, "<w:rPr>", "<w:rPr><w:b/>", 1)
	}
	return strings.Replace(paragraph, "<w:r>", "<w:r><w:rPr><w:b/></w:rPr>", 1)
}

func isClosingChapterTitle(paragraph, plain, name string) bool {
	plain = strings.TrimSpace(strings.ReplaceAll(plain, template.GiBoldMarker, ""))
	plain = leadingChapterNumRE.ReplaceAllString(plain, "")
	if strings.TrimSpace(plain) != name {
		return false
	}
	return strings.Contains(paragraph, template.GiBoldMarker) || strings.Contains(paragraph, "<w:b")
}

func mergeClosingNumberingDefs(numberingXML string, conclusionNum, recommendationsNum int) string {
	if !strings.Contains(numberingXML, "<w:numbering") {
		return numberingXML
	}
	if conclusionNum == 0 {
		conclusionNum = 8
	}
	if recommendationsNum == 0 {
		recommendationsNum = 9
	}
	decimalAbstract := decimalMultiLevelAbstractNumID(numberingXML)
	letterAbstract := letterListAbstractNumID(numberingXML)
	extra := fmt.Sprintf(
		`<w:num w:numId="%d"><w:abstractNumId w:val="%d"/>`+
			`<w:lvlOverride w:ilvl="0"><w:startOverride w:val="%d"/></w:lvlOverride>`+
			`<w:lvlOverride w:ilvl="1"><w:startOverride w:val="1"/></w:lvlOverride></w:num>`+
			`<w:num w:numId="%d"><w:abstractNumId w:val="%d"/>`+
			`<w:lvlOverride w:ilvl="0"><w:startOverride w:val="%d"/></w:lvlOverride>`+
			`<w:lvlOverride w:ilvl="1"><w:startOverride w:val="1"/></w:lvlOverride></w:num>`+
			`<w:num w:numId="%d"><w:abstractNumId w:val="%d"/>`+
			`<w:lvlOverride w:ilvl="0"><w:startOverride w:val="1"/></w:lvlOverride></w:num>`,
		closingConclusionNumID, decimalAbstract, conclusionNum,
		closingRecommendationsNumID, decimalAbstract, recommendationsNum,
		closingLetterNumID, letterAbstract,
	)
	if strings.Contains(numberingXML, `w:numId="200"`) {
		return numberingXML
	}
	return strings.Replace(numberingXML, "</w:numbering>", extra+"</w:numbering>", 1)
}
