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
RUN apk add --no-cache git ca-certificates openssh-client zip

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
ARG TARGETOS
ARG TARGETARCH
ARG BIN_NAME
ARG NO_ARCHIVE
ENV CGO_ENABLED=0
RUN --mount=type=cache,target=/root/.cache \
    --mount=type=cache,target=/go/pkg/mod \
    GIT_VERSION=$(git describe --tags | cut -c 2-) && \
    PKG_NAME=$(go mod graph | head -n 1 | cut -d ' ' -f 1) && \
    xx-go build \
      -o dist/${BIN_NAME} \
      -ldflags="-w -s \
        -X $PKG_NAME/internal/constants.Version=$GIT_VERSION" \
      ./cmd/${BIN_NAME} && \
    xx-verify dist/${BIN_NAME} && \
    if [ -z "${NO_ARCHIVE}" ]; then \
      # on windows add the .exe extension and zip the binary \
      if [ "${TARGETOS}" = "windows" ]; then \
        mv dist/${BIN_NAME} dist/${BIN_NAME}.exe && \
        (cd dist && zip ${BIN_NAME}-${TARGETOS}-${TARGETARCH}.zip ${BIN_NAME}.exe && rm -f ${BIN_NAME}.exe); \
      fi && \
      # if target os is not windows, tar and gzip the binary \
      if [ "${TARGETOS}" != "windows" ]; then \
        tar -C dist -czf dist/${BIN_NAME}-${TARGETOS}-${TARGETARCH}.tar.gz ${BIN_NAME} && rm -f dist/${BIN_NAME}; \
      fi \
    fi

FROM scratch AS export-bin
ARG BIN_NAME
ARG TARGETOS
ARG TARGETARCH
COPY --from=binary /go/src/dist/* /
