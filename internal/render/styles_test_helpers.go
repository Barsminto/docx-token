package render

import (
	"fmt"
	"strings"
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
