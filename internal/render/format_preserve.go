package render

import (
	"regexp"
	"strings"
)

var (
	pPrBlockRE = regexp.MustCompile(`(?s)<w:pPr>.*?</w:pPr>`)
	pPrRPrRE   = regexp.MustCompile(`(?s)<w:pPr>.*?<w:rPr>.*?</w:rPr>`)
	numPrBlockRE = regexp.MustCompile(`(?s)<w:numPr>.*?</w:numPr>`)
	wtContentRE  = regexp.MustCompile(`(<w:t(?:\s[^>]*)?>)([^<]*)(</w:t>)`)
	runBlockRE   = regexp.MustCompile(`(?s)<w:r\b[^>]*>.*?</w:r>`)
	runOpenRE    = regexp.MustCompile(`(<w:r\b[^>]*>)`)
)

func defaultRPrFromParagraph(paragraph string) string {
	if m := pPrRPrRE.FindString(paragraph); m != "" {
		start := strings.Index(m, "<w:rPr>")
		if start >= 0 {
			return m[start:]
		}
	}
	for _, m := range regexp.MustCompile(`(?s)<w:rPr>.*?</w:rPr>`).FindAllString(paragraph, -1) {
		if strings.Contains(m, "rFonts") {
			return m
		}
	}
	return regexp.MustCompile(`(?s)<w:rPr>.*?</w:rPr>`).FindString(paragraph)
}

// stripMarkersInParagraph removes structural tokens only inside paragraph XML (keeps pPr / runs).
func stripMarkersInParagraph(paragraph string, markers ...string) string {
	out := paragraph
	for _, m := range markers {
		if m != "" {
			out = strings.ReplaceAll(out, m, "")
		}
	}
	return out
}

func trimParagraphTextRuns(paragraph string) string {
	return replaceParagraphTextPreservingRuns(paragraph, strings.TrimSpace(paragraphPlainText(paragraph)))
}

// replaceParagraphTextPreservingRuns updates text but keeps w:pPr (indents, numPr) and per-run rPr.
func replaceParagraphTextPreservingRuns(paragraph, text string) string {
	if !wtContentRE.MatchString(paragraph) {
		return rebuildParagraphPlain(paragraph, text)
	}
	first := true
	return wtContentRE.ReplaceAllStringFunc(paragraph, func(m string) string {
		parts := wtContentRE.FindStringSubmatch(m)
		if len(parts) != 4 {
			return m
		}
		if !first {
			return parts[1] + parts[3]
		}
		first = false
		return parts[1] + xmlEscape(text) + parts[3]
	})
}

func rebuildParagraphPlain(paragraph, plain string) string {
	openEnd := strings.Index(paragraph, ">")
	if openEnd < 0 {
		return paragraph
	}
	openTag := paragraph[:openEnd+1]
	closeTag := "</w:p>"

	pPr := pPrBlockRE.FindString(paragraph)
	rPr := defaultRPrFromParagraph(paragraph)
	run := "<w:r>"
	if rPr != "" {
		run += rPr
	}
	run += "<w:t xml:space=\"preserve\">" + xmlEscape(plain) + "</w:t></w:r>"
	return openTag + pPr + run + closeTag
}

func inferTemplateBodyRPr(documentXML string) string {
	for _, p := range paragraphRE.FindAllString(documentXML, -1) {
		rPr := defaultRPrFromParagraph(p)
		if strings.Contains(rPr, "Arial") {
			return rPr
		}
	}
	return ""
}

// inheritParagraphRunFonts copies paragraph-level or sibling run fonts onto runs that lack rFonts.
func inheritParagraphRunFonts(documentXML string) string {
	fallback := inferTemplateBodyRPr(documentXML)
	return paragraphRE.ReplaceAllStringFunc(documentXML, func(p string) string {
		rPr := defaultRPrFromParagraph(p)
		if rPr == "" {
			rPr = fallback
		}
		if rPr == "" {
			return p
		}
		return runBlockRE.ReplaceAllStringFunc(p, func(run string) string {
			if strings.Contains(run, "rFonts") {
				return run
			}
			if strings.Contains(run, "<w:rPr>") {
				return run
			}
			return runOpenRE.ReplaceAllString(run, "$1"+rPr)
		})
	})
}

// trimEmptyParagraphsBeforeHeading removes spacer empty lines before a heading (template page-fill artifacts).
func trimEmptyParagraphsBeforeHeading(documentXML, heading string) string {
	paragraphs := paragraphRE.FindAllString(documentXML, -1)
	indices := paragraphRE.FindAllStringIndex(documentXML, -1)
	if len(paragraphs) != len(indices) {
		return documentXML
	}
	target := -1
	for i, p := range paragraphs {
		plain := strings.TrimSpace(paragraphPlainText(p))
		if strings.EqualFold(plain, strings.TrimSpace(heading)) {
			target = i
			break
		}
	}
	if target <= 0 {
		return documentXML
	}
	removeFrom := target
	for i := target - 1; i >= 0; i-- {
		if isSpacerParagraph(paragraphs[i]) {
			removeFrom = i
			continue
		}
		break
	}
	if removeFrom == target {
		return documentXML
	}
	var b strings.Builder
	last := 0
	for i, idx := range indices {
		if i >= removeFrom && i < target {
			last = idx[1]
			continue
		}
		b.WriteString(documentXML[last:idx[0]])
		b.WriteString(paragraphs[i])
		last = idx[1]
	}
	b.WriteString(documentXML[last:])
	return b.String()
}

func isSpacerParagraph(paragraph string) bool {
	if strings.TrimSpace(paragraphPlainText(paragraph)) != "" {
		return false
	}
	if strings.Contains(paragraph, "w:drawing") || strings.Contains(paragraph, "w:pict") {
		return false
	}
	return true
}

func ensurePageBreakBeforeHeading(documentXML, heading string) string {
	want := strings.TrimSpace(heading)
	return paragraphRE.ReplaceAllStringFunc(documentXML, func(p string) string {
		plain := strings.TrimSpace(paragraphPlainText(p))
		if !strings.EqualFold(plain, want) {
			return p
		}
		if strings.Contains(p, "w:pageBreakBefore") {
			return p
		}
		br := `<w:pageBreakBefore/>`
		if strings.Contains(p, "<w:pPr>") {
			return strings.Replace(p, "<w:pPr>", "<w:pPr>"+br, 1)
		}
		return strings.Replace(p, "<w:p>", "<w:p><w:pPr>"+br+"</w:pPr>", 1)
	})
}
