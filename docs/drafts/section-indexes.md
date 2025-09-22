# Section Indexes

> *Note: This document specifies how section index pages are generated. For a detailed breakdown of content kinds (`Page`, `Article`, `Blog`, `Series`) and their primary purpose, please refer to the `blocks.md` draft.*

---

## Index Page

An index page acts as a container for a section's chronological content.

### Structure and Appearance

-   **Layout:** A single-column, paginated list.
-   **Order:** Chronological, from newest to oldest.
-   **Row Content:** Each entry in the list links to the full content and serves as a preview, containing:
    -   A thumbnail image (from the content's main image).
    -   A summary or excerpt.
    -   Relevant metadata (e.g., publication date, tags).

---

## Aggregation Rules

### Included Content

The following content kinds are always included in their section's index:

-   `Article`
-   `Blog`
-   `Series`

### Excluded Content

-   `Page` content is never listed in a chronological index.

### Section Landings

A section's index can be manually controlled. If a `Page` with the slug `index` exists within a section, the system will render that page directly **instead of** generating the automatic chronological index. This allows for custom, handcrafted landing pages for any section.

It is then the content editor's responsibility to include any index-like listings if desired, as the automatic generation will be bypassed entirely.

---

## Index Scope

### Local Indexes

An index for a specific section (e.g., `/news/`) contains only content belonging directly to that section. This includes content from its dedicated blog path (e.g., posts from `/news/blog/` will appear in the `/news/` index).

### Global Index

The root section's index (`/`) is a special case. It acts as a global aggregator, containing all `Article`, `Blog`, and `Series` content from the entire site. This consolidation includes content from the root section itself, all other sections, and the blogs within those sections (e.g., posts from `/news/blog/` will also appear in the global index at `/`).

---

## Blog Indexes

In addition to the main section indexes, each blog instance has its own dedicated, chronological index.

Unlike section indexes which can aggregate multiple content kinds, a blog index **only** lists `Blog` posts belonging to that specific blog path.

### Root Blog Index

-   **Path:** `/blog/`
-   **Content:** Lists all posts from the root section's blog, paginated and ordered chronologically.

### Section Blog Index

-   **Path:** `/{section-name}/blog/`
-   **Content:** Lists all posts from that specific section's blog, paginated and ordered chronologically.

---

## Series Indexes

Similar to blogs, each series has a dedicated index page that lists only the posts belonging to that series.

The primary ordering for these indexes is determined by the predefined sequence of the posts, which allows for flexible arrangement (e.g., inserting new posts between existing ones). Chronological order should be used as a fallback where a sequence is not explicitly defined.

### Root Series Index

-   **Path:** `/{series-name}/index.html`

### Section Series Index

-   **Path:** `/{section-name}/{series-name}/index.html`

---

## Index URLs

Index pages are generated with the following URL structure:

-   **Root Index:** `/index.html`
-   **Section Index:** `/<section-name>/index.html`
-   **Paginated Index:** `/<section-name>/page/2/index.html`

---
