# Content Processor

This document describes the architecture for converting database content into its final HTML form.

## Core Principle

The database is the single source of truth for content, metadata, and configuration flags (e.g., draft status, comments enabled). File generation (`.md` or `.html`) is treated as output only, never as input to be parsed again. This avoids file parsing overhead and keeps all logic tied to structured, queryable data.

## In-Memory Processing

Content is transformed from database records to HTML in an in-memory pipeline for performance and simplicity:

1.  **Fetch:** Load `Content` objects with metadata (`Title`, `Draft`, `PublishedAt`, etc.) and raw Markdown `Body`.
2.  **Business Logic:** Apply rules directly on the struct (e.g., skip drafts).
3.  **Markdown Conversion:** Pass `Body` to the Markdown converter, producing an HTML string.
4.  **Final Object:** Build an in-memory struct with metadata and rendered HTML, ready for templating.

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

The rendering engine uses a hierarchical layout system to generate final HTML pages. This provides both simplicity for the majority of use cases and flexibility for custom designs.

### Layout Hierarchy and Logic

1.  **Content Inheritance:** All content is associated with a `Section`. Even content at the root of the site belongs to a default `root` section. Content inherits the layout assigned to its parent section.
2.  **Section-Specific Layout:** A `Section` can have a specific layout assigned to it. This allows, for example, a "Blog" section to have a different look and feel from a "Portfolio" section.
3.  **Default Layout:** If a section does not have a specific layout assigned, the system uses a global `default` layout. This ensures that all content can be rendered out-of-the-box.

### The Default Layout

The `default` layout is designed to be minimalist, emphasizing readability with clean typography and generous whitespace.
*   **Typography:** The primary font for body content is **Montserrat**. A second, aesthetically compatible font may be used for contrast in elements like headings. A maximum of two fonts will be used.

### Custom Layouts

Users can upload custom layouts via the web interface. These layouts must be configured to accept the rendered content body and associated metadata. The specifics of how content is injected into the layout will be detailed separately.

### Fallback and Seeding Mechanism

The layout system is designed for resilience, ensuring the site can always be rendered.

*   **Database-First:** The system always attempts to retrieve layouts from the database first.
*   **Embedded Fallback:** A copy of the `default` layout is embedded into the application binary (`assets.FS`). In the exceptional case that the database is unavailable or the default layout is missing, the system will use this embedded version to render content.
*   **Notification:** When the fallback mechanism is triggered, a notification will be logged to alert the administrator, suggesting a manual or automatic re-seeding of the layouts from the embedded source.
*   **Initial Seeding:** During the initial application setup, the database will be seeded with the `default` layout, ensuring it is available from the start.

## Companion Blocks

Companion blocks are dynamic sections of content (e.g., lists of related articles) displayed alongside the main content to provide relevant follow-up reading suggestions. The goal is to enhance user engagement without creating distractions.

### Design Philosophy

-   **Placement:** All companion blocks appear at the bottom of the page, after the main content. This maintains a focused, book-like reading experience, avoiding sidebars or other intrusive elements.
-   **Context-Awareness:** The type and content of the blocks are determined by the `Content-Type` of the main document.

### Behavior by Content-Type

The blocks displayed are specific to the type of content being viewed.

-   **`Page`:**
    -   By default, `Page` content does not display any companion blocks. They are considered structurally complete and self-contained. This behavior is iterative and may be adjusted in the future.

-   **`Article`:**
    -   `Article` content can display up to four types of companion blocks, providing a tiered system of recommendations from most to least specific.

    1.  **Same Section, Similar Tags (Default: ON):**
        -   **Content:** Links to other articles within the *same section* that share similar `tags`.
        -   **Purpose:** Offers the most contextually relevant follow-up reading.
        -   **Control:** Can be disabled on a per-article basis via a flag in the database.

    2.  **Same Section, Any Tag (Default: ON):**
        -   **Content:** Links to other articles within the *same section*, regardless of tags.
        -   **Purpose:** Provides broader, section-level discovery.

    3.  **Different Section, Similar Tags (Default: OFF):**
        -   **Content:** Links to articles in *other sections* that share similar `tags`.
        -   **Purpose:** Encourages cross-pollination of content across the site.
        -   **Control:** Disabled by default. Must be enabled via a hard-coded application configuration flag, as it carries a risk of mismatched context.

    4.  **Different Section, Any Tag (Default: OFF):**
        -   **Content:** Links to articles in *any other section*, regardless of tags.
        -   **Purpose:** A generic "keep reading" feature for maximum content exposure.
        -   **Control:** Disabled by default. Must be enabled via a hard-coded application configuration flag.

