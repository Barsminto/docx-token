package domain

import "fmt"

// BindLoopPrefix writes {prefix_0.field} keys from row maps and optional enriched LineItems.
// listIndex is always 1..n (not taken from the database row).
func BindLoopPrefix(dst map[string]string, itemPrefix string, rows []map[string]interface{}, items []LineItem) {
	for i, row := range rows {
		key := fmt.Sprintf("%s_%d", itemPrefix, i)
		dst[Key(key, "listIndex")] = fmt.Sprintf("%d", i+1)
		if row != nil {
			for k, v := range row {
				dst[Key(key, k)] = fmt.Sprint(v)
			}
		}
		if i < len(items) {
			bindIssue(dst, key, items[i])
		}
	}
}
