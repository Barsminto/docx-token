package render

import (
	"fmt"
	"strings"

	"github.com/liboyang/docx-token/internal/domain"
)

func patchIssueHeadings(documentXML string, items []domain.LineItem) string {
	for i, item := range items {
		prefix := fmt.Sprintf("issue_%d", i)
		documentXML = strings.ReplaceAll(documentXML, fmt.Sprintf("{%s.numberedTitle}", prefix), item.NumberedTitle)
	}
	return documentXML
}
