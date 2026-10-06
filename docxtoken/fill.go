package docxtoken

import (
	"fmt"
	"os"

	"github.com/liboyang/docx-token/internal/render"
	"github.com/liboyang/docx-token/internal/template"
)

// Fill reads templatePath, discovers {BEGIN:path:prefix} regions, expands them using the
// slice at path in data, injects prefix.listIndex per row, then replaces all {tokens}.
//
// data is map[string]interface{}:
//   - scalars / nested maps → flat keys like executiveSummary.identified
//   - slice at each loop path, e.g. case.case_issues for {BEGIN:case.case_issues:record}
func Fill(templatePath, outputPath string, data map[string]interface{}) error {
	raw, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("read template: %w", err)
	}
	specs, err := template.DiscoverLoopSpecsFromDocx(raw)
	if err != nil {
		return fmt.Errorf("discover loops: %w", err)
	}
	issuePath := template.PrimaryRecordDataPath(specs)
	doc, err := BuildFromMap(data, issuePath)
	if err != nil {
		return fmt.Errorf("build from map: %w", err)
	}
	loopCounts, loopPH := BuildLoopPlaceholders(data, specs, doc.Items)
	extras := ScalarPlaceholders(data)
	if err := render.Generate(templatePath, outputPath, doc, extras, specs, loopCounts, loopPH); err != nil {
		return fmt.Errorf("generate: %w", err)
	}
	return nil
}
