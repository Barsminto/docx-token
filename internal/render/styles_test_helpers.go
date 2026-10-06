package render

import (
	"fmt"
	"os"
	"strings"

	"github.com/liboyang/docx-token/internal/domain"
	"github.com/liboyang/docx-token/internal/template"
)

func sampleScalarStrings(data map[string]interface{}) map[string]string {
	out := flattenSampleMap("", data)
	for k := range out {
		if strings.HasPrefix(k, "issue.") || strings.HasPrefix(k, "issue_") {
			delete(out, k)
		}
	}
	return out
}

func flattenSampleMap(prefix string, data map[string]interface{}) map[string]string {
	out := make(map[string]string)
	for k, v := range data {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		switch val := v.(type) {
		case map[string]interface{}:
			for fk, fv := range flattenSampleMap(key, val) {
				out[fk] = fv
			}
		case []interface{}, []map[string]interface{}:
			continue
		default:
			out[key] = fmt.Sprint(val)
		}
	}
	return out
}

func loopContextForTest(tplPath string, doc domain.Document) ([]template.LoopSpec, map[string]int, map[string]string) {
	raw, err := os.ReadFile(tplPath)
	if err != nil {
		return nil, nil, nil
	}
	specs, err := template.DiscoverLoopSpecsFromDocx(raw)
	if err != nil {
		return nil, nil, nil
	}
	primary := template.PrimaryRecordDataPath(specs)
	counts := make(map[string]int)
	ph := make(map[string]string)
	n := len(doc.Items)
	rows := make([]map[string]interface{}, n)
	for _, spec := range specs {
		counts[spec.Begin] = n
		var enriched []domain.LineItem
		if spec.DataPath == primary {
			enriched = doc.Items
		}
		domain.BindLoopPrefix(ph, spec.ItemPrefix, rows, enriched)
	}
	return specs, counts, ph
}
