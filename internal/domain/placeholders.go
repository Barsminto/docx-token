package domain

import (
	"fmt"
	"reflect"
	"strings"
)

// LoopItemPrefix is the object prefix inside {BEGIN_RECORD}…{END_RECORD} before expansion.
const LoopItemPrefix = "Record"

// Key builds a placeholder key such as Report.Title.
func Key(parts ...string) string {
	return strings.Join(parts, ".")
}

// Wrap formats a key as a docx token literal, e.g. {Report.Title}.
func Wrap(key string) string {
	return "{" + key + "}"
}

// LoopField builds a prototype loop field key, e.g. Record.Category.
func LoopField(field string) string {
	return Key(LoopItemPrefix, field)
}

// Placeholders returns a flat map for go-docx: object fields become Object.Field keys.
func (d Document) Placeholders() map[string]string {
	out := make(map[string]string)
	bindObject(out, "Report", d.Report)
	bindObject(out, "Summary", d.Summary)
	for i, item := range d.Items {
		bindObject(out, fmt.Sprintf("%s_%d", LoopItemPrefix, i), item)
	}
	return out
}

func bindObject(dst map[string]string, prefix string, value any) {
	rv := reflect.ValueOf(value)
	rt := rv.Type()
	if rt.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if field.PkgPath != "" {
			continue
		}
		dst[Key(prefix, field.Name)] = fmt.Sprint(rv.Field(i).Interface())
	}
}
