# syntax=docker/dockerfile:1.4

ARG XX_VERSION=1.2.1
ARG ALPINE_VERSION=3.20
ARG GO_VERSION=1.23.1

ARG BIN_NAME

FROM --platform=$BUILDPLATFORM tonistiigi/xx:${XX_VERSION} AS xx

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS build-base
COPY --from=xx / /
RUN apk add --no-cache curl
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin
RUN apk add --no-cache git ca-certificates openssh-client

FROM build-base AS build
ARG TARGETPLATFORM
RUN xx-go --wrap
WORKDIR /go/src/
COPY go.mod ./
COPY go.sum ./
RUN --mount=type=ssh \
    --mount=type=cache,target=/root/.cache \
    --mount=type=cache,target=/go/pkg/mod \
     go mod download
COPY . ./

FROM build AS binary
ENV CGO_ENABLED=0
ARG BIN_NAME
RUN --mount=type=cache,target=/root/.cache \
    --mount=type=cache,target=/go/pkg/mod \
    GIT_VERSION=$(git describe --tags | cut -c 2-) && \
    xx-go build \
      -o dist/${BIN_NAME} \
      -ldflags="-w -s \
        -X {{.PKG_NAME}}/internal/constants.Version=$GIT_VERSION" \
      ./cmd/${BIN_NAME} && \
    xx-verify dist/${BIN_NAME}

FROM scratch AS export-bin
ARG BIN_NAME
ARG TARGETOS
ARG TARGETARCH
COPY --from=binary /go/src/dist/${BIN_NAME} /${BIN_NAME}-${TARGETOS}-${TARGETARCH}
