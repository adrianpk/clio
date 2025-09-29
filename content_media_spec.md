# Content & Media Specification

## Images
### Uploading
- Images can be uploaded through:
  - A **modal uploader** directly from the article editor.
  - A standalone image management interface (CRUD).
- After upload, images are displayed in the modal with:
  - **Thumbnail preview**.
  - **Clickable name/icon** to insert the image tag directly into the Markdown editor.

### Naming Convention
- All uploaded images are renamed automatically.
- Format:  
  ```
  {article-slug}_{timestamp}_{sequence}.{extension}
  ```
  - **article-slug** → unique identifier of the article.
  - **timestamp** → current time for uniqueness.
  - **sequence** → incremental number (001, 002, …).
  - **extension** → original file extension (`.jpg`, `.png`, `.webp`, etc).

### Header Image
- Defined outside of the Markdown body (article-level metadata).
- Naming format:  
  ```
  {article-slug}_header_{timestamp}.{extension}
  ```
- When detected, the system associates it automatically as the **article header image**.

### Content Images
- Inserted directly inside the Markdown body.
- Rendered as `<img>` tags in the final HTML.
- No special metadata needed, just standard Markdown image syntax:
  ```
  ![alt text](/images/{article-slug}_{timestamp}_{sequence}.jpg)
  ```

## Integration
- **Header image** = global, article-level metadata.
- **Content images** = inline, part of Markdown.
- Both stored consistently with the same renaming rules.

## UX Behavior
- Image upload, preview, and Markdown insertion are handled dynamically via **HTMX**.
- No full page reloads; updates happen seamlessly in the editor.


# Sections & Blog Headers

## Section Headers
- Managed via a dedicated section interface.
- Each section can have **only one header image** at a time.
- Uploading a new header replaces the previous one.
- Naming convention:  
  ```
  {section-slug}_header_{timestamp}.{extension}
  ```

## Blog Headers (per Section)
- Each section with a blog (`/section-name/blog`) can also define a **blog header image**.
- Each section can have **only one blog header image** at a time.
- Uploading a new blog header replaces the previous one.
- Naming convention:  
  ```
  {section-slug}_blog_{timestamp}.{extension}
  ```

## Uploading & UX
- Upload happens via the same interface as section headers.
- User selects whether the uploaded image is:
  - The **section header**, or
  - The **blog header** for that section.
- Only one image is stored per type; new uploads overwrite the previous one.


# Technical Architecture

## Directory Structure
Images follow the content hierarchy structure:

```
images/
├── article-slug-123/                    # Content en section root (/)
│   ├── article-slug-123_header_timestamp.jpg
│   ├── article-slug-123_timestamp_001.jpg
│   └── article-slug-123_timestamp_002.jpg
├── news/
│   ├── breaking-news-456/               # Content en section "news"
│   │   ├── breaking-news-456_header_timestamp.jpg
│   │   └── breaking-news-456_timestamp_001.jpg
│   └── section_header_timestamp.jpg     # Section header de "news"
├── blog/
│   ├── my-post-789/                     # Blog post en root blog
│   │   └── my-post-789_timestamp_001.jpg
│   └── blog_header_timestamp.jpg        # Blog header de root
└── tech/
    ├── tech-article-101/                # Content en section "tech"
    ├── blog/
    │   ├── tech-blog-post-202/          # Blog post en tech blog
    │   └── blog_header_timestamp.jpg    # Blog header de tech
    └── section_header_timestamp.jpg     # Section header de "tech"
```

### Directory Rules:
- **Root content**: `/images/{content.Slug()}/`
- **Section content**: `/images/{section.Path}/{content.Slug()}/`
- **Root blog**: `/images/blog/{content.Slug()}/`
- **Section blog**: `/images/{section.Path}/blog/{content.Slug()}/`
- **Section headers**: `/images/{section.Path}/section_header_timestamp.ext`
- **Blog headers**: `/images/{section.Path}/blog_header_timestamp.ext`

## ImageManager Component

### Location
`internal/feat/ssg/imagemanager.go`

### Responsibilities
- **File naming/renaming**: Generate consistent names following conventions
- **Directory management**: Create/manage directory structure
- **Metadata processing**: Extract and store image metadata
- **File operations**: Upload, move, delete operations
- **Quality management**: (Future) Compression, format conversion
- **Variant generation**: (Future) Thumbnails, different sizes

### Architecture Pattern
- Embeds `am.Core` for logging and config access
- Constructor: `NewImageManager(opts ...am.Option)`
- Service delegation pattern: Service calls ImageManager methods
- No direct database access - returns data for Service to persist

### Example Usage
```go
imageManager := NewImageManager()
result, err := imageManager.ProcessUpload(file, content, imageType)
// Service persists result to database
```

## Database Strategy

### Content Images
- **No foreign keys**: Use naming convention for association
- Store image paths/metadata in Content model if needed
- Header image path stored as Content metadata

### Section Images
- Consider foreign key relationships for easier management
- Section.HeaderImagePath field
- Section.BlogHeaderImagePath field

### Cleanup Policy
- **Flag-controlled**: `CLEANUP_ORPHANED_IMAGES = false` (constant)
- Future: Zombie image detection and cleanup utility
- Manual cleanup for now, automated later

## Implementation Scope

### This PR includes:
1. ImageManager component creation
2. Modal uploader with HTMX
3. Content editor integration
4. Section editor integration
5. Basic upload/naming functionality

### Future PRs:
1. Image variants generation
2. Advanced cropping/quality controls
3. Social media optimization
4. Automated cleanup utilities
5. Bulk image management

## Frontend Integration

### Modal Uploader
- HTMX-powered for seamless UX
- Thumbnail previews
- Click-to-insert functionality
- Progress indicators
- Error handling

### Editor Integration
- Content editor: Header + content images
- Section editor: Section + blog headers
- Markdown insertion helpers
- Preview functionality
