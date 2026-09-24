# Mina favoritrecept

A small recipe blog, published at <https://recipes.johnblaberg.com>.

The site is a handful of markdown files rendered to static HTML by a Go
generator in this repo, styled with Tailwind and hosted on GitHub Pages.

## Layout

| Path | What it is |
| --- | --- |
| `recipes/*.md` | The content. YAML front matter with a `title`, then markdown. **The filename is the slug**, so `recipes/tomatsoppa.md` becomes `/tomatsoppa/`. |
| `cmd/generator/` | The generator. `templates/` holds the Go HTML templates, `input.css` is the Tailwind entry point. |
| `backend/` | A small development server for previewing the built site. |
| `public/` | Generated output. Gitignored, and rebuilt from scratch — never edit it by hand. |

## Requirements

Go (see the version in `go.mod`), plus Node and npm for the Tailwind CLI.
CI builds with Node 22.

## Building

```sh
npm ci
go run ./cmd/generator
```

**Run both from the repository root.** The generator resolves its templates,
`input.css` and its output path relative to the working directory, so it fails
immediately from anywhere else.

## Previewing locally

```sh
go run ./backend
```

Then open <http://localhost:8080>. This serves `public/` at the root, which is
how the site is served in production — the pages use root-relative links, so
previewing them from a subpath would break navigation.

## Adding a recipe

Create `recipes/<slug>.md`:

```markdown
---
title: Namnet på receptet
---

### Ingredienser
- ...

### Gör så här
1. ...
```

Rebuild, and it is picked up automatically and linked from the start page.

## Deploying

Pushing to `master` builds the site and publishes it to GitHub Pages. Note the
`paths:` filter in `.github/workflows/deploy.yaml` — only changes under
`cmd/`, `recipes/` or to the build manifests trigger a deploy. Use the
workflow's manual *Run workflow* button to republish without a commit.

The custom domain is configured in the repository's Settings → Pages, not by a
`CNAME` file in this repo. There is deliberately no `CNAME` file; adding one is
not needed.
