package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"text/template"
)

//go:generate go run main.go

const (
	DatabaseFile = "esv.json"
	TemplateFile = "verses.go.tmpl"
	OutputFile   = "../../../internal/esv/verses.go"
)

type BooksConfig struct {
	Books []BookConfig `json:"books"`
}

type BookConfig struct {
	Name   string  `json:"name"`
	Verses [][]int `json:"verses"`
}

func main() {
	esvj, err := ioutil.ReadFile(DatabaseFile)
	if err != nil {
		panic(err)
	}

	var bookConfig BooksConfig
	err = json.Unmarshal(esvj, &bookConfig)
	if err != nil {
		panic(err)
	}

	tmplBytes, err := ioutil.ReadFile(TemplateFile)
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

	err = tmpl.Execute(fh, bookConfig)
	if err != nil {
		panic(err)
	}

	err = exec.Command("go", "fmt", OutputFile).Run()
	if err != nil {
		panic(err)
	}
}
