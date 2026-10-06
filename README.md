# docx-token

Go 库：用 **`map[string]interface{}`** 填充 **Group Investigation（GI）** 风格的 Word **`.docx`** 模板，保留 Word 多级列表（1.1 / a. b. c. / 3.1…）与页眉占位符。

**其他项目只需依赖公开包** `github.com/liboyang/docx-token/docxtoken`，不要 import `internal/*`。

---

## 快速集成

### 1. 添加依赖

```bash
go get github.com/liboyang/docx-token/docxtoken
```

（若仓库尚未发布 tag，可在 `go.mod` 中使用 `replace` 指向本地或 fork 路径。）

### 2. 准备模板文件

将本仓库中的模板复制到你的项目（路径自定），例如：

| 文件 | 用途 |
|------|------|
| `templates/template.docx` | 正式填充用（无文首说明） |
| `templates/reference-template.docx` | 含全部标记说明，做新模板时对照 |

在本仓库重新生成模板：

```bash
go run ./cmd/gentemplate                          # → templates/template.docx
go run ./cmd/gentemplate -reference               # → templates/reference-template.docx
```

### 3. 调用 API 生成报告

**推荐：`Generate` + `Options`（模板路径、输出路径由调用方指定）**

```go
package myapp

import (
    "path/filepath"

    "github.com/liboyang/docx-token/docxtoken"
)

func WriteGIReport(templateDir, outputFile string, payload map[string]interface{}) error {
    return docxtoken.Generate(docxtoken.Options{
        TemplatePath: filepath.Join(templateDir, "template.docx"),
        OutputPath:   outputFile,
        Data:         payload,
    })
}
```

**简写：`Fill`**

```go
err := docxtoken.Fill(
    "/data/templates/gi-template.docx",
    "/data/output/case-123.docx",
    payload,
)
```

`TemplatePath` / `OutputPath` 可以是：

- 配置项里的目录 + 文件名
- `embed` 解压后的临时路径（见下文）
- 对象存储下载到本地的路径

---

## 数据格式（`map[string]interface{}`）

支持 **嵌套 map**；标量会展平为 `executiveSummary.identified` 等形式参与替换。  

**列表数据**：路径由模板里的 `{BEGIN:点路径:变量名}` 决定（例如 `case.case_issues` + 区内 `{record.xxx}`）。  
内置 GI 模板固定为 `case.issues` + 变量名 `issue`（`{issue.brief}`）。每条循环自动带 **`{变量.listIndex}`**（1、2、3…，不由数据库提供）。

### 最小示例

```go
data := map[string]interface{}{
    "projectName":          "PROJECT ALPHA",
    "draftDate":            "01 February 2026",
    "accountableExecutive": "Name, Title",
    "subjectList":          "Subject A (PSID 001), Subject B (PSID 002)",
    "background.period":    "February 2026",
    "case": map[string]interface{}{
        "name": "Investigation case title",
        "issues": []map[string]interface{}{
            {
                "title":   "CONCERN 1",
                "brief":   "Executive summary list line",
                "summary": "Text after Concern: in INVESTIGATION block",
            },
            {
                "title":   "CONCERN 2",
                "brief":   "...",
                "summary": "...",
            },
        },
    },
    "executiveSummary": map[string]interface{}{
        "identified":      "...",
        "concernRef":      "...",
        "substantiation":  "not substantiated",
        "evidence":        "did not establish",
        "recommendations": "...",
    },
    "conclusion": map[string]interface{}{
        "support":        "...",
        "policyBreaches": "...",
    },
    "recommendation": map[string]interface{}{
        "1": "...",
        "2": "...",
        "3": "...",
    },
}
```

### Concern 列表的兼容键名

解析顺序（任选其一即可）：

1. `data["case"]["issues"]` — **推荐**
2. `data["records"]` / `data["items"]`（见 `docxtoken.RecordsKey`）

