# Environment Variables and Configuration Keys Summary

This document summarizes the environment variables and configuration keys identified in the project, originating from `.envrc`, `makefile`, and `internal/am/key.go`.

## Direct Environment Variables (`.envrc` and `makefile`)

These variables are directly exported as environment variables, following the `CLIO_XXX_YYY_ZZZ` pattern.

*   **`CLIO_APP_ENV`**: Application environment (e.g., `dev`).
*   **`CLIO_SERVER_WEB_HOST`**: Web server host (e.g., `localhost`).
*   **`CLIO_SERVER_WEB_PORT`**: Web server port (e.g., `8080`).
*   **`CLIO_SERVER_API_HOST`**: API server host (e.g., `localhost`).
*   **`CLIO_SERVER_API_PORT`**: API server port (e.g., `8081`).
*   **`CLIO_SERVER_INDEX_ENABLED`**: Enables/disables the server index (e.g., `true`).
*   **`CLIO_DB_SQLITE_DSN`**: Connection string for the SQLite database (e.g., `"file:clio.db?cache=shared&mode=rwc"`).
*   **`CLIO_SEC_CSRF_KEY`**: Key for CSRF protection.
*   **`CLIO_SEC_ENCRYPTION_KEY`**: Encryption key.
*   **`CLIO_SEC_HASH_KEY`**: Hash key.
*   **`CLIO_SEC_BLOCK_KEY`**: Block key.
*   **`CLIO_NOTIFICATION_SUCCESS_STYLE`**: CSS style for success notifications.
*   **`CLIO_NOTIFICATION_INFO_STYLE`**: CSS style for info notifications.
*   **`CLIO_NOTIFICATION_WARN_STYLE`**: CSS style for warning notifications.
*   **`CLIO_NOTIFICATION_ERROR_STYLE`**: CSS style for error notifications.
*   **`CLIO_NOTIFICATION_DEBUG_STYLE`**: CSS style for debug notifications.
*   **`CLIO_BUTTON_STYLE_GRAY`**: CSS style for gray button.
*   **`CLIO_BUTTON_STYLE_BLUE`**: CSS style for blue button.
*   **`CLIO_BUTTON_STYLE_RED`**: CSS style for red button.
*   **`CLIO_BUTTON_STYLE_GREEN`**: CSS style for green button.
*   **`CLIO_BUTTON_STYLE_YELLOW`**: CSS style for yellow button.
*   **`CLIO_RENDER_WEB_ERRORS`**: Enables/disables web error rendering.
*   **`CLIO_RENDER_API_ERRORS`**: Enables/disables API error rendering.
*   **`CLIO_SSG_BLOCKS_MAXITEMS`**: Maximum number of items in SSG blocks.
*   **`CLIO_SSG_INDEX_MAXITEMS`**: Maximum number of items in the SSG index.
*   **`CLIO_SSG_SEARCH_GOOGLE_ENABLED`**: Enables/disables Google search in SSG.
*   **`CLIO_SSG_SEARCH_GOOGLE_ID`**: Google search ID for SSG.

## Web-Configurable SSG Parameters

The following SSG parameters are intended to be configurable via the web interface. They are stored in the `params` table and can be modified by users at runtime.

- **`ssg.blocks.maxitems`**: Maximum number of items in SSG blocks.
- **`ssg.index.maxitems`**: Maximum number of items in the SSG index.
- **`ssg.search.google.enabled`**: Enables/disables Google search in SSG.
- **`ssg.search.google.id`**: Google search ID for SSG.
- **`ssg.publish.repo.url`**: The URL of the repository where the site will be published (e.g., `git@github.com:user/repo.git`).
- **`ssg.publish.branch`**: The branch to which the site will be published (e.g., `gh-pages`).
- **`ssg.publish.pages.subdir`**: The subdirectory within the branch where the site will be published (e.g., `/`).
- **`ssg.publish.auth.method`**: The authentication method to use for publishing (e.g., `token`).
- **`ssg.publish.auth.token`**: The authentication token to use for publishing.
- **`ssg.publish.commit.user.name`**: The name of the user to use for the commit.
- **`ssg.publish.commit.user.email`**: The email of the user to use for the commit.
- **`ssg.publish.commit.message`**: The default commit message to use when publishing.


### Content Versioning (Future)

These parameters will be used for versioning content in a separate Git repository.

- **`ssg.content.repo.url`**: The repository URL for storing and versioning the markdown content.
- **`ssg.content.branch`**: The branch in the content repository.

## Configuration Keys in `internal/am/key.go`

This file defines a `Keys` struct and a `Key` variable containing configuration keys in `xxx.yyy.zzz` format. These `CLIO_XXX_YYY_ZZZ` environment variables are read and mapped to the internal `xxx.yyy.zzz` properties somewhere in the code (likely within the `am` package during configuration loading).

