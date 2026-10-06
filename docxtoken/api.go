package docxtoken

import "fmt"

// Options is the public entry for report generation.
// TemplatePath and OutputPath are absolute or relative paths on the local filesystem.
// Data is the full report payload (nested maps and case.issues / records slice).
type Options struct {
	TemplatePath string
	OutputPath   string
	Data         map[string]interface{}
}

// Generate expands template loops, applies Word numbering patches, replaces all
// {placeholders}, and writes OutputPath.
//
// It is safe to call from another service or HTTP handler; paths are chosen by the caller
// (e.g. templates next to your binary, a shared volume, or a temp file).
func Generate(opts Options) error {
	if opts.TemplatePath == "" {
		return fmt.Errorf("docxtoken: TemplatePath is required")
	}
	if opts.OutputPath == "" {
		return fmt.Errorf("docxtoken: OutputPath is required")
	}
	return Fill(opts.TemplatePath, opts.OutputPath, opts.Data)
}
