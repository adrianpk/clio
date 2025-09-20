# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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