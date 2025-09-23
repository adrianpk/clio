# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- Migrated direct Tailwind CSS classes and inline styles from `assets/ssg/layout/layout.html` and `assets/ssg/partial/list.tmpl` to `assets/static/css/prose.css`.
- Encapsulated Tailwind utility classes into custom CSS classes within `prose.css` using `@apply` directives.
- Updated `layout.html` and `list.tmpl` to use the new custom CSS classes from `prose.css`.
- Maintained Tailwind CSS CDN import in `layout.html` as per user request.

### Fixed
- Resolved build errors related to Goldmark `ast.Kind` constants by correctly aliasing `github.com/yuin/goldmark/ast` as `gmast` and `github.com/yuin/goldmark/extension/ast` as `extast`.
- Updated `renderEmphasis` to correctly handle strong emphasis (level 2) and removed the redundant `renderStrong` function.
- Ensured all rendering functions correctly use `gmast.Node`, `gmast.WalkStatus`, and `gmast.Text` where appropriate.

### Added
- **Preview Server:** A new web server running on port 8082 by default now serves the generated static site from `_workspace/documents/html`. This allows for a more realistic preview of the site with correct asset paths, without requiring a server restart for content changes.
- **SSG:** Implement a configurable limit (`CLIO_SSG_BLOCKS_MAXITEMS`) for the maximum number of items displayed in content blocks.
- **SSG:** Enforce a cascading hierarchy in block generation to ensure content appears only in the most relevant block.

## [2025-09-20]

### Added
- Implemented HTML generation from Markdown, rendering content within a template layout to create full pages.
- Added a dynamic navigation menu to the layout, generated from site sections.
- Implemented an asset pipeline for the static site generator:
    - Copies the embedded placeholder header image to the output directory.
    - Handles post-specific header images, copying them to the correct per-post directory.
    - Generates relative paths for assets to ensure links work on both local filesystems and web servers.
- Added a global configuration (`ssg.header.style`) to control the header layout style.
- Added support for four header styles: `stacked` (default), `overlay`, `text-only`, and `boxed` (which uses a full-width frosted-glass effect).

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