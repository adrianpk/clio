# Content Metadata

This document outlines the comprehensive metadata fields associated with a piece of content. These fields provide information for presentation, SEO, and content management.

## Core Identity & Content

Fields that define the content's address and provide summaries.

- **slug**: The URL-friendly identifier for the content, used in the path.
- **permalink**: The full, absolute, and permanent URL for the content.
- **excerpt**: A short, direct extract from the content body, typically used for previews in list views.
- **summary**: A self-contained summary of the content, which may be written independently and is suitable for external use (e.g., social media).

## Presentation & Layout

Fields that control the visual presentation of the content.

- **image**: The main featured image for the content.
- **social-image**: A specific image to be used when the content is shared on social media platforms.
- **layout**: Specifies the template file to be used for rendering the content page.
- **table-of-contents**: A boolean flag to control the automatic generation of a table of contents from headings in the body.

## SEO & Discovery

Fields that provide information to search engines and other crawlers.

- **description**: A short description of the content, primarily for use in search engine result snippets.
- **keywords**: A list of keywords relevant to the content.
- **robots**: Instructions for search engine crawlers (e.g., `noindex, nofollow`).
- **canonical-url**: The preferred URL for the content, used to combat duplicate content issues.
- **sitemap**: Configuration for how the content should be treated in the sitemap (e.g., priority, change frequency).

## Features & Engagement

Fields that enable or disable specific on-page features.

- **featured**: A boolean flag to mark the content as "featured," allowing it to be highlighted on the site.
- **share**: A boolean flag to control the visibility of social sharing buttons.
- **comments**: A boolean flag to enable or disable the comments section.

## Localization

Fields related to internationalization.

- **locale**: The language and region code for the content (e.g., `en-US`, `es-AR`).

## Timestamps & Status

Auditing, publication dates, and status for the content.

- **draft**: A boolean flag to indicate if the content is a draft. If true, the content will not be published, regardless of the `published-at` date.
- **published-at**: The date and time the content is officially published and made public.
- **created-at**: The date and time the content was originally created.
- **updated-at**: The date and time the content was last modified.

## Architecture

This section outlines the implementation plan for the metadata feature.

- **Metadata Form:** A "Meta" button in the content editor will trigger a full-screen modal to edit all metadata fields. This approach avoids cluttering the main editor interface.

- **Data Flow:** The modal will submit the metadata to a dedicated API handler. The handler will process and persist the data.

- **Database:**
  - A new `meta` table will be created to store the metadata fields.
  - This table will have a one-to-one relationship with the `content` table.
  - A new database migration will be created to reflect this change.

- **Data Access:**
  - The `GetAllContentWithTags` function will be renamed to `GetAllContentWithMeta`.
  - This function will fetch a content item along with all its related data, including `Section`, `Tags`, and `Meta`.
  - The new name will be used consistently across the data access layer, services, and handlers.

- **Markdown Generation:** The static site generator will query the content with its metadata and inject all fields into the frontmatter of the corresponding markdown file.