-   **`Blog`:**
    -   Blogs are chronological collections of posts within a larger `Section` (e.g., `/{section-name}/blog`). Navigation is kept simple and focused.
    -   **Block 1: Recent Posts:** Displays a list of the 'N' most recent posts from the same blog.
    -   **Block 2: Related by Tag:** Shows a list of posts within the same blog that share the same `tags` as the current post.

-   **`Series`:**
    -   This content type prioritizes a focused, sequential reading experience, minimizing distractions.
    -   **Block 1: Sequential Navigation:**
        -   Provides a `<- Previous` link to the prior article in the series, if one exists.
        -   Provides a `Next ->` link to the subsequent article in the series, if one exists.
    -   **Block 2: Direct Navigation:**
        -   Displays a complete list of all articles in the series for direct access.
        -   Articles that come before the current one are prefixed with `<-`.
        -   Articles that come after the current one are suffixed with `->`.

## Header and Hero Image Strategy

To ensure visual consistency and address the common design challenge of placing text over images, the site employs a global (site-wide) strategy for rendering header images and titles. The desired style is set via a single application configuration key (e.g., `ssg.header.style`). This approach prevents visual fragmentation where different sections might use conflicting styles.

The following styles are supported:

### 1. `separate` (Default)

This is the default, minimalist style, prioritizing maximum clarity and readability.

-   **Layout:** The header image is displayed as a distinct block. The page title (`<h1>`) is rendered immediately below it, with no overlap.
-   **Pros:** Guarantees perfect text legibility regardless of the image content.
-   **Cons:** Presents a more traditional, stacked layout which may feel less modern than integrated "hero" designs.

### 2. `overlay`

This style achieves a modern "hero" look by placing the title directly over the image while ensuring legibility.

-   **Layout:** The header image serves as a full-width background for the hero section. A semi-transparent dark overlay is applied on top of the entire image. The page title is then rendered over this overlay.
-   **Pros:** Creates a cohesive, immersive hero unit. Text contrast is maintained against a darkened background.
-   **Cons:** The overall brightness and color vibrancy of the source image are slightly reduced.

### 3. `boxed-text`

This is a sophisticated hybrid approach that preserves the integrity of the main image while ensuring title legibility.

-   **Layout:** The header image is displayed in its entirety in the hero area. The page title is placed inside a "box" positioned at the bottom of the image, floating slightly above the edge.
-   **Functionality:**
    -   The background of this text box is not a solid color, but a `blur` effect applied to the portion of the image directly behind it. This creates a frosted-glass effect that isolates the text.
    -   The main focal point of the image remains sharp and unaffected.
-   **Customization:** This style is controlled by several metadata properties to adapt to different images:
    -   `header_text_align`: Aligns the text within the box (`left`, `center`, `right`).
    -   `header_text_color`: Sets the title's font color (`light` or `dark`).
    -   `header_blur_style`: Determines if the blur effect behind the text is light-toned or dark-toned (`light` or `dark`) to best suit the text color.

#### Default Placeholder Image

To ensure design consistency, the system utilizes a default placeholder image for any content that does not have a specific `header_image` defined. This placeholder is a subtle, high-quality image embedded within the application binary, ensuring that every page can be rendered with a complete header, even without user-provided images.

## Phased Implementation

*   **Phase 1 – Raw HTML:** Convert Markdown from the database to HTML and save to `html/` directory. Ignore layouts.
*   **Phase 2 – Templating:** Feed render-ready objects into layouts to produce full HTML pages.
*   **Phase 3 – Media:** Handle images:
    1.  Users write Markdown image syntax (`![alt](/images/foo.png)`).
    2.  Source files go in `~/documents/assets/images/`.
    3.  Processor converts to `<img>` tags.
    4.  Build step copies assets to the public output directory.

Future work: image optimization, responsive sizes (`srcset`), WebP conversion.