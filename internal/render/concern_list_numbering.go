package render

import (
	"regexp"
	"strings"
)

var (
	concernListTokenRE = regexp.MustCompile(`\{[a-zA-Z_][a-zA-Z0-9_]*(?:_\d+)?\.listIndex\}`)
	concernListFilledRE = regexp.MustCompile(`Concern \d+:`)
)

// reapplyConcernListNumIDs: executive list numId 1 (a…), background list numId 2 (restarts at a).
func reapplyConcernListNumIDs(documentXML string) string {
	paragraphs := paragraphRE.FindAllString(documentXML, -1)
	if len(paragraphs) == 0 {
		return documentXML
	}

	seenBackground := false
	for i, p := range paragraphs {
		plain := paragraphPlainText(p)
		if strings.Contains(plain, "BACKGROUND") {
			seenBackground = true
			continue
		}
		if !isConcernListLine(plain) {
			continue
		}
		numID := 1
		if seenBackground {
			numID = 2
		}
		paragraphs[i] = injectNumPrIntoParagraph(p, numID, 0)
	}

	return replaceParagraphsInOrder(documentXML, paragraphs)
}

func isConcernListLine(plain string) bool {
	if !strings.Contains(plain, "Concern ") {
		return false
	}
	return concernListTokenRE.MatchString(plain) || concernListFilledRE.MatchString(plain)
}
