package main

import (
	"fmt"
	"os"
	"text/template"
)

const (
	TemplateFile = "./README.md.tmpl"
	OutputFile   = "README.md"
)

func main() {
	tmpl, err := os.ReadFile(TemplateFile)
	if err != nil {
		panic(err)
	}

	parsedTemplate := template.Must(template.New("README").Parse(string(tmpl)))

	// Run after parsing so if there are any errors, they will be caught before truncating the existing README.

	o, err := os.Create("../" + OutputFile)
	if err != nil {
		panic(err)
	}

	collections := ParseDirectories("../")

	err = parsedTemplate.Execute(o, collections)
	if err != nil {
		panic(err)
	}
	fmt.Println("Done")

}
