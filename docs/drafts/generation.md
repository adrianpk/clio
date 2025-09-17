# Content Generation from Database

This document outlines the process of converting the content stored in the database into the actual documentation files.

*   **Single Source of Truth:** The database is and will be the single source of truth for all content.
*   **Reconstructability:** It must be possible to reconstruct the entire content, including its organization and categorization (tags), solely from the database.
*   **Database Engine:** The database is a single SQLite file.
*   **Exclusions:** Images will not be stored in the database.
*   **Versioning via Git:** Exporting content to files allows for versioning all content components together in Git/GitHub. This provides these key benefits:
    *   **Versioning:** Markdown files and their associated assets (e.g., images) are versioned together to maintain consistency and prevent mismatches.
    *   **Database:** The database only stores the latest version of a document's text and metadata. Git is responsible for the full version history, keeping the database lean.
    *   **Automation:** Clio will be responsible for committing all content changes (text and assets) to GitHub, based on user configuration (manual, automatic, or hybrid).
*   **Database Reconstruction:** Eventually, it will also be possible to reconstruct the database from a clone of the repository where the base docs are stored. This feature is not a priority for the initial implementation.
*   **Scope:** This process concerns the base Markdown content, not the final generated site that Clio will also version and publish.
