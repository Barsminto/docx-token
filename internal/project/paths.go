package project

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	DirTemplates = "templates"
	DirOutput    = "output"

	FileTemplate = "template.docx"
	FileOutput   = "report.docx"
)

// Root 从当前工作目录向上查找包含 go.mod 的目录作为项目根。
func Root() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("未找到项目根（go.mod）")
		}
		dir = parent
	}
}

func TemplatePath() (string, error) {
	root, err := Root()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, DirTemplates, FileTemplate), nil
}

func OutputPath() (string, error) {
	root, err := Root()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, DirOutput, FileOutput), nil
}
