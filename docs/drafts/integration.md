# Integration Test Plan: Real GitHub Repository

This document outlines a detailed plan for implementing an integration test designed to interact with a real GitHub repository. This test is intended for **on-demand execution** rather than continuous integration (CI/CD) pipelines, allowing for repeatable validation of the full lifecycle without relying on constant network connectivity or GitHub service availability. Its "deactivatable" nature, combined with automated setup and teardown procedures, ensures a predictable initial state and facilitates easy, repeated testing without manual intervention.

## 1. Test Repository Definition

*   A dedicated private GitHub repository will be required for testing purposes.
*   The repository name and owner/organization will be configured via environment variables to ensure flexibility and avoid hardcoding.

## 2. Credential Configuration (Environment Variables)

*   The test will rely on environment variables for GitHub credentials (Personal Access Token) and test repository details (owner, repository name).
*   These variables will be loaded at the beginning of the test. If any are missing, the test will automatically be skipped (`t.Skip()`) with an informative message.

## 3. Test Structure (Table-Based)

*   The test will be table-based to allow for multiple publishing scenarios (e.g., first commit, content update, content deletion).
*   Each table entry will define:
    *   `name`: A descriptive name for the scenario.
    *   `initialContent`: The initial content of the file to be published (or `nil` if the file does not exist initially).
    *   `publishContent`: The content to be published in the scenario.
    *   `expectedCommitMessage`: The expected commit message.
    *   `expectedFileContent`: The expected final content of the file in the repository.
    *   `expectedError`: The expected error (or `nil` if no error is expected).

## 4. Test Setup (`TestMain` or Setup Function)

*   **Environment Variable Verification:** At the start, the presence of necessary environment variables (e.g., `GITHUB_TEST_REPO_OWNER`, `GITHUB_TEST_REPO_NAME`, `GITHUB_TEST_TOKEN`) will be checked. If any are missing, an informative message will be printed, and the test will be skipped (`t.Skip()`).
*   **GitHub Client Initialization:** An instance of `am.GitClient` will be created using the real GitHub client (`internal/am/github.NewClient`) with credentials obtained from environment variables.
*   **Initial Repository Cleanup:**
    *   The test repository will be cloned to a temporary local directory.
    *   All existing files and directories in the main branch (or designated test branch) of the cloned repository will be deleted, except for the `.git` directory.
    *   An empty commit (or a cleanup commit) will be performed, and a force push will be executed to the main branch to ensure a clean and predictable state before each test run.
    *   A temporary working branch will be created for the test (e.g., `test-branch-<timestamp>`) to isolate changes for each execution.

## 5. Test Teardown (`t.Cleanup()` Function)

*   **Working Branch Cleanup:** At the end of each sub-test (or at the end of `TestMain`), the temporary working branch created for the test will be deleted, both locally and in the remote repository.
*   **Temporary Directory Deletion:** The local temporary directory where the repository was cloned will be removed.

## 6. Logic for Each Sub-Test (within the table's `for` loop)

*   **Repository Preparation for the Scenario:**
    *   If `initialContent` is defined for the scenario, a file with that content will be created in the cloned temporary directory, committed, and pushed to the temporary working branch.
*   **`Publisher` Execution:**
    *   An instance of `ssg.Publisher` will be created with the real GitHub client.
    *   The `Publish` method of the `Publisher` will be called with the `publishContent` and other relevant parameters (e.g., `path`, `commitMessage`).
*   **Results Verification:**
    *   **Errors:** Verify if the error returned by `Publish` matches `expectedError`.
    *   **File Content:**
        *   The repository will be cloned again (or a `git pull` will be performed on the existing clone) to get the latest state.
        *   The content of the published file in the repository will be read.
        *   The read content will be compared with `expectedFileContent`.
    *   **Commit History:**
        *   The last commit message on the temporary working branch will be verified to ensure it matches `expectedCommitMessage`.

## 7. Additional Considerations

*   **Timeout:** Set a reasonable timeout for the test to prevent it from hanging due to network issues or GitHub API problems.
*   **Retries:** Consider implementing retries for network operations that might intermittently fail (e.g., cloning, pushing).
*   **Isolation:** Ensure that each sub-test is completely independent and does not affect other sub-tests.
*   **Logging:** Add detailed logging for debugging in case of failures.
*   **Test File Name:** `internal/feat/ssg/publisher_integration_test.go` to distinguish it from unit tests.