Here are the defined keys, along with their environment variable equivalent:

*   `CLIO_APP_ENV` => `app.env`
*   `CLIO_SERVER_WEB_HOST` => `server.web.host`
*   `CLIO_SERVER_WEB_PORT` => `server.web.port`
*   `CLIO_SERVER_WEB_ENABLED` => `server.web.enabled`
*   `CLIO_SERVER_API_HOST` => `server.api.host`
*   `CLIO_SERVER_API_PORT` => `server.api.port`
*   `CLIO_SERVER_API_ENABLED` => `server.api.enabled`
*   `CLIO_SERVER_RES_PATH` => `server.res.path`
*   `CLIO_SERVER_INDEX_ENABLED` => `server.index.enabled`
*   `CLIO_SERVER_PREVIEW_HOST` => `server.preview.host`
*   `CLIO_SERVER_PREVIEW_PORT` => `server.preview.port`
*   `CLIO_SERVER_PREVIEW_ENABLED` => `server.preview.enabled`
*   `CLIO_DB_SQLITE_DSN` => `db.sqlite.dsn`
*   `CLIO_SEC_CSRF_KEY` => `sec.csrf.key`
*   `CLIO_SEC_CSRF_REDIRECT` => `sec.csrf.redirect`
*   `CLIO_SEC_ENCRYPTION_KEY` => `sec.encryption.key`
*   `CLIO_SEC_HASH_KEY` => `sec.hash.key`
*   `CLIO_SEC_BLOCK_KEY` => `sec.block.key`
*   `CLIO_SEC_BYPASS_AUTH` => `sec.bypass.auth`
*   `CLIO_BUTTON_STYLE_GRAY` => `button.style.gray`
*   `CLIO_BUTTON_STYLE_BLUE` => `button.style.blue`
*   `CLIO_BUTTON_STYLE_RED` => `button.style.red`
*   `CLIO_BUTTON_STYLE_GREEN` => `button.style.green`
*   `CLIO_BUTTON_STYLE_YELLOW` => `button.style.yellow`
*   `CLIO_NOTIFICATION_SUCCESS_STYLE` => `notification.success.style`
*   `CLIO_NOTIFICATION_INFO_STYLE` => `notification.info.style`
*   `CLIO_NOTIFICATION_WARN_STYLE` => `notification.warn.style`
*   `CLIO_NOTIFICATION_ERROR_STYLE` => `notification.error.style`
*   `CLIO_NOTIFICATION_DEBUG_STYLE` => `notification.debug.style`
*   `CLIO_RENDER_WEB_ERRORS` => `render.web.errors`
*   `CLIO_RENDER_API_ERRORS` => `render.api.errors`
*   `CLIO_SSG_WORKSPACE_PATH` => `ssg.workspace.path`
*   `CLIO_SSG_DOCS_PATH` => `ssg.docs.path`
*   `CLIO_SSG_MARKDOWN_PATH` => `ssg.markdown.path`
*   `CLIO_SSG_HTML_PATH` => `ssg.html.path`
*   `CLIO_SSG_LAYOUT_PATH` => `ssg.layout.path`
*   `CLIO_SSG_HEADER_STYLE` => `ssg.header.style`
*   `CLIO_SSG_ASSETS_PATH` => `ssg.assets.path`
*   `CLIO_SSG_IMAGES_PATH` => `ssg.images.path`
*   `CLIO_SSG_BLOCKS_MAXITEMS` => `ssg.blocks.maxitems`
*   `CLIO_SSG_INDEX_MAXITEMS` => `ssg.index.maxitems`
*   `CLIO_SSG_SEARCH_GOOGLE_ENABLED` => `ssg.search.google.enabled`
*   `CLIO_SSG_SEARCH_GOOGLE_ID` => `ssg.search.google.id`
*   `CLIO_SSG_PUBLISH_REPO_URL` => `ssg.publish.repo.url`
*   `CLIO_SSG_PUBLISH_BRANCH` => `ssg.publish.branch`
*   `CLIO_SSG_PUBLISH_PAGES_SUBDIR` => `ssg.publish.pages.subdir`
*   `CLIO_SSG_PUBLISH_AUTH_METHOD` => `ssg.publish.auth.method`
*   `CLIO_SSG_PUBLISH_AUTH_TOKEN` => `ssg.publish.auth.token`
*   `CLIO_SSG_PUBLISH_COMMIT_USER_NAME` => `ssg.publish.commit.user.name`
*   `CLIO_SSG_PUBLISH_COMMIT_USER_EMAIL` => `ssg.publish.commit.user.email`
*   `CLIO_SSG_PUBLISH_COMMIT_MESSAGE` => `ssg.publish.commit.message`
