# Go tools - runx version

This is a set of multiple tools to work with Go codebases.

## Usage

### Linter

`lint` action will run `golangci-lint` on the codebase.

Go packages and cache are shared with `golangci-lint` to speed up the process.

### Build

`build` action will build the codebase using `docker buildx` and output as (multi-platform) binaries.

By convention if the binary name is `tool`, it will try to build `cmd/tool/` and the output will be inside `dist` directory.

To make it more convenient while working on a code base, create a file `.docker/runx.yaml` with the following content:

```yaml
ref: eunomie/go
images:
  eunomie/go:
    actions:
      build:
        opts:
          bin_name: <bin name>
          platforms: <platforms to build against>
```

That way you will be able to run `docker runx build` without to specify anything else.

### Mocks

`mocks` action will generate mocks for the codebase using `mockery`.
