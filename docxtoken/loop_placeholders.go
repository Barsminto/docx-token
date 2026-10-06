package docxtoken

import (
	"github.com/liboyang/docx-token/internal/domain"
	"github.com/liboyang/docx-token/internal/template"
)

// BuildLoopPlaceholders resolves each {BEGIN:path:prefix} region against data and produces
// keys like record_0.title, record_0.listIndex (listIndex is always 1..n).
func BuildLoopPlaceholders(data map[string]interface{}, specs []template.LoopSpec, items []domain.LineItem) (map[string]int, map[string]string) {
	counts := make(map[string]int, len(specs))
	out := make(map[string]string)
	primaryPath := template.PrimaryRecordDataPath(specs)

	for _, spec := range specs {
		rows, err := SliceAtPath(data, spec.DataPath)
		if err != nil {
			continue
		}
		counts[spec.Begin] = len(rows)
		var enriched []domain.LineItem
		if spec.DataPath == primaryPath && len(items) == len(rows) {
			enriched = items
		}
		domain.BindLoopPrefix(out, spec.ItemPrefix, rows, enriched)
	}
	return counts, out
}
