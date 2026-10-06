package docxtoken

import (
	"fmt"
	"strings"
)

// SliceAtPath reads a []map[string]interface{} at a dot path (e.g. case.case_issues).
func SliceAtPath(data map[string]interface{}, path string) ([]map[string]interface{}, error) {
	if data == nil || path == "" {
		return nil, nil
	}
	cur := interface{}(data)
	for _, seg := range strings.Split(path, ".") {
		m, ok := cur.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("path %q: not a map at %q", path, seg)
		}
		cur = m[seg]
	}
	return coerceMapSlice(cur)
}

func coerceMapSlice(raw interface{}) ([]map[string]interface{}, error) {
	if raw == nil {
		return nil, nil
	}
	switch list := raw.(type) {
	case []map[string]interface{}:
		return list, nil
	case []interface{}:
		out := make([]map[string]interface{}, 0, len(list))
		for _, elem := range list {
			m, ok := elem.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("slice entries must be map[string]interface{}")
			}
			out = append(out, m)
		}
		return out, nil
	default:
		return nil, fmt.Errorf("expected slice at path, got %T", raw)
	}
}
