# Content and Asset Generation Strategy

This document details the configuration and conventions for generating content, sections, and their associated assets, particularly images.

## Section and Content-Type Structure

The site is organized into sections. A section can be the root (`/`) or have a specific name (e.g., `news`, `articles`). Within each section, different types of content are handled with specific path structures:

-   **Pages and Articles:** Rendered directly under their section. Example: `/articles/article-slug`.
-   **Blog Posts:** Rendered within a `blog` sub-path in their section. Example: `/news/blog/post-slug`.
-   **Series Posts:** Rendered within a sub-path named after the series. Example: `/articles/my-series-name/post-slug`.

## Asset and Image Location

Images associated with a piece of content are stored in a corresponding `img` subdirectory. If a post's final path is `/articles/my-article`, its images will be located in the output directory at `/articles/my-article/img/`.

-   **Header Image:** If a specific header image is provided for a post, it must be named `header.ext` (where `.ext` can be `png`, `jpg`, `webp`, etc.).
-   **Social Media Images:** Social media card images follow a specific naming convention, such as `social-fb.ext` or `social-x.ext`, depending on the platform. Multiple versions may exist to accommodate different aspect ratios.

## User-Uploaded Image Naming

Images uploaded by a user will be automatically renamed to include the content's slug, a timestamp, and a partial unique identifier.
-   **Example:** `my-article-1632512345-uuidsegment.ext`.

## Embedded Placeholder Image

A default placeholder image is embedded in the Go binary. This image is used as a fallback when no specific header image is provided for a piece of content.

-   **Conditional Copying:** The placeholder image is only copied to the final output directory if the configured header style requires an image. If the style does not use a header image, the placeholder is not copied, saving disk space.

## Processing and Final Generation

The generation processor is responsible for transforming Markdown and handling assets. It places the final images in their correct paths within the output directory, ensuring that all image links work correctly both when viewed on a local filesystem during development and when deployed to a web server in production.
