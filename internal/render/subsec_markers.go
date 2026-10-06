package render

import (
	"strings"

	"github.com/liboyang/docx-token/internal/template"
)

// stripSubsectionMarkers removes leftover structural tags after numbering is applied.
func stripSubsectionMarkers(documentXML string) string {
	out := documentXML
	for _, m := range []string{template.DotMarker, template.TildeMarker} {
		out = strings.ReplaceAll(out, m, "")
	}
	return out
}
