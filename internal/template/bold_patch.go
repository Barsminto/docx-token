package template

import (
	"regexp"
	"strings"
)

// GiBoldMarker marks chapter titles: bold body text, not Heading 1 (stripped on save).
const GiBoldMarker = "{GI_BOLD}"

var headingStyleRE = regexp.MustCompile(`<w:pStyle w:val="Heading1"\s*/>`)

func stripGiBoldMarkers(documentXML string) string {
	return paragraphBlockRE.ReplaceAllStringFunc(documentXML, func(p string) string {
		if !strings.Contains(p, GiBoldMarker) {
			return p
		}
		p = strings.ReplaceAll(p, GiBoldMarker, "")
		p = headingStyleRE.ReplaceAllString(p, "")
		return ensureRunBold(p)
	})
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
