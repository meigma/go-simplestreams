# Welcome to the Meigma Go Template

This repository was generated from `template-go`, the standard starter for Meigma Go projects.
It is meant to give new repositories a working baseline on day one: a small Go CLI, Moon task orchestration, pinned CI, dependency automation, repository security defaults, and a dormant release pipeline that can be enabled when the project needs it.

Delete this file after you finish the first-repository setup checklist below.
It is only here to orient the initial project owner.

## What This Template Provides

- A minimal Go module at `github.com/meigma/template-go`.
- A Cobra/Viper CLI skeleton under `cmd/template-go` and `internal/cli`.
- Moon tasks for `format`, `lint`, `build`, `test`, and `check`.
- `golangci-lint` wired through Proto and Moon.
- CI that delegates to `moon ci --summary minimal` with pinned actions and minimal token permissions.
- Dependabot coverage for GitHub Actions, Go modules, and the docs npm project.
- Docusaurus docs scaffolding under `docs/`.
- Repository settings for signed commits, squash-only merges, immutable releases, private vulnerability reporting, and protected tags.
- Disabled release workflows for Release Please, GoReleaser, checksums, SBOMs, and GitHub artifact attestations.
- A root `ghd.toml` package manifest so released binaries can be installed with `ghd`.

## How It Works

Moon is the main entrypoint for local development and CI:

```sh
moon run root:check
```

That aggregate check runs the Go formatter/linter/build/tests plus the docs typecheck and build.
The GitHub Actions CI workflow runs the same path through:

```sh
moon ci --summary minimal
```

The release machinery is intentionally present but disabled.
Workflows live in `.github/workflows.disabled/` so they cannot run by accident.
When a project is ready to release binaries, move the relevant workflow files into `.github/workflows/`, configure the release app credentials and tag-ruleset bypass, then run the dry-run workflow before publishing.

## First Setup Checklist

1. Rename the Go module:

   ```sh
   go mod edit -module github.com/meigma/YOUR_REPO
   ```

2. Rename the binary directory:

   ```sh
   mv cmd/template-go cmd/YOUR_BINARY
   ```

3. Replace template placeholders:

   ```sh
   rg "template-go|TEMPLATE_GO|github.com/meigma/template-go"
   ```

   Update Go imports, Moon metadata, GoReleaser config, Release Please config, `ghd.toml`, README text, docs text, and CLI environment variable prefixes.

4. Refresh module metadata:

   ```sh
   go mod tidy
   ```

5. Run the full local check:

   ```sh
   moon run root:check
   ```

6. Decide what to do with releases:

   - Keep `.github/workflows.disabled/` if this project may release binaries later.
   - Delete the disabled release files if this project will never publish releases.
   - Enable them only after reviewing `.goreleaser.yaml`, `release-please-config.json`, and the workflow comments.
   - Keep and update `ghd.toml` if the binary should be installable through `ghd`.

7. Update project-facing docs:

   - Rewrite `README.md` for the actual project.
   - Review `CONTRIBUTING.md` and `SECURITY.md`.
   - Add a real license before publishing the repository.

8. Delete this file:

   ```sh
   rm DELETE_ME.md
   ```
