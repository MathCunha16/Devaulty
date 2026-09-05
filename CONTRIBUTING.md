# Contributing to Devaulty

Thank you for your interest in contributing to Devaulty. Contributions are
welcome through issues and pull requests that improve the application, its
security, documentation, or developer experience.

## Read the project documentation first

Before contributing, read the documentation in [`/docs`](docs/). It contains
important project rules, architectural decisions, security guidance, and
conventions that contributors are expected to follow.

At minimum, review:

- [`docs/architecture/`](docs/architecture/) for architectural decisions and
  project conventions;
- [`docs/security/`](docs/security/) for security and sensitive-data guidance;
- the relevant [`backend/README.md`](backend/README.md) or
  [`frontend/README.md`](frontend/README.md) for the area being changed.

These documents are part of the project's technical record. If a contribution
changes an existing decision, convention, or security behavior, explain the
reason in the pull request and update the relevant documentation.

## Before opening an issue

- Search existing issues to avoid duplicates.
- Use the bug report template for reproducible bugs.
- Do not include passwords, API keys, private keys, vault contents, database
  files, personal data, or other sensitive information.
- Do not report security vulnerabilities in a public issue. Use GitHub's
  private vulnerability reporting when enabled, or contact the maintainers
  privately before disclosing the issue.

## Before opening a pull request

1. Open or identify an issue for non-trivial changes.
2. Keep the pull request focused on one problem or feature.
3. Explain the motivation, implementation, and validation performed.
4. Update the relevant documentation when behavior, configuration, or
   architecture changes.
5. Add or update tests when the change affects behavior.
6. Do not commit generated artifacts, local databases, IDE metadata, secrets,
   or environment-specific files.

## Development expectations

- Follow the existing structure, naming, formatting, and patterns in the
  repository.
- Prefer small, understandable changes over speculative refactors.
- Preserve Devaulty's local-first and security-focused behavior.
- Treat credentials and vault data as sensitive at all times.
- A pull request must pass the repository's existing CI checks before merging.

## AI-assisted development

AI tools may be used as development aids, but blindly submitting generated
code is not acceptable. Do not submit "vibe-coded" changes that you cannot
explain, test, maintain, or defend.

Contributors are responsible for:

- reviewing every AI-generated suggestion before committing it;
- understanding the security and behavioral impact of the submitted code;
- verifying dependencies, licenses, and copied snippets;
- writing or running appropriate tests and checks;
- disclosing meaningful AI assistance in the pull request when it affected the
  design or implementation.

AI assistance does not replace code review, testing, or responsibility for the
contribution.

## Commit and pull request guidance

- Use a clear, imperative description for commits.
- Keep commits logically organized.
- Do not rewrite shared history after review has started unless coordinated
  with the maintainers.
- Be respectful and constructive when reviewing or discussing contributions.

## Questions

If you are unsure whether a change fits the project, open a discussion or an
issue describing the goal before investing in a large implementation.
