package render

import (
	"fmt"
	"strings"

	"github.com/liboyang/docx-token/internal/domain"
)

// patchIssueHeadings replaces heading tokens in document.xml before go-docx runs.
// go-docx skips runs with two placeholders like "{Record.Seq}. {Record.Title}".
func patchIssueHeadings(documentXML string, items []domain.LineItem) string {
	for i, item := range items {
		prefix := fmt.Sprintf("Record_%d", i)
		combined := fmt.Sprintf("{%s.Seq}. {%s.Title}", prefix, prefix)
		documentXML = strings.ReplaceAll(documentXML, combined, item.NumberedTitle)
		documentXML = strings.ReplaceAll(documentXML, fmt.Sprintf("{%s.NumberedTitle}", prefix), item.NumberedTitle)
	}
	return documentXML
}
