FROM golang:1.13.8-alpine3.11@sha256:1ff752199f17b70e5f4dc2ad7f3e7843c456eb7e1407ed158ed8c237dbf1476a as builder

LABEL maintainer="Joakim Lazakis <joakim@kenza.ai>"

ARG APP_VERSION

# Install git for go mod and ca-certificates and tzdata for the bare image
RUN set -ex \
    apk update && apk add --no-cache \ 
    ca-certificates \
    git \ 
    tzdata \
    && update-ca-certificates

# Add non-privileged user
RUN adduser -D -g '' progress

WORKDIR /kenza

# Use caches dependencies if possible
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

# Build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags="-w -s -X main.version=${APP_VERSION} -extldflags '-static'" -o progress/progress progress/*.go

# Trim image down
FROM alpine

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /kenza/progress/progress /kenza/progress
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Use non-privileged user added earlier
USER progress
