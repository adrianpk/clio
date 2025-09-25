# Clio → GitHub Pages Publishing

## Goal
Enable Clio to publish locally generated static HTML to GitHub Pages.

## High-level
- Build: Clio generates HTML into a dedicated `html/` folder.
- Publish: Everything under `html/` is pushed to a user-defined GitHub repository.

## Authentication options
- Option A: GitHub Personal Access Token (preferred).
  - Transport: HTTPS.
  - Mechanism: Send `Authorization: token <PAT>` header on each GitHub API request. For git push, use HTTPS with embedded credentials or an authenticated transport.
  - Scope: Minimum required to push to a single repo. Keep granular. User-supplied.
  - Storage: The token must be stored securely (e.g., as an environment variable). It must never be committed to the repository or exposed in logs.
- Option B: SSH.
  - Transport: SSH for git operations.
  - Assumption: User already has an SSH key pair and their agent set up on their OS.

## Config and environment
- Support both environment variables and config values, consistent with existing system conventions.
- Tentative names, adjust to project standards:
  - `CLIO_PUB_GH_REPO`
  - `CLIO_PUB_GH_BRANCH_PUB`
  - `CLIO_PUB_GH_AUTH_METHOD`
  - `CLIO_PUB_GH_TOKEN`
- For SSH, rely on the platform’s SSH agent. No key storage in app.

## Branching strategy options
- **Option 1: Direct Publish.** The app pushes the generated content directly to the branch configured for GitHub Pages (e.g., `gh-pages`).
  - Fast and ergonomic.
  - This is the default preference for v1 to keep the UX simple.
- **Option 2: Indirect Publish.** The app pushes the content to an alternative branch (e.g., `clio-publish`). The user must then manually merge or rebase these changes into the publishing branch via the GitHub UI.
  - Adds more control for versioning and gated releases, but introduces bureaucracy.
  - To be considered as an advanced mode in a future version.

## Default stance for v1
- Choose ergonomic simplicity.
- Publish directly to the Pages branch.
- If the user wants advanced versioning, they can create branches in GitHub themselves from the current published state.

## Publishing Client in Go

### Interface
The client will be defined by an interface to allow for mocking in tests. It will live in the `ssg` package.

```go
package ssg

type PublishAuthMethod int

const (
    AuthToken PublishAuthMethod = iota
    AuthSSH
)

type PublishConfig struct {
    Repo            string // "owner/repo"
    PublishBranch   string // e.g. "gh-pages" or "main"
    PagesSubdir     string // "" for root, or "docs"
    AuthMethod      PublishAuthMethod
    Token           string // if AuthToken
}

type Publisher interface {
    // Validate configuration and environment. Return detailed errors.
    Validate(cfg PublishConfig) error

    // Build the workspace content into the output HTML folder.
    // This method is environment-agnostic; it receives all necessary paths.
    BuildHTML() error

    // Publish the html/ folder to the configured repo and branch.
    // The commit message can be auto-generated.
    Publish(cfg PublishConfig, commitMessage string) error

    // Dry-run to show what would be pushed, without changing remote.
    Plan(cfg PublishConfig) (PlanReport, error)
}

type PlanReport struct {
    Added    []string
    Modified []string
    Removed  []string
    Summary  string
}
```

### Default Implementation
Resides in the `ssg` package as the default production publisher.
- **Token mode**: Use HTTPS remote with PAT via `Authorization` header for API calls. For `git push`, use HTTPS with embedded credentials or an authenticated transport.
- **SSH mode**: Use git over SSH, relying on the user’s SSH agent.
- **Both modes**:
    - Stage contents of `html/` into the repo worktree of the target branch or subdirectory.
    - Preserve dotfiles present in `html/` if any.
    - Avoid deleting unrelated files outside the publishing scope. If serving from `/docs`, operate only within `/docs`.
    - Write, commit, and push changes.

### Fake Implementation
A fake implementation of the `Publisher` will be created for testing purposes.
- It implements the `Publisher` interface using in-memory structures, avoiding network or filesystem operations.
- It records calls and inputs, allowing tests to verify the logic by simulating plan, build, and publish actions.

### Architecture: Publisher vs. Github Client
The system will have two distinct components: a high-level `Publisher` and a low-level `GithubClient`. The `Publisher` is responsible for the publishing logic, while the `GithubClient` provides a generic, reusable interface for interacting with GitHub repositories (e.g., clone, commit, push). The `Publisher` will *use* the `GithubClient` to perform its tasks. The `GithubClient` itself will be agnostic to the concept of 'publishing'.

## Build and Publish Workflow
The overall process is orchestrated by a higher-level component that understands workspaces (`test`/`production`). This component provides the environment-specific paths to the publishing client, which remains agnostic.

1.  **Generate HTML**: The client builds content into the specified `html/` output directory.
2.  **Prepare Target**:
    - If publishing to the repository root, map `html/*` to the root of the Pages branch.
    - If publishing to `/docs`, map `html/*` to the `/docs` subdirectory.
3.  **Plan (Dry-Run)**: If requested, stage changes and show a plan of what will be added, modified, or removed.
4.  **Commit**: Commit the changes with a standardized message (e.g., including a timestamp).
5.  **Push**: Push to the remote repository using the configured authentication method.
6.  **Report**: On success, confirm with the commit hash and a summary of changes.

## UX Considerations
- The default path should be a zero-config publish once the repository and authentication are set.
- Provide a visible dry-run “Plan” button for user confidence before publishing.
- Show a compact diff summary after the build and before the push.
- Advanced users can opt into a manual branch/PR flow in a future version.

## Error Handling
Provide clear, actionable errors:
- Invalid repository format.
- Branch not found or missing permissions.
- Token missing, has the wrong scope, or is expired.
- SSH agent not available or key not authorized.
- Build failures during Markdown-to-HTML generation.
- Detected divergence on the target branch that prevents a fast-forward push.
- Provide remediation hints where possible, without printing secrets.

## Logging
- Use structured logs with levels.
- Key events to log: build start/end, file counts, target repo/branch, and push success/failure.
- **Crucially, do not log secrets.**

## Security Notes
- Never persist the PAT unencrypted in project files.
- Do not echo tokens in the UI, logs, or error messages.
- Prefer environment variables or a system keychain if available for token storage.
- SSH private keys are managed by the OS and its agent, not by the application.

## Acceptance Criteria for v1
- Given a correct PAT or a working SSH agent, the app can push the `html/` directory to the configured branch, publishing the site.
- The default configuration uses token authentication and direct publishing.
- The `Publisher` interface is fully mockable and covered by tests for `Validate`, `BuildHTML`, `Plan`, and `Publish`.

## Open Questions (Deferred)
- Should the app assist in creating the target branch if it's missing?
- Should it auto-create the `/docs` folder and warn if Pages points elsewhere?
- Should a helper be added to validate repository Pages settings via the GitHub API?
- How to handle very large sites that might require chunked pushes?

---
*This document is an initial draft and will likely require further pruning and refinement.*
