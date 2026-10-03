package docxtoken

import (
	"fmt"

	"github.com/liboyang/docx-token/internal/render"
)

// Fill reads templatePath, expands loops, replaces tokens from data, writes outputPath.
//
// data is map[string]interface{}:
//   - scalars: "Report.Title", "Summary.Overview", or nested {"Report": {"Title": "..."}}
//   - records: []map[string]interface{}{{"Title":"...", "Category":"..."}}
func Fill(templatePath, outputPath string, data map[string]interface{}) error {
	doc, err := BuildFromMap(data)
	if err != nil {
		return fmt.Errorf("build from map: %w", err)
	}
	extras := ScalarPlaceholders(data)
	if err := render.Generate(templatePath, outputPath, doc, extras); err != nil {
		return fmt.Errorf("generate: %w", err)
	}
	return nil
}
