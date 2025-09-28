# Image (Asset Registry) Draft

This document outlines the proposed structure and considerations for an image asset registry, serving as a scaffolding for its first draft implementation.

## Identity
- **`id`**: UUID, primary key.
- **`short_code`**: Human-readable identifier.
- **`content_hash`**: SHA-256, unique, for deduplication and integrity.

## Technical Metadata
- **`mime`**: MIME type (e.g., `image/jpeg`).
- **`width`**: Image width in pixels.
- **`height`**: Image height in pixels.
- **`filesize_bytes`**: File size in bytes.
- **Deferred**: `raw_exif`, `raw_iptc`, `color_profile`, `dominant_color`, `palette` are deferred for a later iteration.

## Accessibility
- **`title`**: Image title (TEXT, nullable).
- **`alt_text`**: Alternative text for accessibility (TEXT, nullable).
- **`alt_lang`**: Language of alt_text (BCP-47, TEXT, nullable).
- **`long_description`**: Extended description for complex images (TEXT, nullable).
- **`caption`**: Image caption (TEXT, nullable).
- **`decorative`**: Boolean, true if the image is purely decorative and `alt_text` is not required.
- **`described_by_id`**: UUID, external reference to another image's `long_description` or an external resource.

## Rights & Provenance
- **Deferred**: `source_url`, `credit`, `copyright_notice`, `license`, `license_url` are deferred for a later iteration.

## Management
- **`etag`**: Derived from `content_hash`, for caching and integrity.
- **`created_at`**: Timestamp of creation.
- **`updated_at`**: Timestamp of last update.
- **`created_by`**: UUID of user who created the image.
- **`updated_by`**: UUID of user who last updated the image.

## Awareness
- This table (`images`) should remain unaware of usage (no FKs to posts, layouts, etc.). This ensures domain-agnosticism.

---

## Image Variants

- **`image_variants` table**: Stores different renditions of an image.
- **`image_id`**: UUID, foreign key to `images.id`.
- **`kind`**: Enum/TEXT in `{original, web, thumb}`.
- **`width`**: Variant width in pixels.
- **`height`**: Variant height in pixels.
- **`filesize_bytes`**: Variant file size in bytes.
- **`mime`**: Variant MIME type.
- **`blob_ref`**: Abstract reference to the stored file (e.g., S3 key, local path).

### Rules for Variants
- One variant per (`image_id`, `kind`).
- The `original` variant is never recompressed.

---

## External References (How other domains link to images)

- **Layouts**: Store `header_image_id` pointing to `images.id`.
- **Posts**: Integration with posts (e.g., via a `post_images` join table) is postponed for a later iteration.
- **Other Domains**: Manage their own references via FKs to `images.id`.

---

## Rules & Invariants

- **Idempotency**: Deduplicate images by `content_hash`.
- **Accessibility**: `alt_text` required unless `decorative = true`. Extended descriptions in `long_description` or via `described_by_id`, consistent with WCAG/WAI-ARIA.
- **Deletion**: Only allowed if no external references exist; enforced via `ON DELETE RESTRICT` (for FKs) and verified by a garbage collection process.
- **Updates**: Editing metadata in `images` should not affect consuming domains, as they only store `image_id`.

---

## Considerations (Prioritization)

- **Accessibility fields**: Included in the first iteration for compliance.
- **Deferred Metadata**: `raw_exif`, `raw_iptc`, `color_profile`, `dominant_color`, `palette`, `source_url`, `credit`, `copyright_notice`, `license`, `license_url` are deferred for a later iteration.
- **Posts Integration**: Postponed for a later iteration.
- **Domain-agnostic registry**: Avoids entanglement; consuming domains decide relationships.
- **Garbage collection**: Optional job, not critical for early drafts.
- **`etag` from hash**: Provides cache busting and stability.
