# Menu Generation

This document outlines the architecture and design philosophy for constructing the site's navigation menus. It is closely related to the concepts described in `processor.md`.

## Design Philosophy

The menu design adheres to a minimalist, almost "vintage newspaper" aesthetic, but with a modern sensibility. The core principles are:

-   **Clarity over Density:** Avoid information overload. Menus should be clean, with ample whitespace, and present only essential navigation.
-   **Relaxed Experience:** The page structure is designed to be calming and focused. The user should feel relaxed, as if reading a high-quality digital book. The structure is strictly: Top Menu -> Main Content -> Bottom Navigation Blocks.
-   **Modern Typography:** While the spirit is vintage, the typography is modern and clean, consistent with the overall site design.

## Main Menu (First Line)

The main menu is the primary navigation bar at the top of the page. It is composed of the following elements in order:

1.  **Home:** A mandatory first link that always points to the site's root (`/`).

2.  **Sections:** The menu dynamically includes a link for each defined `Section`.
    -   **Ordering:** *Note:* Currently, the order of these sections is not guaranteed and will likely appear in the order they are retrieved from the database. A future Pull Request will introduce an `order` or `position` field to the `Section` model to allow for explicit manual ordering.

3.  **Contact (Optional):**
    -   A final `Contact` link can be displayed.
    -   **Activation:** This link is disabled by default and must be enabled via a hard-coded application configuration flag.
    -   **Content:** Initially, this will link to a static page whose content (e.g., email address, social links) is drawn from the application configuration. In the future, this could be powered by a dynamic `Page` content type. Any forms (like a contact form) would be handled via embedded third-party services, as Clio itself does not provide backend form processing.

## Submenu (Second Line)

A secondary, more subtle menu appears below the main menu under specific conditions.

-   **Trigger:** The submenu is displayed only if the current section (including the `root` section) contains content of type `Blog` or `Series`.
-   **Aesthetics:** The submenu is visually lighter and less prominent than the main menu, acting as a contextual navigation aid.
-   **Links:**
    -   **Blog:** If blog posts exist in the section, a `Blog` link will point to the index page for that section's blog (e.g., `/{section-name}/blog`).
    -   **Series:** If series exist in the section, a `Series` link will point to an index page listing all available series within that section.
