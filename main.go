package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/rwestlund/gotex"
	"github.com/spf13/pflag"
)

//go:embed resume.tex.tmpl
var defaultTemplate string

// A package-scoped variable is used here so that it can be overridden at compile time
var defaultOutputFilename string = "output.pdf"

func main() {
	inputFilename := pflag.StringP("input", "i", "resume.json", "Resume input filename")
	outputFilename := pflag.StringP("output", "o", defaultOutputFilename, "Output PDF filename")
	texCommand := pflag.StringP("command", "c", "xelatex", "TeX command")
	templateFilename := pflag.StringP("template", "t", "", "Path to template file")
	shouldRedact := pflag.BoolP("redact", "r", false, "Redact sensitive information in resume")

	shouldShowDefault := pflag.BoolP("show-default", "", false, "Print default template contents to stdout")

	shouldHelp := pflag.BoolP("help", "h", false, "Display help message")
	pflag.CommandLine.MarkHidden("help")

	pflag.Parse()

	if *shouldHelp {
		pflag.PrintDefaults()
		os.Exit(0)
	}

	if *shouldShowDefault {
		fmt.Print(defaultTemplate)
		os.Exit(0)
	}

	data, err := os.ReadFile(*inputFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", *inputFilename, err)
		os.Exit(1)
	}

	var r Resume

	// Default to unmarshalling as JSON
	var unmarshalFunc func([]byte, interface{}) error = json.Unmarshal

	// Unmarshal as TOML if the file extension is "toml"
	if filepath.Ext(strings.ToLower(*inputFilename)) == ".toml" {
		unmarshalFunc = toml.Unmarshal
	}

	err = unmarshalFunc(data, &r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing file: %v\n", err)
		os.Exit(1)
	}

	if *shouldRedact {
		r.redact()
	}

	funcMap := template.FuncMap{
		"sub": func(a, b int) int {
			return a - b
		},
	}

	var templateContents string
	switch *templateFilename {
	case "":
		templateContents = defaultTemplate
	default:
		fileContents, err := os.ReadFile(*templateFilename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		templateContents = string(fileContents)
	}

	t, err := template.New("resume.template").Funcs(funcMap).Parse(templateContents)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	b := bytes.Buffer{}
	err = t.Execute(&b, r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing template: %v\n", err)
		os.Exit(1)
	}

	// Write executed template to stdout if --output == "-"
	// Otherwise, run LaTeX and write the PDF to a file
	switch *outputFilename {
	case "-":
		io.Copy(os.Stdout, &b)
	default:
		filename := *outputFilename
		file, err := os.Create(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating file %s: %v\n", filename, err)
			os.Exit(1)
		}
		defer file.Close()

		pdfBytes, err := gotex.Render(b.String(), gotex.Options{Command: *texCommand})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error rendering PDF: %v\n", err)
			os.Exit(1)
		}

		err = os.WriteFile(filename, pdfBytes, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file %s: %v\n", filename, err)
			os.Exit(1)
		}
	}
}
