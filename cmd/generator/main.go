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

type Link struct {
	Title string
	URL   string
}

func main() {
	markdown := goldmark.New()
	if err := os.MkdirAll("public", os.ModePerm); err != nil {
		log.Fatalf("Failed to create public directory: %v", err)
	}

	links := make([]Link, 0)
	if err := filepath.WalkDir("recipes", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}
		dirName := fmt.Sprintf("public/%s", strings.TrimSuffix(d.Name(), ".md"))
		if err := os.MkdirAll(dirName, os.ModePerm); err != nil {
			return err
		}
		bs, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		m, l, err := parseFrontMatter(path, string(bs))
		if err != nil {
			return err
		}
		links = append(links, Link{
			Title: m.Title,
			URL:   strings.TrimSuffix(d.Name(), ".md"),
		})

		f, err := os.Create(fmt.Sprintf("%s/index.html", dirName))
		if err != nil {
			return err
		}
		defer f.Close()
		var buf bytes.Buffer
		if err := markdown.Convert([]byte(l), &buf); err != nil {
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
		log.Fatalf("Failed to build recipe pages: %v", err)
	}
	f, err := os.Create("public/index.html")
	if err != nil {
		log.Fatalf("Failed to create public/index.html: %v", err)
	}
	defer f.Close()
	if err := startPage.Execute(f, struct {
		Title   string
		Links   []Link
		CSSFile string
	}{
		Title:   "Johns Recept",
		Links:   links,
		CSSFile: "./output.css",
	}); err != nil {
		log.Fatalf("Failed to render the start page: %v", err)
	}
	cmd := exec.Command("npx", "@tailwindcss/cli", "-i", "input.css", "-o", "../../public/output.css", "-m")
	cmd.Dir = "cmd/generator"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to build the CSS with Tailwind: %v", err)
	}
}

type Metadata struct {
	Title string `yaml:"title"`
}

// parseFrontMatter splits a recipe into its YAML front matter and the markdown
// body that follows it. path is only used to name the file in error messages.
func parseFrontMatter(path, content string) (Metadata, string, error) {
	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		return Metadata{}, "", fmt.Errorf("%s: invalid front matter format", path)
	}
	var meta Metadata
	if err := yaml.Unmarshal([]byte(parts[1]), &meta); err != nil {
		return Metadata{}, "", fmt.Errorf("%s: failed to parse front matter: %w", path, err)
	}

	return meta, parts[2], nil // Metadata and the remaining content
}
