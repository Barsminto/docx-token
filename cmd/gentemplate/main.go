package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/liboyang/docx-token/internal/project"
	"github.com/liboyang/docx-token/internal/template"
)

func main() {
	defaultOut, err := project.TemplatePath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "解析项目路径失败: %v\n", err)
		os.Exit(1)
	}
	out := flag.String("out", defaultOut, "输出模板路径（默认项目 templates/template.docx）")
	flag.Parse()

	if err := os.MkdirAll(filepath.Dir(*out), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "创建目录失败: %v\n", err)
		os.Exit(1)
	}
	if err := template.WriteReportTemplate(*out); err != nil {
		fmt.Fprintf(os.Stderr, "生成模板失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("已覆盖生成模板: %s（会覆盖已有文件）\n", *out)
}
