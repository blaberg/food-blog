package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"gopkg.in/yaml.v2"
)

//go:embed templates/*.html
var templates embed.FS

func main() {
	templs, err := template.New("").ParseFS(templates, "templates/*.html")
	if err != nil {
		panic(err)
	}
	markdown := goldmark.New()
	if err := os.MkdirAll("public", os.ModePerm); err != nil {
		panic(err)
	}
	start, err := os.Create("public/index.html")
	if err != nil {
		panic(err)
	}
	defer start.Close()
	if err := filepath.WalkDir("recipes", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		dirName := fmt.Sprintf("public/%s", strings.TrimSuffix(d.Name(), ".md"))
		if err := os.MkdirAll(dirName, os.ModePerm); err != nil {
			panic(err)
		}
		start.Write([]byte(fmt.Sprintf("<a href=\"%s\">%s</a>", strings.TrimSuffix(d.Name(), ".md"), strings.TrimSuffix(d.Name(), ".md"))))
		bs, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		m, l := parseFrontMatter(string(bs))

		f, err := os.Create(fmt.Sprintf("%s/index.html", dirName))
		if err != nil {
			return err
		}
		defer f.Close()
		var buf bytes.Buffer
		context := parser.NewContext()
		if err := markdown.Convert([]byte(l), &buf, parser.WithContext(context)); err != nil {
			panic(err)
		}
		return templs.ExecuteTemplate(f, "base.html", struct {
			Title   string
			Content template.HTML
		}{
			Title:   m.Title,
			Content: template.HTML(strings.TrimSpace(buf.String())),
		})
	}); err != nil {
		panic(err)
	}
}

type Metadata struct {
	Title string `yaml:"title"`
}

func parseFrontMatter(content string) (Metadata, string) {
	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		log.Fatalf("Invalid front matter format")
	}
	var meta Metadata
	err := yaml.Unmarshal([]byte(parts[1]), &meta)
	if err != nil {
		log.Fatalf("Failed to parse front matter: %v", err)
	}

	return meta, parts[2] // Metadata and the remaining content
}
