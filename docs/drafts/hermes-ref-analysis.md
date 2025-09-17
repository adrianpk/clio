# Analysis of Hermes v1 Generation Logic

This document contains notes and analysis from the reference implementation provided in the `./hermesone` directory. The focus is on understanding the core data transformation logic, while consciously avoiding the architectural patterns (e.g., verbose command/service split) that should not be replicated in Clio.

## Core Concepts from Reference

The `hermesone` implementation was file-based (reading from `.md` files) rather than database-first. However, the logic used to handle paths and metadata from the file's front matter can be adapted for our `db -> file` generation process.

### 1. Metadata Structure (`hermes/content.go`)

The `Meta` struct provides a comprehensive list of fields that were used to describe a piece of content. This is a good reference for the data we need to pull from Clio's database for each `Content` entity.

**Key Fields for Path Generation:**
*   `Type`: (Article, Blog, Series, etc.) - Directly influences the path structure.
*   `Section`: The name of the section folder.
*   `Slug`: The filename of the content (without extension).

Other notable fields include `date`, `published-at`, `header-image`, `tags`, etc., which will be needed for generating the front matter in the final Markdown files.

### 2. Path and Asset Logic (`hermes/gen.go`)

The most relevant piece of logic is in the `modifyImagePath` function. Although its purpose was to rewrite image paths during the `md -> html` conversion, it reveals the underlying path rules that were enforced.

**Inferred Logic:**

1.  **Base Asset Path:** There was a global `ImgDir` variable, analogous to our `.../assets/images/` directory.
2.  **Content-Relative Paths:** The function calculates the final image path based on the path of the markdown file containing it (`mdPath`).
3.  **Path Parsing:** It uses `strings.Split` on the markdown file's path to break it into components (e.g., `[section, content-type, slug]`).
4.  **Conditional Path Building:** It uses `switch` and `if` statements on the content type (`parts[1] == ContentType.Blog`) to decide on the final directory structure. For example, if the type was `Blog` or `Series`, it would include that component in the final image path, mirroring the logic we designed.

**Takeaway for Clio:**
Our `db -> markdown` generator will need to implement a similar path-building logic. For each `Content` record from the database, we will:

1.  Fetch the `Content.Type`, `Content.Slug`, and its associated `Section.Path`.
2.  Use these fields to construct the final destination path for the `.md` file, following the rules in our `generation.md` design document.
3.  Similarly, construct the corresponding asset folder path (e.g., `.../assets/images/{section-path}/{slug}/`).

### 3. Content Parsing (`hermes/gen.go`)

The `Parse` function shows a simple but effective way to handle content with front matter. It splits the file content using a `---` separator. 

**Takeaway for Clio:**
When we generate the `.md` files, we will need to do the reverse:
1.  Gather all the metadata for a `Content` item from the database.
2.  Marshal this metadata into a YAML string.
3.  Prepend it to the markdown body, wrapped in `---` separators.

---

**Conclusion:** The reference code, particularly `hermes/gen.go` and `hermes/content.go`, provides a solid conceptual foundation for the logic required to translate database records into a structured directory of Markdown files with associated metadata. The key is to adapt the path-building and YAML generation logic to our database-driven models within Clio's existing architecture.