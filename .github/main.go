package main

import (
	"fmt"
	"os"
	"text/template"
)

const (
	TemplateFile = "./README.md.tmpl"
	OutputFile   = "../README.md"
)

var collections = []Collection{
	{
		Name: "Anime",
		Files: []File{
			{Url: "https://img.freepik.com/free-photo/anime-moon-landscape_23-2151645908.jpg", Description: "Anime moon landscape"},
			{Url: "https://images.squarespace-cdn.com/content/v1/5fe4caeadae61a2f19719512/0c94a8d5-9587-4b1e-a818-252e22deaa88/Screenshot+%281728%29.jpg", Description: "Anime landscape"},
			{Url: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSRLi839Gq2WnEPxDe_nl7xf_qwQaTcOuiRnw&s", Description: "Anime landscape"},
		},
		BrowseUrl:   "https://github.com/natsukagami/anime-wallpapers",
		Description: "Anime wallpapers",
	},
	{
		Name: "Catppuccin",
		Files: []File{
			{Url: "https://github.com/zhichaoh/catppuccin-wallpapers/blob/main/solids/bkg1.png?raw=true", Description: "Catppuccin solid 1"},
			{Url: "https://github.com/zhichaoh/catppuccin-wallpapers/blob/main/solids/bkg5.png?raw=true", Description: "Catppuccin solid 5"},
		},
		BrowseUrl:   "https://github.com/zhichaoh/catppuccin-wallpapers",
		Description: "Catppuccin wallpapers",
	},
}

func main() {
	tmpl, err := os.ReadFile(TemplateFile)
	if err != nil {
		panic(err)
	}

	parsedTemplate := template.Must(template.New("README").Parse(string(tmpl)))

	// Run after parsing so if there are any errors, they will be caught before truncating the existing README.

	o, err := os.Create(OutputFile)
	if err != nil {
		panic(err)
	}

	err = parsedTemplate.Execute(o, collections)
	if err != nil {
		panic(err)
	}
	fmt.Println("Done")

}
