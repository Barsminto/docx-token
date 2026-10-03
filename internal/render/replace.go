package render

import docx "github.com/lukasjarosch/go-docx"

// ReplaceFromMap 使用 map 中的键值对替换文档内全部匹配的占位符。
func ReplaceFromMap(document *docx.Document, values map[string]string) error {
	placeholders := docx.PlaceholderMap{}
	for key, value := range values {
		placeholders[key] = value
	}
	return document.ReplaceAll(placeholders)
}
