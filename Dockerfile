FROM golang:1.23-alpine AS builder

WORKDIR /src

# Download dependencies as a separate step to take advantage of Docker's caching.
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=bind,source=go.mod,target=go.mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    go mod download

# Build project using bind mounts to avoid having to copy everything into the container.
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 go build -o /bin/upload-service .

FROM alpine:3.20

RUN apk --no-cache --no-progress add ca-certificates tzdata \
    && update-ca-certificates \
    && rm -rf /var/cache/apk/*

RUN adduser \
    --disabled-password \
    --home /dev/null \
    --no-create-home \
    --shell /sbin/nologin \
    --gecos upload-service \
    --uid 10000 \
    upload-service

USER upload-service

COPY --from=builder /bin/upload-service /bin

ENTRYPOINT ["/bin/upload-service"]
