FROM registry.gitlab.com/gitlab-org/gitlab-build-images:golangci-lint-alpine as builder

WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o auth ./cmd/main.go


# generate clean, final image for end users
FROM quay.io/jitesoft/alpine:3.11.3

COPY --from=builder /build/auth /app/
COPY --from=builder /build/configs /app/configs/

# executable
ENTRYPOINT [ "/app/auth", "/app/prod" ]
