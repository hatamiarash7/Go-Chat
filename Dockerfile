FROM --platform=$BUILDPLATFORM golang:1.19-alpine as builder

ARG APP_VERSION="undefined@docker"

WORKDIR /src

COPY *.go .
COPY go.* .

ENV LDFLAGS="-s -w -X github.com/hatamiarash7/go-chat/internal/pkg/version.version=$APP_VERSION"
ENV GO111MODULE=on
ARG TARGETOS TARGETARCH
ENV GOOS $TARGETOS
ENV GOARCH $TARGETARCH

RUN set -x \
    && go version \
    && CGO_ENABLED=0 go build -trimpath -ldflags "$LDFLAGS" -o /src/go-chat .

FROM scratch

ARG APP_VERSION="undefined@docker"

LABEL \
    org.opencontainers.image.title="go-chat" \
    org.opencontainers.image.description="Simple & Encrypted Chat" \
    org.opencontainers.image.url="https://github.com/hatamiarash7/go-chat" \
    org.opencontainers.image.source="https://github.com/hatamiarash7/go-chat" \
    org.opencontainers.image.vendor="hatamiarash7" \
    org.opencontainers.image.author="hatamiarash7" \
    org.opencontainers.version="$APP_VERSION" \
    org.opencontainers.image.created="$DATE_CREATED" \
    org.opencontainers.image.licenses="MIT"

COPY --from=builder /src /

ENV START_MODE server
ENV HOST 0.0.0.0

CMD ["/go-chat"]
