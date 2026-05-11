# go-simplestreams

`go-simplestreams` is a Go library for working with the simplestreams protocol.
It provides protocol types and helpers for building and validating simplestreams metadata.

## Install

```sh
go get github.com/meigma/go-simplestreams
```

The package name is `simplestreams`.

## Development

Prerequisites:

- Go 1.26.2
- Moon 2.x
- Node.js 22.22.2 for the Docusaurus docs project

```sh
moon run root:check
```

Useful individual checks:

```sh
go build ./...
go test ./...
moon run root:lint
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## Security

See [SECURITY.md](SECURITY.md).

## License

Licensed under either of:

- [Apache License, Version 2.0](LICENSE-APACHE)
- [MIT License](LICENSE-MIT)
