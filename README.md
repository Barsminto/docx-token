# docx-token

Go library: fill **docx** templates from **`map[string]interface{}`**.

## Library usage

```go
import "github.com/liboyang/docx-token/docxtoken"

err := docxtoken.Fill("templates/template.docx", "output/report.docx", map[string]interface{}{
    "Report.Title":       "FY26 Q1 Report",
    "Summary.Author":     "Jane Doe",
    "Report.GeneratedAt": "2026-10-03 15:30", // optional; default time.Now()
    "Summary.Overview":   "Executive summary text.",
    docxtoken.RecordsKey: []map[string]interface{}{
        {"Title": "INC-1001 Hardware", "Category": "Hardware"},
        {"Title": "INC-1002 License", "Category": "Software"},
    },
})
```

Nested maps are supported:

```go
map[string]interface{}{
    "Report": map[string]interface{}{"Title": "FY26 Q1"},
    docxtoken.RecordsKey: []map[string]interface{}{...},
}
```

### Map keys

| Key | Meaning |
|-----|---------|
| `Report.*`, `Summary.*` | Scalar placeholders (`Report.Title`, `Report.GeneratedAt`, …) |
| `records` (or `items`) | Issue rows for `{BEGIN_ISSUE_LIST}` / `{BEGIN_RECORD}` loops |
| Row fields | `Title` (required), `Category` (detail token) |

Issue section numbers are read from the template (`1. Summary` → issues start at `2`). Heading text uses `{Record.NumberedTitle}` or `{Record.Seq}. {Record.Title}`.

## CLI (demo)

```bash
go run ./cmd/gentemplate
go run ./cmd/docx-token -skip-template
```

## Dependencies

- [lukasjarosch/go-docx](https://github.com/lukasjarosch/go-docx) — `{token}` replacement
- [gomutex/godocx](https://github.com/gomutex/godocx) — optional built-in template generator
