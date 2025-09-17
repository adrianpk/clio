# Content Generation from Database

This document outlines the process of converting the content stored in the database into the actual documentation files.

*   **Single Source of Truth:** The database is and will be the single source of truth for all content.
*   **Reconstructability:** It must be possible to reconstruct the entire content, including its organization and categorization (tags), solely from the database.
*   **Database Engine:** The database is a single SQLite file.
*   **Exclusions:** Images will not be stored in the database.
*   **Versioning via Git:** Exporting content to files allows for versioning the base content in Git/GitHub. This includes both the Markdown files and their associated assets (like images) to ensure content integrity. This provides two key benefits:
    *   **Versioning:** The database will only store the latest version of any given document. We are not interested in creating a massive database file that includes every single change. Git will handle the version history.
    *   **Automation:** Clio will be responsible for committing all content changes (text and assets) to GitHub, based on user configuration (manual, automatic, or hybrid).
*   **Database Reconstruction:** Eventually, it will also be possible to reconstruct the database from a clone of the repository where the base docs are stored. This feature is not a priority for the initial implementation.
*   **Scope:** This process concerns the base Markdown content, not the final generated site that Clio will also version and publish.
*   **Generation Strategy:** The process will be implemented in two phases. Initially, a full regeneration of all Markdown documents will be performed in each batch. Later, this process will be optimized to only regenerate new or modified documents, improving efficiency. A `force` flag will be included to allow for a complete, on-demand regeneration at any time.

## Metadata, Publishing Logic, and UI

### Schema Approach

Considering that the primary database engine is SQLite, which lacks the native JSONB support found in databases like PostgreSQL, the application will favor a well-defined database schema with explicit columns over a flexible JSON-based field for metadata. This approach ensures type safety and data integrity across all supported database engines. The comprehensive `Meta` struct from the reference implementation will serve as a solid reference for all the fields to include in the `content` table, ensuring type safety and data integrity.

This means fields like `title`, `description`, `summary`, `type`, `slug`, `layout`, `header_image`, etc., will eventually become dedicated columns in the database and fields in the Go struct.

The immediate priority for implementation will be to add the following critical fields:
*   `draft` (boolean)
*   `published_at` (timestamp)

### Properties Editor UI

To manage this metadata, a "properties editor" will be available in the content creation and editing views. 

*   **Visibility:** This editor will be **hidden by default** to maintain a clean, simple, and non-distracting writing interface.
*   **Future Enhancements:** A more dynamic metadata editor may be considered in the future.

### Publishing Logic

*   **Draft Status:** A new `draft` boolean field will be the master control for publication. If `draft` is `true`, the content will not be published, regardless of any other settings.
*   **Publication Date:** A `publication_date` field will allow scheduling posts for the future or back-dating them. However, the `draft` status always takes precedence.

## Directory Structure Generation from Database

The process of "dumping" the database content into Markdown files will follow a specific structure to ensure content is organized logically. The file paths will be derived from the `Section` and `Content` entities.

### Guiding Principle: Simplicity by Default

While many options for content organization (Sections, Content Types, etc.) will be available, they must be implemented as ergonomically as possible. The primary goal is to provide a straightforward workflow for users who simply want to publish text online quickly.

*   **Default Layout:** A single, minimalist default layout will be provided out-of-the-box. Users will not be required to create or configure layouts to get started.
*   **Optional Sections:** Sections will not be mandatory. If no sections are configured, all content will be published to the root directory, mimicking a traditional, flat publishing system.
*   **Default Content Type:** The `Article` content type will be the default.

This will ensure a "just write" experience: a user will only need to provide a title and their content to publish it instantly.

### Section-based Paths

The `Section` entity will dictate the base directory for its associated content.

*   **Root Section:** There will be a main, root section which has no `path` defined. Content assigned to this section will reside in the root of the markdown output directory.
*   **Other Sections:** All other sections will have a `path` field (e.g., `/news`, `/tech`). This path will be used to create a corresponding subdirectory. Content assigned to these sections will be placed within their respective directories.
*   **Nesting:** The section structure will be limited to a single level. There will be no nested sections (e.g., `/tech/hardware` is not a valid structure).

### Content-based Paths

The `Content` entity's `type` will determine the final path, sometimes adding an extra level of nesting within its section's directory. There will be four content types initially: `Page`, `Article`, `Blog`, and `Series`.

*   **`Page` and `Article` Types:** Content of these types will be placed directly within their assigned section's path.
    *   Root section: `/{content-slug}.md`
    *   Other section: `/{section-path}/{content-slug}.md`

*   **`Blog` and `Series` Types:** These types will add an additional subdirectory for organization.
    *   **`Blog`:** Will add a `/blog` subdirectory.
        *   Root section: `/blog/{content-slug}.md`
        *   Other section: `/{section-path}/blog/{content-slug}.md`
    *   **`Series`:** Will add a `/{series-name}` subdirectory. The `{series-name}` slug will be defined by a new, dedicated "Series Name" field in the content's database record. This will allow multiple, uniquely named content items to be grouped under the same series.
        *   Root section: `/{series-name}/{content-slug}.md`
        *   Other section: `/{section-path}/{series-name}/{content-slug}.md`

## Asset Management: Images

### Upload Mechanism

A mechanism will be provided to upload images directly from the content editor. This will avoid the need for a separate, pre-populated image repository. The system will generate the final path for the image on-the-fly, allowing it to be inserted into the document instantly (similar to the Jira content editor).

### Storage and Directory Structure

*   **Base Location:** All images will be stored within the `.../assets/images/` directory defined by the workspace.
*   **Mirrored Structure:** The directory structure for images will mirror the content's markdown file structure. A dedicated folder will be created for each content item that has images. 
*   **Example:** A post located at `/{section-path}/{post-slug}.md` will have its associated images stored in a folder at `/{section-path}/{post-slug}/`. All images for that post (e.g., `image1.png`, `photo.jpg`) will reside inside this folder.

#### Alternative Structure Considered

An alternative approach was considered where images would be stored "side-by-side" with their corresponding markdown file, either in the same directory or in a dedicated subdirectory (e.g., `img/`). For example, a post at `.../news/my-post.md` would have its images in `.../news/my-post/img/`. This co-location can simplify versioning and management of content and its assets together. 

However, for the initial implementation, the decision is to maintain a centralized `assets/images` directory that mirrors the content structure, keeping a clean separation between markdown and binary assets at the top level.

### Image Naming Conventions

While most images can have any name, certain filenames will be reserved for specific functionality:

*   **Header Image:** An image named `header.{jpg|png|webp|etc}` within a post's image folder will be automatically used as that post's main header image.
*   **General Images:** Any other image can have a user-defined name.
*   **Future Consideration:** A similar convention will be defined for social media thumbnail images (e.g., `thumbnail-og.png`).
