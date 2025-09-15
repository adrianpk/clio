# Tagging System

This document outlines the architecture and implementation of the content tagging system. Tagging allows for flexible content organization and discovery by associating keywords with content items.

## Implementation Details

- A `tag` table stores tag definitions (ID, name).
- A `content_tag` table links content and tags (`content_id`, `tag_id`), establishing a many-to-many relationship.
- The system supports dynamic tag creation. If a tag doesn't exist when being assigned to a content item, it's created automatically.
- The UI for content editing includes a field for managing tags, allowing users to add or remove them.
- The API handles the logic of associating tags with content, including creating new tags on the fly.
