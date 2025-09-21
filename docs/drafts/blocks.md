# Bottom Navigation Blocks

To maintain the minimalist reading focus, the main content area is not cluttered with sidebars or intrusive widgets. Instead, contextual navigation and engagement elements are placed in "blocks" at the bottom of the page, after the main content.

These blocks serve as "engagement helpers" and can include:

-   **Related Content:** Links to other articles or pages that are relevant to the current one.
-   **Next/Previous:** Sequential navigation for content within a series or a chronological sequence.

The specific blocks to be displayed can be configured at the layout level, allowing for different sets of blocks for different types of content.

## Philosophy and Purpose

Blocks are primarily concerned with **navigation** and **engagement**.

It is worth noting that Clio's core philosophy does not focus on stimulating user interaction. However, this principle is not rigid, and specific setups can be configured to encourage more engagement as needed.

> *Note: The information in the following "Core Concepts" sections is partly repeated from other drafts. It is consolidated here to provide a unified context for the ideas described below.*

---

## Core Concepts: Sections

*   **Root Section:** A special section with the path `/`, which otherwise behaves like any other section.
*   **Section Path:** Every section has a name and a path, following the format `/{section-name}`.
*   **Flat Structure:** Sections cannot be nested.

---

## Core Concepts: Content Types and URLs

There are four main types of content, each with a specific purpose and URL structure.

### 1. Pages
- **Purpose:** Structural or general informational content. They provide a flexible container for different kinds of information and don't necessarily follow a narrative or editorial logic.
- **URL Structure:** Hangs directly from a section.
    - **Root Section:** `/slug/index.html`
    - **Other Sections:** `/{section-name}/slug/index.html`

### 2. Articles
- **Purpose:** Discursive content meant to be read as a complete piece (e.g., essays, commentary, narratives).
- **URL Structure:** Same as Pages. Hangs directly from a section.
    - **Root Section:** `/slug/index.html`
    - **Other Sections:** `/{section-name}/slug/index.html`

### 3. Blog
- **Purpose:** Chronologically organized posts.
- **URL Structure:** Hangs from a dedicated `/blog` path within a section.
    - **Root Section:** `/blog/slug/index.html`
    - **Other Sections:** `/{section-name}/blog/slug/index.html`

### 4. Series
- **Purpose:** Ordered posts that form a sequence (e.g., guides, tutorials). The system will provide automatic navigation between parts of the series.
- **URL Structure:** Hangs from a path defined by the series' name within a section.
    - **Root Section:** `/{series-name}/slug/index.html`
    - **Other Sections:** `/{section-name}/{series-name}/slug/index.html`

---

## Block Behavior by Content Type

This section defines which blocks are displayed for each type of content.

### 1. Pages
By default, **Pages do not display any blocks**.

Navigation away from a Page is handled by the main site menu or through links manually embedded within the page's content. This is because Pages are considered structural elements, often serving as editable portals to other parts of the site rather than pieces of a larger narrative.

### 2. Articles
Articles have four potential blocks, designed to show a progression of relevance to the reader. The goal is to signal how "strong" the connection is between the current article and the suggested content.

The blocks are:

1.  **Tag-Related Content (Same Section)**
    *   **Logic:** Displays links to articles *within the same section* that share common tags.
    *   **UI:** The matching tag(s) should be displayed as small labels to indicate the reason for the suggestion.
    *   **Relevance:** Strongest.

2.  **Recent Content (Same Section)**
    *   **Logic:** Displays links to the most recent articles *within the same section*, regardless of tags.
    *   **Relevance:** Medium.

3.  **Tag-Related Content (All Sections)**
    *   **Condition:** This block is only displayed if a specific "allow cross-section" flag is enabled in the configuration.
    *   **Logic:** Displays links to articles *from any section* that share common tags.
    *   **Relevance:** Weak.

4.  **Recent Content (All Sections)**
    *   **Condition:** Requires the same "allow cross-section" flag to be enabled.
    *   **Logic:** Displays links to the most recent articles *from any section*, regardless of tags.
    *   **Relevance:** Weakest.

### 3. Blog
Blog posts have two potential blocks, focused on discovery within the same blog (i.e., the blog belonging to the same section).

1.  **Tag-Related Posts (Same Blog)**
    *   **Logic:** Displays links to other posts *within the same blog* that share common tags.

2.  **Recent Posts (Same Blog)**
    *   **Logic:** Displays a configurable number of recent posts from the *same blog*, sorted from newest to oldest. The current post is excluded from this list.

### 4. Series
Series have two blocks, both focused on navigating the sequence of posts.

1.  **Simple Sequential Navigation**
    *   **Logic:** Provides simple "previous" and "next" links to move one step backward or forward in the series.
    *   **UI:** Renders as `<-` for previous and `->` for next.
    *   **Condition:** The `<-` link will not render on the first post, and the `->` link will not render on the last post.

2.  **Full Series Index**
    *   **Logic:** This block allows for non-sequential navigation by providing a complete, ordered list of all posts in the series, relative to the current one.
    *   **UI:** The block is divided into two parts:
        *   **Upcoming Posts:** A list of all subsequent posts in the series, ordered by progression. Each link is suffixed with a forward arrow (e.g., `Post Title ->`).
        *   **Previous Posts:** A list of all prior posts, ordered from the most recent to the oldest (i.e., from the post just before the current one back to the first post). Each link is prefixed with a backward arrow (e.g., `<- Post Title`).

---
## Implementation and Performance Notes

Taking a cue from Kent Beck’s line, “Make it work, make it right, make it fast,” the first version of the `BlockBuilder` is going to lean on clarity and correctness before worrying about speed.

- **Initial approach (make it work):** For every content item, the builder will scan through the full list of site content to find matches. That’s about **O(N²)** in terms of complexity. For Clio’s expected scale (personal blogs, small project sites), this is fine and won’t really affect build times.

- **Later optimization (make it fast):** If at some point we need to deal with a much larger number of posts, the path forward is straightforward. Instead of doing a full scan each time, we’d add a one-time indexing step at the start of the build (an O(N) cost). By mapping content by section, tag, and so on, the `BlockBuilder` can look up related items almost instantly. That gets us closer to **O(N)** overall.

This step-by-step way of doing it gives us a feature that works correctly today, with a clear option to scale performance when it actually matters.
