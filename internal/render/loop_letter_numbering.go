package render

import (
	"fmt"
	"strings"

	"github.com/liboyang/docx-token/internal/template"
)

const (
	loopLetterNumIDBase   = 500
	loopLetterNumIDStride = 200
	recordLetterNumIDBase = 350
)

func loopLetterNumID(specIndex, copyIndex int) int {
	return loopLetterNumIDBase + specIndex*loopLetterNumIDStride + copyIndex
}

// applyLoopCopyLetterNumID gives each duplicated loop block its own letter list (a., b., c. restart).
// Mark lines with {~} in the template (same as CONCLUSION letter items).
func applyLoopCopyLetterNumID(chunk string, numID int) string {
	if !strings.Contains(chunk, template.TildeMarker) {
		return chunk
	}
	return paragraphRE.ReplaceAllStringFunc(chunk, func(p string) string {
		if !strings.Contains(p, template.TildeMarker) {
			return p
		}
		p = stripMarkersInParagraph(p, template.TildeMarker)
		return injectNumPrIntoParagraph(p, numID, 0)
	})
}

func mergeLoopLetterNumberingDefs(numberingXML string, specs []template.LoopSpec, counts map[string]int) string {
	if !strings.Contains(numberingXML, "<w:numbering") {
		return numberingXML
	}
	letterAbstract := letterListAbstractNumID(numberingXML)
	var extra strings.Builder
	for si, spec := range specs {
		n := counts[spec.Begin]
		for i := 0; i < n; i++ {
			numID := loopLetterNumID(si, i)
			if strings.Contains(numberingXML, `w:numId="`+itoa(numID)+`"`) {
				continue
			}
			extra.WriteString(fmt.Sprintf(
				`<w:num w:numId="%d"><w:abstractNumId w:val="%d"/>`+
					`<w:lvlOverride w:ilvl="0"><w:startOverride w:val="1"/></w:lvlOverride></w:num>`,
				numID, letterAbstract,
			))
		}
	}
	if extra.Len() == 0 {
		return numberingXML
	}
	return strings.Replace(numberingXML, "</w:numbering>", extra.String()+"</w:numbering>", 1)
}

func mergeRecordLetterNumberingDefs(numberingXML string, itemCount int) string {
	if itemCount == 0 || !strings.Contains(numberingXML, "<w:numbering") {
		return numberingXML
	}
	letterAbstract := letterListAbstractNumID(numberingXML)
	var extra strings.Builder
	for i := 0; i < itemCount; i++ {
		numID := recordLetterNumIDBase + i
		if strings.Contains(numberingXML, `w:numId="`+itoa(numID)+`"`) {
			continue
		}
		extra.WriteString(fmt.Sprintf(
			`<w:num w:numId="%d"><w:abstractNumId w:val="%d"/>`+
				`<w:lvlOverride w:ilvl="0"><w:startOverride w:val="1"/></w:lvlOverride></w:num>`,
			numID, letterAbstract,
		))
	}
	if extra.Len() == 0 {
		return numberingXML
	}
	return strings.Replace(numberingXML, "</w:numbering>", extra.String()+"</w:numbering>", 1)
}
