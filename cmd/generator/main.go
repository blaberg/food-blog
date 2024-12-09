package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"gopkg.in/yaml.v2"
)

var startPage = template.Must(template.ParseFiles(
	"cmd/generator/templates/base.html",
	"cmd/generator/templates/start.html",
))

var recipePage = template.Must(template.ParseFiles(
	"cmd/generator/templates/base.html",
	"cmd/generator/templates/recipe.html",
))

func main() {
	markdown := goldmark.New()
	if err := os.MkdirAll("public", os.ModePerm); err != nil {
		panic(err)
	}

	links := make([]string, 0)
	if err := filepath.WalkDir("recipes", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		dirName := fmt.Sprintf("public/%s", strings.TrimSuffix(d.Name(), ".md"))
		if err := os.MkdirAll(dirName, os.ModePerm); err != nil {
			panic(err)
		}
		links = append(links, strings.TrimSuffix(d.Name(), ".md"))
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
			return err
		}
		return recipePage.Execute(f, struct {
			Title   string
			Content template.HTML
			CSSFile string
		}{
			Title:   m.Title,
			Content: template.HTML(strings.TrimSpace(buf.String())),
			CSSFile: "../output.css",
		})
	}); err != nil {
		panic(err)
	}
	f, err := os.Create("public/index.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := startPage.Execute(f, struct {
		Title   string
		Recipes []string
		CSSFile string
	}{
		Title:   "John's Recipes",
		Recipes: links,
		CSSFile: "./output.css",
	}); err != nil {
		panic(err)
	}
	if err := exec.Command("tailwindcss", "-i", "cmd/generator/input.css", "-o", "public/output.css").Run(); err != nil {
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
