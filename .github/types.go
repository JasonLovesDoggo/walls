package main

import "fmt"

type File struct {
	Url         string
	Description string
}

type Collection struct {
	Name        string
	Files       []File
	BrowseUrl   string
	Description string
}

func (c *Collection) String() string {
	return fmt.Sprintf("%s: %s", c.Name, c.Description)
}
