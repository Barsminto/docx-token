package render

import (
	"regexp"
	"strings"
)

var placeholderKeyRE = regexp.MustCompile(`\{[A-Za-z0-9_.]+\}`)

// repairFragmentedPlaceholders merges Word runs that split tokens like {Background.Period} across <w:t> nodes.
func repairFragmentedPlaceholders(documentXML string) string {
	return paragraphRE.ReplaceAllStringFunc(documentXML, func(p string) string {
		plain := paragraphPlainText(p)
		if !strings.Contains(plain, "{") {
			return p
		}
		for _, key := range placeholderKeyRE.FindAllString(plain, -1) {
			if strings.Contains(p, key) {
				continue
			}
			return rebuildParagraphPlain(p, plain)
		}
		return p
	})
}
