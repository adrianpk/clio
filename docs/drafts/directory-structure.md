# Directory Structure

This document defines the intended directory structure for the generated content and its purpose.

## Operating System Considerations

Initially, the focus is on a Linux-based environment as it is the primary development OS. Support for other operating systems will be explicitly and automatically handled later. However, the system will be configurable from the start to allow easy adaptation to any OS.

## Operating Modes

Clio can run in two modes, determined by the `CLIO_APP_ENV` environment variable.

### Production Mode (`prod`, Default)

This is the default mode when `CLIO_APP_ENV` is not set or is set to `prod`. It's intended for regular use after the application has been installed. It uses standard system-wide directories.

*   **Configuration:** `~/.config/clio`
*   **Database:** `~/.clio`
*   **Documents Root:** `~/Documents/Clio/`
    *   **Markdown Source:** `markdown/`
    *   **Generated HTML:** `html/`
    *   **Assets:** `assets/images/`

### Development Mode (`dev`)

This mode is activated by setting `CLIO_APP_ENV=dev`. It's designed for development and testing, keeping all files within the project's working directory.

*   **Base Path:** `./_workspace/`
*   **Configuration:** `./_workspace/config/`
*   **Database:** `./_workspace/db/`
*   **Documents Root:** `./_workspace/documents/`
