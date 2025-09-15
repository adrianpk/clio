# Content to Markdown to HTML Generation - Draft

## Core Approach & Key Considerations:

*   **On-Demand Conversion is Key:** We're not looking for Clio to eagerly convert Markdown to HTML every time content is saved. Instead, this conversion needs to be an explicit, on-demand action. Think of it as something the user triggers when they're ready.

*   **Two-Step Generation Pipeline:** We'll break down the content generation into two distinct phases:
    1.  **Database to Markdown (Filesystem):** Our `clio.db` remains the single source of truth. We'll export content records from the database into Markdown files on the local filesystem. This involves taking the metadata from our DB columns and embedding it as YAML "frontmatter" within these Markdown files. This step gives us a solid filesystem backup and acts as a necessary intermediate stage before HTML generation.
    2.  **Markdown (Filesystem) to HTML (Filesystem):** Once we have our Markdown files on the filesystem, we'll process them to create static HTML files. These HTML files will be structured and ready for direct deployment to static hosting platforms, with GitHub Pages being our primary target.

*   **Optimizations:** To keep things snappy and avoid unnecessary work, we'll build in some smart checks:
    *   **DB to MD Export:** We'll only regenerate a `.md` file from the database if the `updated_at` timestamp of its corresponding content record in `clio.db` has actually changed since the last `.md` file was generated. Of course, if the `.md` file somehow goes missing, we'll regenerate it.
    *   **MD to HTML Generation:** Similarly, an `.html` file will only be regenerated if its source `.md` file has been modified since the last `.html` generation. We'll also include a "force flag" for users who want to explicitly override this and force a full regeneration.

*   **Directory Structure:**
    *   All our generated content will live under a configurable root directory, starting with `~/Documents/Clio/`.
    *   Inside this root, we'll have two main subdirectories:
        *   `md/`: This is where all our exported Markdown files (with frontmatter) will reside.
        *   `html/`: This will contain the final static HTML site, all set for deployment.

*   **Focus on Text First:** For this initial phase, we're concentrating purely on text-based content (Markdown and HTML generation). Image management, things like asset referencing and copying, will be tackled in a later development phase.

*   **Content Paths:** The way we organize our Markdown and HTML output will directly follow Clio's content structure:
    *   **`root` Section:** We'll have a special `root` section for the main landing page of the site, corresponding to the base path (`/`).
    *   **Named Sections:** Other content sections will be defined by a `section-name` (e.g., `cars`, `news`) and stored in the database. These will logically map to distinct URL paths (e.g., `/cars`, `/news`). Our filesystem structure within the `md/` and `html/` directories will mirror this logical routing.

## Action Plan:

### Initial Setup

1.  **`docs/draft.md` (This Document):** This document itself serves as a plan and context guide.
2.  **Review & Adapt `Content` Struct:**
    *   **Goal:** Make sure the `Content` struct in `internal/feat/ssg/content.go` can hold all the necessary metadata that will eventually become YAML frontmatter in our Markdown files. This will likely mean adding some new fields.
    *   **Action:** Add fields for both the original Markdown content (`BodyMarkdown`) and the generated HTML (`BodyHTML`) to the `Content` struct.

### Implementing the DB → Markdown (Filesystem) Export

1.  **Define the "Export to Markdown" Action:**
    *   **Goal:** Create the user-facing way to kick off the DB-to-MD export.
    *   **Action:** Implement a new web `handler` function or integrate this functionality into an existing API endpoint (e.g., `/api/export/md`).
2.  **Develop the Export Logic:**
    *   **Goal:** Get content from the database, build our Markdown files with frontmatter, and write them to the `md/` directory.
    *   **Action:**
        *   We'll loop through `Content` records from `clio.db`.
        *   For each record, we'll dynamically build YAML frontmatter using its metadata fields.
        *   Then, we'll combine this frontmatter with the `BodyMarkdown` content.
        *   We'll figure out the right target file path within `~/Documents/Clio/md/`, respecting the content's section and slug.
        *   We'll apply our optimization logic: checking the `updated_at` timestamp against the existing `.md` file's modification time, and making sure the `.md` file is actually there before writing.
3.  **Error Handling & Logging:**
    *   **Goal:** Ensure a clean export process.
    *   **Action:** We'll integrate comprehensive error handling and detailed logging throughout the export process.

### Implementing the Markdown (Filesystem) → HTML (Filesystem) Generation

1.  **Define the "Generate HTML" Action:**
    *   **Goal:** Create the user-facing way to kick off the MD-to-HTML generation.
    *   **Action:** Implement a new web `handler` function or integrate this functionality into an existing or new web `handler` function.
2.  **Develop the Generation Logic:**
    *   **Goal:** Read our Markdown files, convert them to HTML, and write the output to the `html/` directory.
    *   **Action:**
        *   Loop through the Markdown files in `~/Documents/Clio/md/`.
        *   Apply our optimization logic: checking the modification date of the source `.md` file against the existing `.html` file's modification time, and making sure the `.html` file is there before regenerating.
        *   Read each `.md` file, parse its YAML frontmatter, and pull out the raw Markdown content.
        *   We'll use a suitable Go library for Markdown to HTML conversion (we've looked at `blackfriday/v2` and `goldmark`).
        *   Implement a custom renderer (taking inspiration from `hermes`'s original `TailwindRenderer`) to apply specific HTML styling (like CSS classes) and adjust internal links/image paths to fit Clio's output structure.
        *   We'll figure out the right target file path within `~/Documents/Clio/html/`, respecting the content's section and slug.
        *   Finally, we'll write the generated HTML content to its target file.

### Integration, Configuration & Testing

1.  **Integrate with `internal/feat/ssg/service.go`:**
    *   **Goal:** Potentially link our content database operations with the export process.
    *   **Action:** If it makes sense, we'll adapt our `CreateContent` and `UpdateContent` functions to trigger (or set up for triggering) the Markdown export processes in addition to their main job of saving to the database.
2.  **Centralized Path Configuration:**
    *   **Goal:** Keep our content paths consistent and flexible.
    *   **Action:** We'll set up a centralized configuration for defining the root content directory (`~/Documents/Clio/`), and make sure the mapping logic from database sections to filesystem paths (in both `md/` and `html/` directories) is consistent and configurable.
3.  **Testing:**
    *   **Goal:** Guarantee that our new features are correct and no regression are produced.
    *   **Action:** write thorough unit tests for individual components and integration tests to validate the entire content generation pipeline.


