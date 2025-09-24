# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **External Site Search Support**: Implemented Google Custom Search integration for site-wide content search, allowing users to easily find content.
- **Pagination Controls**: Implemented navigation controls (Previous/Next, Page X of Y) for index pages.

### Changed
- **Asset Separation**: Refactored SSG static assets into a dedicated `assets/ssg/static/` structure.
- **Asset Paths**: Updated SSG templates and Go code to use absolute and simplified asset paths (`/static/`) in generated HTML, removing the `/ssg` prefix.

### Fixed
- Resolved various internal issues to ensure stable and correct rendering of site content and features.

## [2025-09-23]

### Added
- **Section Indexes:** Implemented functionality to generate and display section indexes, providing organized listings of content within different categories.
- **Static Site Preview Server:** Added a dedicated server for previewing the generated static site, improving the development workflow by allowing real-time content and style verification.
- **SSG Enhancements:** Implemented a configurable limit (`CLIO_SSG_BLOCKS_MAXITEMS`) for the maximum number of items displayed in content blocks and enforced a cascading hierarchy in block generation.
- **Pagination Controls**: Implemented navigation controls (Previous/Next, Page X of Y) for SSG index pages.

### Improved
- **Styling and Consistency:** Enhanced the overall visual presentation by centralizing placeholder images for content headers and section indexes, and refining the display of content pages and section index cards. 

## [2025-09-20]

### Added
- Implemented HTML generation from Markdown, rendering content within a template layout to create full pages.
- Added a dynamic navigation menu to the layout, generated from site sections.
- Implemented an asset pipeline for the static site generator:
    - Copies the embedded placeholder header image to the output directory.
    - Handles post-specific header images, copying them to the correct per-post directory.
    - Generates relative paths for assets to ensure links work on both local filesystems and web servers.
- Added a global configuration (`ssg.header.style`) to control the header layout style.
- Added support for four header styles: `boxed` (default), `overlay`, `text-only`, and `stacked` (which uses a full-width frosted-glass effect).

### Docs
- Updated the gallery with screenshots and descriptions of the new header styles.


## [2025-09-18]

### Added
- Implemented a comprehensive metadata system for content.
- The web UI now includes a modal for managing various metadata fields, including publishing status, SEO attributes (description, keywords, robots), and content features (ToC, sharing, comments).
- The static site generator now marshals all metadata into YAML frontmatter for each generated markdown file.

## [2025-09-16]

### Added
- Implemented **Zen Mode** for the markdown editor, providing a fullscreen, distraction-free writing canvas.
- Implemented a **Dark Mode** for the editor, available only within Zen Mode.
- Added keyboard shortcuts for toggling Zen Mode (`Alt+Z`) and Dark Mode (`Alt+D`).
- Created a dual-button system: static buttons for entering Zen Mode and floating buttons for exiting and controlling Dark Mode.

### Changed
- Refactored editor enhancement logic into a single `editor-enhancements.js` file.
- Refined button positioning and styles for a cleaner user experience.