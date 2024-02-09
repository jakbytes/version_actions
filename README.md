# version_actions

## Warning: This is a work in progress, you can probably accomplish what you need with a combination of other github actions and bash scripts.

[![Go Report Card](https://goreportcard.com/badge/github.com/jakbytes/version_actions)](https://goreportcard.com/report/github.com/jakbytes/version_actions)
[![codecov](https://codecov.io/gh/jakbytes/version_actions/graph/badge.svg?token=QTT22V13C4)](https://codecov.io/gh/jakbytes/version_actions)

`version_actions` is a suite of composite GitHub Actions written in Go, designed to facilitate various semantic versioning activities within GitHub repositories. Leveraging conventional commits, these actions provide robust tools for version incrementing, changelog generation, tag creation, and automated pull request generation.

Each action is designed to do one task and do it well. This allows for the creation of workflows that can be easily customized to meet the needs of your project.

There are a number of other actions you may be interested in that accomplish similar tasks. You should evaluate each depending on your needs. Here are a few:
- [google-github-actions/release-please-action](https://github.com/google-github-actions/release-please-action): automated releases based on conventional commits
- [anothrNick/github-tag-action](https://github.com/anothrNick/github-tag-action): A GitHub Action to tag a repo on merge.


## Features

This action is currently in pre-release form, however the following features are planned for the initial release:

- **Version Incrementing**: Automatically increments repository version based on conventional commit messages and provide it as an output
- **Changelog Generation**: Creates comprehensive changelogs from conventional commit messages.
- **Tag Creation**: Facilitates the creation of new tags aligned with semantic versioning.
- **Pull Request Generation**: Automates the creation and setting of pull requests on pushes to arbitrary branches.
- **Generating Releases**: Autmomates the creation of releases 

## Usage

To use `version_actions` in your GitHub workflows, include them as steps in your `.yml` workflow files. [Examples](https://github.com/jakbytes/version_actions/blob/main/README.md#examples) for each action can be found below.

## Configuration

Each action (pull_request, version, release) may have its unique configuration options. All actions require minimally:

- token: A GitHub or Personal Access Token.
- action: The specific action to be executed (pull request, version, release).

## Workflows

### Pull Request

This workflow automates the creation of a pull request to merge a specified branch into the main branch, and includes an up to date changelog. For manual edits to the pull request, ensure they are made above the `## Changelog` header. As long as `## Changelog` is present in the body, the content below this header will be automatically removed and updated by the workflow.

This example creates a draft PR from any non-main branch to the main branch, if you have multiple release branches you may ignore other branches by adding them to the `branches-ignore` list.

```yaml
on:
  push:
    branches-ignore:
      - 'main'
jobs:
  semver_job:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Repository
        uses: actions/checkout@v4

      - name: Automated Pull Request
        id: semver_action
        uses: jakbytes/version-action/pull_request@v0.1.4
        with:
          token: ${{ secrets.GITHUB_TOKEN }}
          base: "main"
```

## Requirements

- Some workflows require a PERSONAL_ACCESS_TOKEN with specific permissions
- Conventional Commits: Adherence to conventional commit standards is required for effective changelog generation and version management. In the future we may support other basic commit messaging formats to handle versioning.

## [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/)

We use conventional commits because using a standardized set of rules to write commits, makes commits easier to read, and enforces writing descriptive commits.

### Why?

- Automatically generating CHANGELOGs.
- Automatically determining a semantic version bump (based on the types of commits landed).
- Communicating the nature of changes to teammates, the public, and other stakeholders.
- Triggering build, release, publish processes.
- Making it easier for people to contribute to your projects, by allowing them to explore a more structured commit history.

### Types

- `fix`: a commit of the type fix patches a bug in your codebase (this correlates with PATCH in Semantic Versioning).
- `feat`: a commit of the type feat introduces a new feature to the codebase (this correlates with MINOR in Semantic Versioning).
- `BREAKING CHANGE`: a commit that has a footer BREAKING CHANGE:, or appends a ! after the type/scope, introduces a breaking API change (correlating with MAJOR in Semantic Versioning). A BREAKING CHANGE can be part of commits of any type.
- **types other than fix: and feat**: are allowed, for example @commitlint/config-conventional (based on the Angular convention) recommends build:, chore:, ci:, docs:, style:, refactor:, perf:, test:, and others.
footers other than BREAKING CHANGE: <description> may be provided and follow a convention similar to git trailer format.

#### Changelog

In version_actions, commits are automatically categorized in the changelog under the following types, each with its own header:

- `fix`: Bug fixes, corresponding to PATCH in Semantic Versioning (SemVer).
- `feat`: New features, corresponding to MINOR in SemVer.
- `docs`: Changes exclusively to documentation.
- `style`: Code style changes (e.g., whitespace, formatting, missing semi-colons) without altering code functionality.
- `refactor`: Code changes that are neither bug fixes nor feature additions.
- `perf`: Enhancements that improve performance.
- `test`: Additions or corrections to existing tests.
- `build`: Modifications affecting the build system or external dependencies (examples: pip, docker, npm).
- `ci`: Changes to CI configuration files and scripts (examples: GitLabCI).

#### Tools

Here are some tools used in the creation and maintenance of this repository:

- [pre-commit](https://pre-commit.com/index.html): a tool that simplifies code quality checks by automatically running predefined hooks before each commit. 
- [Commitizen](https://commitizen-tools.github.io/commitizen/): a cli tool to generate conventional commits

## License
version_actions is MIT licensed, as found in the [LICENSE](https://github.com/jakbytes/version_actions/blob/main/LICENSE) file.

## Support and Issues
For support or to report issues, please create an issue in the GitHub repository.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
