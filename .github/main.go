package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

const (
	RootTemplateFile     = "./README.md.tmpl"
	CategoryTemplateFile = "./README.md.category.tmpl"
	OutputFile           = "README.md"
)

func main() {
	// Collect all collections
	collections := ParseDirectories("../")

	// Create root README
	err := createRootReadme(collections)
	if err != nil {
		panic(err)
	}

	fmt.Println("Done")
}

func createRootReadme(collections []Collection) error {
	// Note: if you wish to have the README in the root directory, ensure to edit RootTemplateFile to remove the
	/// ../ references
	readmePath := filepath.Join(BaseDir, ".github", OutputFile)

	// Parse the root README template
	tmpl, err := template.New("RootReadme").Parse(RootReadmeTemplate)
	if err != nil {
		return err
	}

	// Create the README file
	f, err := os.Create(readmePath)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	// Execute the template
	err = tmpl.Execute(f, collections)
	if err != nil {
		return err
	}

	return nil
}