每条 issue / record 支持的字段（大小写兼容见 `mapdata.go`）：

| 字段 | 模板占位符（循环内） |
|------|----------------------|
| `title` | `{issue.title}` |
| `brief` | `{issue.brief}` |
| `summary` | `{issue.summary}` |
| （自动） | `{issue.listIndex}`、`{issue.numberedTitle}`、`{issue.sectionNumber}` 等由程序赋值 |

### 程序自动写入的占位符

| 键 | 说明 |
|----|------|
| `{Conclusion.Number}` | 结论章号 = 首个调查章号 + concern 条数 |
| `{Recommendations.Number}` | 建议章号 = 结论 + 1 |
| `{Report.Title}` / `{Report.Date}` / `{Report.GeneratedAt}` | 来自 `projectName` 与生成时间（可覆盖） |

页眉可使用：`{projectName}`、`{case.name}`、`{draftDate}` 等（与正文相同规则）。

---

## 模板约定（自定义模板必读）

### 循环区域（可自定义数据路径与变量名）

语法：

```text
{BEGIN:<map 中的点路径>:<循环变量名>}
  … 段落 …
{END:<同一段路径>:<同一变量名>}
```

示例（数据在 `case.case_issues`，模板里用 `record.*`）：

```text
{BEGIN:case.case_issues:record}
{record.code}: {record.title}
{END:case.case_issues:record}
```

```go
"case": map[string]interface{}{
    "case_issues": []map[string]interface{}{
        {"code": "C-1", "title": "Alpha", "brief": "…"},
    },
},
```

- **数据路径** `<点路径>`：从 `map` 里取 `[]map[string]interface{}`（与路径一致）。
- **循环变量** `<变量名>`：区内写 `{record.xxx}`；展开后变为 `{record_0.xxx}`、`{record_1.xxx}` …
- **`listIndex`**：每条自动注入 `{record.listIndex}` → `1`、`2`、`3`…（**不读数据库**）。若库里有同名字段，以自动序号为准。
- 仍支持内置 GI：`{BEGIN:case.issues}` / `:background` / `:record`（数据路径均为 `case.issues`，变量名 **`issue`**）。

| 区域 | 开始 | 结束 |
|------|------|------|
| Executive concern 列表 | `{BEGIN:case.issues}` | `{END:case.issues}` |
| Background concern 列表 | `{BEGIN:case.issues:background}` | `{END:case.issues:background}` |
| 每条 INVESTIGATION 整块 | `{BEGIN:case.issues:record}` | `{END:case.issues:record}` |

循环内一行示例（内置变量名 `issue`）：

`Concern {issue.listIndex}: {issue.brief}`

（不要手打 `a.`；字母由 Word 列表生成。）

程序还会为每条 concern 计算 `numberedTitle`、`sectionNumber` 等（与 `issue` / 主 record 循环同路径时，写到 `{record_0.numberedTitle}` 或 `{issue_0.numberedTitle}`）。

### 小节标记

| 标记 | 含义 |
|------|------|
| `{.}` | 十进制小节（1.1、3.1、8.1）— 行首；内置 `gentemplate` 生成后会去掉文字并挂 Word 列表 |
| `{~}` | 字母项（a.、b.），多用于 CONCLUSION 下 |

在 Word 里手改模板时：**章标题加粗即可**，不必写 `{GI_BOLD}`（仅内置生成器使用）。

### 章号逻辑

- 模板中 `{BEGIN:case.issues:record}` **之前**的 `1.`、`2.` … 决定第一个调查章从几开始（通常 **3**）。
- 每条 concern 占一章（3、4、5…）；结论/建议章号随 concern **条数**变化，不依赖模板里写死的 8/9。

完整标记列表见 **`templates/reference-template.docx`**（文首说明区仅用于对照，正式出报告请用 `template.docx`）。

---

## 项目结构

