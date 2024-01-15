package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"text/template"

	"gopkg.in/yaml.v3"
)

//go:generate go run main.go

const (
	DatabaseFile = "esv.json"
	CategoryFile = "categories.yaml"
	TemplateFile = "verses.go.tmpl"
	OutputFile   = "../../../pkg/ref/canonical.go"
)

type BooksConfig struct {
	Books []BookConfig `json:"books"`
}

type BookConfig struct {
	Name   string  `json:"name"`
	Verses [][]int `json:"verses"`
}

type CategoriesConfig struct {
	Categories map[string][]string `yaml:"categories"`
}

func loadDatabase() (*BooksConfig, error) {
	esvj, err := os.ReadFile(DatabaseFile)
	if err != nil {
		return nil, err
	}

	var bookConfig BooksConfig
	err = json.Unmarshal(esvj, &bookConfig)
	if err != nil {
		return nil, err
	}

	return &bookConfig, nil
}

func loadCategories() (*CategoriesConfig, error) {
	catj, err := os.ReadFile(CategoryFile)
	if err != nil {
		return nil, err
	}

	var catConfig CategoriesConfig
	err = yaml.Unmarshal(catj, &catConfig)
	if err != nil {
		return nil, err
	}

	return &catConfig, nil
}

func main() {
	bookConfig, err := loadDatabase()
	if err != nil {
		panic(err)
	}

	catConfig, err := loadCategories()
	if err != nil {
		panic(err)
	}

	tmplBytes, err := os.ReadFile(TemplateFile)
	if err != nil {
		panic(err)
	}

	tmpl := template.New("verses")
	tmpl.Funcs(map[string]interface{}{
		"Mod": func(a, b int) bool { return a%b == 0 },
	})
	_, err = tmpl.Parse(string(tmplBytes))
	if err != nil {
		panic(err)
	}

	fh, err := os.Create(OutputFile)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(fh, struct {
		Books      []BookConfig
		Categories map[string][]string
	}{
		Books:      bookConfig.Books,
		Categories: catConfig.Categories,
	})
	if err != nil {
		panic(err)
	}

	err = exec.Command("go", "fmt", OutputFile).Run()
	if err != nil {
		panic(err)
	}
}
