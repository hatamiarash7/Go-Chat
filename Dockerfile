##################################### Build #####################################

FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS builder

ARG APP_VERSION="dev"
ARG TARGETOS
ARG TARGETARCH

WORKDIR /src

# Cache dependencies first for faster rebuilds.
COPY go.mod go.sum ./
RUN go mod download

# Copy source code.
COPY cmd/ cmd/
COPY internal/ internal/

ENV CGO_ENABLED=0
ENV GOOS=${TARGETOS}
ENV GOARCH=${TARGETARCH}

RUN go build -trimpath \
    -ldflags="-s -w -X github.com/hatamiarash7/go-chat/internal/version.version=${APP_VERSION}" \
    -o /src/go-chat \
    ./cmd/go-chat

##################################### Compression ########################################

FROM hatamiarash7/upx:1.1.0 AS compressor

COPY --from=builder /src/go-chat /workspace/app

RUN upx --best --lzma -o /workspace/app-compressed /workspace/app

##################################### Final ########################################

FROM scratch

ARG APP_VERSION="dev"
ARG BUILD_DATE

LABEL \
    org.opencontainers.image.title="go-chat" \
    org.opencontainers.image.description="Simple & Encrypted Chat Server" \
    org.opencontainers.image.url="https://github.com/hatamiarash7/go-chat" \
    org.opencontainers.image.source="https://github.com/hatamiarash7/go-chat" \
    org.opencontainers.image.vendor="hatamiarash7" \
    org.opencontainers.image.version="$APP_VERSION" \
    org.opencontainers.image.created="$BUILD_DATE" \
    org.opencontainers.image.licenses="MIT"

COPY --from=compressor /workspace/app-compressed /go-chat

ENV START_MODE=server
ENV HOST=0.0.0.0
ENV PORT=12345

EXPOSE 12345

ENTRYPOINT ["/go-chat"]
