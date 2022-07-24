**mkresume** - A tool for generating resumes from a JSON file 

## Overview
`mkresume` is used to generated a resume PDF, given a LaTeX file template and a JSON file containing the information to be used by the template.

The steps taken are:
* The contents of the JSON file (default:`./resume.json`) are read.
* The data from the JSON file is inserted into the template.
  A default template is embedded into the `mkresume` binary at compile time, but a template file can also be chosen at runtime with the `--template` flag.
* The executed template file is written to a temporary directory.
* A LaTeX command (default: `xelatex`) is used to generate a PDF from the executed template.
* The contents of the generated PDF is written to the specified output file (default: `./output.pdf)

## Compilation
By default, this repository can be compiled with the Go compiler simply by running `go build`.
A few things should be noted, however:

### Specifying a new default template
 During compilation, the file `./resume.tex.tmpl` is read and hardcoded into the `mkresume` binary.
 This is done so that no default template need exist on the filesystem.

 While the `resume.tex.tmpl` that is included with this repo is a working example, it is not a very aesthetically-pleasing resume;
 therefore, you may wish to recompile `mkresume` with your own resume template.

### Specifying a new default output file
By default, generated PDFs are named `output.pdf`, unless another name is specified at runtime with the `--output` flag.
This default can be overridden at compile time by instructing the Go linker to set the value of `main.defaultOutputFilename`.

For example, someone named John Doe may wish to compile by running `go build -ldflags "-X 'main.defaultOutputFilename=john_doe_resume.pdf'"`.
This will set the default output filename to `john_doe_resume.pdf` rather than `output.pdf`.