```
docx-token/
├── docxtoken/          # 公开 API（其他项目只 import 这里）
│   ├── api.go          # Generate(Options)
│   ├── api.go              # Generate(Options)
│   ├── fill.go             # Fill — 发现循环、展开、写文件
│   ├── loop_placeholders.go # 按 BEGIN 路径绑定 prefix_0.* + listIndex
│   ├── path.go             # SliceAtPath(data, "case.case_issues")
│   └── mapdata.go          # map → Document（章号用主 record 循环路径）
├── internal/
│   ├── domain/         # Document、LineItem、占位符键、章号
│   ├── render/         # 循环展开、OOXML 列表补丁、go-docx 替换
│   ├── template/       # 内置 GI 模板生成、numbering.xml、页眉
│   ├── sample/         # 演示数据 FillMap
│   └── project/        # 本仓库 CLI 默认路径（库用户可忽略）
├── templates/
│   ├── template.docx
│   └── reference-template.docx
├── cmd/
│   ├── docx-token/     # 命令行填充示例
│   └── gentemplate/    # 生成/参考模板
└── output/             # CLI 默认输出目录（可 .gitignore）
```

### 填充流水线（实现概要）

1. **`DiscoverLoopSpecsFromDocx`**：解析模板中所有 `{BEGIN:path:prefix}`。
2. **`BuildFromMap` + `BuildLoopPlaceholders`**：按 path 取 slice，生成 `record_0.*` / `issue_0.*`（含 **listIndex**）。
3. **`FirstIssueSeq` + `AssignIssueSeq`**：调查章、列表序号（主 record 循环的 path）。
4. **`expandRecordBlocks`**：每个 BEGIN 区域按对应 slice 长度复制 XML，`{prefix.` → `{prefix_0.`。
4. **`patchDocumentXML`**：调查/结论/建议的 Word `numId`、concern 字母列表、页眉占位符修复等。
5. **`go-docx` `ReplaceAll`**：替换正文 + **页眉/页脚**中的 `{token}`。
6. 写出 `OutputPath`。

---

## 在其他项目中的常见用法

### 配置目录

```go
type ReportConfig struct {
    TemplateDir string // e.g. "/opt/myapp/templates" or "./assets/docx"
    OutputDir   string
}

func (c ReportConfig) Generate(caseID string, data map[string]interface{}) error {
    out := filepath.Join(c.OutputDir, caseID+".docx")
    return docxtoken.Generate(docxtoken.Options{
        TemplatePath: filepath.Join(c.TemplateDir, "template.docx"),
        OutputPath:   out,
        Data:         data,
    })
}
```

### 与 `embed` 一起用

```go
// 启动时把 embed 的模板写到临时目录，或解压到固定缓存路径，再将路径传给 Generate。
// 库本身不读 embed；只接受文件路径。
```

### HTTP 服务示例（伪代码）

```go
func (h *Handler) DownloadReport(w http.ResponseWriter, r *http.Request) {
    data := h.buildPayloadFromDB(r.Context(), r.PathValue("id"))
    out := filepath.Join(os.TempDir(), "report-"+id+".docx")
    if err := docxtoken.Generate(docxtoken.Options{
        TemplatePath: h.cfg.GITemplatePath,
        OutputPath:   out,
        Data:         data,
    }); err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    http.ServeFile(w, r, out)
}
```

---

## 本仓库 CLI

```bash
# 生成内置模板
go run ./cmd/gentemplate
go run ./cmd/gentemplate -reference

# 用默认 templates/template.docx + 示例数据填充
go run ./cmd/docx-token -skip-template -out output/report.docx

# 指定模板与输出路径
go run ./cmd/docx-token -skip-template \
  -template /path/to/template.docx \
  -out /path/to/report.docx
```

---

## 开发与测试

```bash
go test ./...
```

---

## 许可证

与仓库根目录 LICENSE 一致（若未添加 LICENSE，发布到 GitHub 前请补充）。
