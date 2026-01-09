package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"gopkg.in/yaml.v3"
)

//go:generate go run main.go

const (
	DatabaseFile              = "esv.json"
	CategoryFile              = "categories.yaml"
	AbbreviationsFile         = "abbr.yaml"
	VerseTemplateFile         = "verses.go.tmpl"
	AbbreviationsTemplateFile = "abbrs.go.tmpl"
	VerseOutputFile           = "../../../pkg/ref/canonical.go"
	AbbreviationsOutputFile   = "../../../pkg/ref/abbr.go"
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

type OrdinalConfig struct {
	Standard string   `yaml:"standard"`
	Accept   []string `yaml:"accept"`
}

type BookAbbrConfig struct {
	Name     string   `yaml:"name"`
	Standard string   `yaml:"standard"`
	Singular string   `yaml:"singular"`
	Ordinal  string   `yaml:"ordinal"`
	Accept   []string `yaml:"accept"`
}

type AbbreviationsConfig struct {
	Ordinals []OrdinalConfig   `yaml:"ordinals"`
	Books    []*BookAbbrConfig `yaml:"books"`
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

func loadAbbreviations() (*AbbreviationsConfig, error) {
	abbrj, err := os.ReadFile(AbbreviationsFile)
	if err != nil {
		return nil, err
	}

	var abbrConfig AbbreviationsConfig
	err = yaml.Unmarshal(abbrj, &abbrConfig)
	if err != nil {
		return nil, err
	}

	for _, abbr := range abbrConfig.Books {
		if abbr.Ordinal != "" {
			found := false
			for _, ord := range abbrConfig.Ordinals {
				if ord.Standard == abbr.Ordinal {
					found = true

					accept := make([]string, 0, len(ord.Accept)*len(abbr.Accept))
					for _, acc := range abbr.Accept {
						for _, ordAcc := range ord.Accept {
							accept = append(accept, ordAcc+acc)
						}
					}

					abbr.Accept = accept

					break
				}
			}

			if !found {
				return nil, fmt.Errorf("book named %q has bad ordinal configuration", abbr.Name)
			}
		}
	}

	return &abbrConfig, nil
}

func templateVerses() error {
	bookConfig, err := loadDatabase()
	if err != nil {
		return err
	}

	catConfig, err := loadCategories()
	if err != nil {
		return err
	}

	return applyTemplate(
		"verses",
		VerseTemplateFile,
		VerseOutputFile,
		struct {
			Books      []BookConfig
			Categories map[string][]string
		}{
			Books:      bookConfig.Books,
			Categories: catConfig.Categories,
		},
	)
}

func templateAbbreviations() error {
	abbrConfig, err := loadAbbreviations()
	if err != nil {
		return err
	}

	return applyTemplate(
		"abbrs",
		AbbreviationsTemplateFile,
		AbbreviationsOutputFile,
		struct {
			Abbreviations []*BookAbbrConfig
		}{
			Abbreviations: abbrConfig.Books,
		},
	)
}

func applyTemplate(
	name string,
	tmplFile string,
	outFile string,
	vars any,
) error {
	tmplBytes, err := os.ReadFile(tmplFile)
	if err != nil {
		return err
	}

	tmpl := template.New(name)
	tmpl.Funcs(map[string]any{
		"Mod": func(a, b int) bool { return a%b == 0 },
	})
	_, err = tmpl.Parse(string(tmplBytes))
	if err != nil {
		return err
	}

	fh, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer fh.Close()

	// Write generated file comment (obfuscated to avoid triggering tools that scan for these phrases)
	obfuscated := "Ly8gQ29kZSBnZW5lcmF0ZWQgYnkgLi90b29scy9nZW4vdmVyc2VzL21haW4uZ287IERPIE5PVCBFRElULgoK"
	comment, err := base64.StdEncoding.DecodeString(obfuscated)
	if err != nil {
		return err
	}
	_, err = fh.Write(comment)
	if err != nil {
		return err
	}

	err = tmpl.Execute(fh, vars)
	if err != nil {
		return err
	}

	err = fh.Close()
	if err != nil {
		return err
	}

	err = exec.Command("go", "fmt", outFile).Run()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := templateVerses()
	if err != nil {
		panic(err)
	}

	err = templateAbbreviations()
	if err != nil {
		panic(err)
	}
}
