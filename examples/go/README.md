# Go tools - runx version

This is a set of multiple tools to work with Go codebases.

## Usage

### Linter

`lint` action will run `golangci-lint` on the codebase.

Go packages and cache are shared with `golangci-lint` to speed up the process.

### Build

`build` action will build the codebase using `docker buildx` and output as binary.

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
```

That way you will be able to run `docker runx build` without to specify anything else.

### Multi platform builds

`build:all` allows to create multi-platform builds using `docker buildx`.

The generated binaries will be compressed and stored in the `dist` directory.

To make it more convenient, define the following in `.docker/runx.yaml`:

```yaml
ref: eunomie/go
images:
  eunomie/go:
    actions:
      build:all:
        opts:
          bin_name: <bin name>
          platforms: <comma separated list of platforms>
```

### Mocks

`mocks` action will generate mocks for the codebase using `mockery`.
