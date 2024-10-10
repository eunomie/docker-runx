# `docker runx` examples

This directory contains examples used to decorate some images with `runx`.

## Build

### [hello](hello)

```
$ cd hello
$ docker runx decorate alpine -t NAMESPACE/hello
```

An already built image is available under `eunomie/runx-hello`.

### [golangci](golangci)

```
$ cd golangci
$ docker runx decorate golangci/golangci-lint:latest -t NAMESPACE/golangci
```

An already built image is available under `eunomie/runx-golangci`.

### [go](go)

```
$ cd go
$ docker runx decorate scratch -t NAMESPACE/go
```

An already built image is available under `eunomie/go`.

## Usage

Once you built your images, explore them using `--docs` and `--list` options.
