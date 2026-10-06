package render

import (
	"fmt"
	"os"

	docx "github.com/lukasjarosch/go-docx"

	"github.com/liboyang/docx-token/internal/domain"
	"github.com/liboyang/docx-token/internal/template"
)

// Generate reads issue sequence from the template, expands the loop, and replaces tokens.
// extras are merged into placeholders (for custom scalar keys from map[string]interface{}).
func Generate(templatePath, outputPath string, doc domain.Document, extras map[string]string) error {
	raw, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("read template: %w", err)
	}

	firstSeq, err := template.FirstIssueSeq(raw)
	if err != nil {
		return fmt.Errorf("read issue sequence from template: %w", err)
	}
	doc.AssignIssueSeq(firstSeq)
	concNum := firstSeq + len(doc.Items)
	recNum := concNum + 1

	expanded, err := expandRecordBlocks(raw, len(doc.Items))
	if err != nil {
		return err
	}
	expanded, err = patchDocumentXML(expanded, doc.Items, concNum, recNum)
	if err != nil {
		return err
	}

	document, err := docx.OpenBytes(expanded)
	if err != nil {
		return fmt.Errorf("open docx: %w", err)
	}
	defer document.Close()

	placeholders := doc.Placeholders()
	for k, v := range extras {
		placeholders[k] = v
	}
	placeholders["Conclusion.Number"] = fmt.Sprintf("%d", concNum)
	placeholders["Recommendations.Number"] = fmt.Sprintf("%d", recNum)
	if err := ReplaceFromMap(document, placeholders); err != nil {
		return fmt.Errorf("replace placeholders: %w", err)
	}
	if err := document.WriteToFile(outputPath); err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	return nil
}
