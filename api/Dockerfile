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

WORKDIR /kenza

# Use cached dependencies if possible
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

# Build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags="-w -s -X main.version=${APP_VERSION} -extldflags '-static'" -o api/apid api/*.go

# Trim image down
FROM alpine

ARG KENZA_DB_MIGRATIONS_HOST_PATH
ENV KENZA_DB_MIGRATIONS_HOST_PATH ${KENZA_DB_MIGRATIONS_HOST_PATH}

ARG KENZA_DB_MIGRATIONS_CONTAINER_PATH
ENV KENZA_DB_MIGRATIONS_CONTAINER_PATH ${KENZA_DB_MIGRATIONS_CONTAINER_PATH}

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /kenza/api/apid /kenza/apid
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /kenza/${KENZA_DB_MIGRATIONS_HOST_PATH} ${KENZA_DB_MIGRATIONS_CONTAINER_PATH}