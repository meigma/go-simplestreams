# template-go

`template-go` is the reusable Go repository starter for Meigma projects.
It includes a small Go CLI skeleton, Moon tasks, pinned CI, Dependabot, baseline repository security settings, and a dormant Release Please plus GoReleaser release layer.

## Local Bootstrap

Prerequisites:

- Go 1.26.2
- Moon 2.x
- Node.js 22.22.2 for the Docusaurus docs project

After creating a new repository from this template, replace the placeholder names before doing feature work:

```sh
go mod edit -module github.com/meigma/YOUR_REPO
mv cmd/template-go cmd/YOUR_BINARY
```

Then update `template-go` references in the Moon tasks, GoReleaser config, `ghd.toml`, README, and package docs.

## Common Tasks

Moon is the standard task front door:

```sh
moon run root:format
moon run root:lint
moon run root:build
moon run root:test
moon run root:check
```

CI runs the same aggregate check:

```sh
moon ci --summary minimal
```

The starter CLI is intentionally small:

```sh
go run ./cmd/template-go --version
go run ./cmd/template-go --message "hello from cobra"
go test ./...
```

The CLI entrypoint uses Cobra and Viper in the same shape as other Meigma CLIs: `cmd/template-go` stays thin, `internal/cli` owns command construction, and Viper-backed flags can also be supplied through `TEMPLATE_GO_*` environment variables.

## CI and Security

The default CI workflow keeps permissions minimal, pins external actions, disables checkout credential persistence, and delegates checks to Moon.
Dependabot covers GitHub Actions, the root Go module, and the docs npm project.

Repository settings live in `.github/repository-settings.toml`.
They default to immutable releases, private vulnerability reporting, signed commits, squash-only merges, and protected tags.

## Release Layer

Release automation is included but disabled by default.
Projects that release binaries can enable it by moving the templates from `.github/workflows.disabled/` into `.github/workflows/`, then configuring the release app credentials and tag-ruleset bypass documented in those files.

The release path is:

- Release Please opens and maintains the release PR.
- Release Please creates a draft GitHub release and tag after merge.
- GoReleaser builds binaries, checksums, and SBOMs without publishing directly.
- The release workflow uploads assets to the draft release and creates a GitHub-hosted attestation for `checksums.txt`.
- A human inspects the draft release before publication.

The root `ghd.toml` matches the default GoReleaser output so generated projects can be installed with `ghd` once the release workflow is enabled.
After cloning this template, update `provenance.signer_workflow`, package names, asset patterns, and binary paths to match the new repository and binary name.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines, local setup expectations, and pull request workflow.

## Security

See [SECURITY.md](SECURITY.md) for supported versions and the private vulnerability reporting path.

## License

Add the repository license before publishing a project generated from this template.
