# Content Processor

This document describes the architecture for converting database content into its final HTML form.

## Core Principle

The database is the single source of truth for content, metadata, and configuration flags (e.g., draft status, comments enabled). File generation (`.md` or `.html`) is treated as output only, never as input to be parsed again. This avoids file parsing overhead and keeps all logic tied to structured, queryable data.

## In-Memory Processing

Content is transformed from database records to HTML in an in-memory pipeline for performance and simplicity:

1. **Fetch:** Load `Content` objects with metadata (`Title`, `Draft`, `PublishedAt`, etc.) and raw Markdown `Body`.
2. **Business Logic:** Apply rules directly on the struct (e.g., skip drafts).
3. **Markdown Conversion:** Pass `Body` to the Markdown converter, producing an HTML string.
4. **Final Object:** Build an in-memory struct with metadata and rendered HTML, ready for templating.

This eliminates writing Markdown to disk or parsing frontmatter, since metadata is already in memory.

## Components

### Service

The existing `ssg.Service` orchestrates the pipeline:
- Fetch content from the repository.
- Apply high-level rules.
- Invoke `MarkdownProcessor`.
- Pass render-ready objects to the templating engine.

### Processor

`MarkdownProcessor` handles Markdown-to-HTML conversion:

- **Responsibility:** Input Markdown, output HTML.
- **Implementation:** `internal/feat/ssg/processor.go`.
- **Library:** `goldmark` with extensions (tables, syntax highlighting).

## Layouts and Templating

Rendering uses layouts, selected by hierarchy:

1. **Section Layout:** If defined in the content’s section.
2. **Default Layout:** Otherwise, use a global default.

The default layout is minimalist, emphasizing readability with clean typography and whitespace. Users can also upload custom layouts and assign them to sections if the defaults do not fit their needs.

### Fallback

Layouts are database-driven but also embedded in the binary as defaults. If the database layouts are unavailable, the system falls back to the embedded ones. A command will allow re-seeding database templates from embedded sources.

## Phased Implementation

* **Phase 1 – Raw HTML:** Convert Markdown from the database to HTML and save to `html/` directory. Ignore layouts.
* **Phase 2 – Templating:** Feed render-ready objects into layouts to produce full HTML pages.
* **Phase 3 – Media:** Handle images:
  1. Users write Markdown image syntax (`![alt](/images/foo.png)`).
  2. Source files go in `~/documents/assets/images/`.
  3. Processor converts to `<img>` tags.
  4. Build step copies assets to the public output directory.

Future work: image optimization, responsive sizes (`srcset`), WebP conversion.

