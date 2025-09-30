# Image Management System

This document details the image management system for content, sections, and static site generation.

## Architecture Overview

The system manages images through a database-driven approach with relationship tables, supporting content images, section headers, and blog headers with proper metadata and accessibility features.

### Database Schema

- **images**: Core image registry with metadata (alt_text, caption, dimensions, etc.)
- **image_variants**: Different sizes/formats of the same image
- **content_images**: Links images to content pieces
- **section_images**: Links images to sections (header, blog_header)

## Content Structure

### Section and Content-Type Organization

The site is organized into sections. A section can be the root (`/`) or have a specific name (e.g., `news`, `articles`). Within each section, different types of content are handled with specific path structures:

- **Pages and Articles:** Rendered directly under their section. Example: `/articles/article-slug`.
- **Blog Posts:** Rendered within a `blog` sub-path in their section. Example: `/news/blog/post-slug`.
- **Series Posts:** Rendered within a sub-path named after the series. Example: `/articles/my-series-name/post-slug`.

## Image Types and Naming

### Content Images
- Inserted directly inside the Markdown body
- Rendered as `<img>` tags in the final HTML
- Uploaded through modal uploader from content editor
- Standard Markdown image syntax: `![alt text](/images/{content-slug}_{timestamp}_{sequence}.jpg)`

### Header Images
- Content-level metadata (outside Markdown body)
- Section-level headers for section pages
- Blog headers for section blog pages
- Naming format: `{slug}_header_{timestamp}.{extension}`

### User-Uploaded Image Naming
Images uploaded by users are automatically renamed with consistent conventions:
- **Content images:** `{content-slug}_{timestamp}_{sequence}.{extension}`
- **Section headers:** `{section-slug}_header_{timestamp}.{extension}`
- **Blog headers:** `{section-slug}_blog_{timestamp}.{extension}`

## Asset Location Strategy

### Development/Upload Storage
Images are stored in the uploads directory during development and content management.

### Generated Site Output
For the final generated site, images are placed in corresponding subdirectories:
- If a post's final path is `/articles/my-article`, its images will be located at `/articles/my-article/img/`
- **Header Image:** Named `header.ext` in the content's img directory
- **Social Media Images:** Named `social-fb.ext`, `social-x.ext`, etc.

### Directory Structure Example
```
output/
├── articles/
│   └── my-article/
│       ├── index.html
│       └── img/
│           ├── header.jpg
│           ├── social-fb.jpg
│           └── content-image-001.jpg
└── news/
    ├── index.html
    ├── img/
    │   └── header.jpg          # Section header
    └── blog/
        ├── index.html
        ├── img/
        │   └── header.jpg      # Blog header
        └── post-slug/
            └── img/
                └── header.jpg  # Post header
```

## Technical Implementation

### ImageManager Component
Location: `internal/feat/ssg/imagemanager.go`

Responsibilities:
- File naming/renaming following conventions
- Directory management and structure creation
- Metadata extraction and processing
- File operations (upload, move, delete)
- Future: Quality management, variant generation

### Database Strategy
- **Content Images:** Associated through content_images relationship table
- **Section Images:** Associated through section_images relationship table with purpose field
- **Cleanup Policy:** Flag-controlled with future automated cleanup utilities

## UX and Frontend Integration

### Modal Uploader
- HTMX-powered for seamless experience
- Thumbnail previews with accessibility metadata
- Click-to-insert functionality for content images
- Progress indicators and error handling

### Editor Integration
- Content editor: Header + content images
- Section editor: Section + blog headers
- Markdown insertion helpers
- Real-time preview functionality

## Embedded Placeholder Image

A default placeholder image is embedded in the Go binary as fallback when no specific header image is provided.

- **Conditional Copying:** Only copied to output if the configured header style requires an image
- Saves disk space when header images are not used

## Processing and Final Generation

The generation processor transforms content and handles assets, placing final images in correct paths within the output directory. This ensures image links work correctly both during development and in production deployment.
