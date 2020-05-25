FROM golang:1.13.8-alpine3.11@sha256:1ff752199f17b70e5f4dc2ad7f3e7843c456eb7e1407ed158ed8c237dbf1476a as builder

LABEL maintainer="Joakim Lazakis <joakim@kenza.ai>"

ARG APP_VERSION

# Install git for go mod
RUN set -ex \
    apk update && apk add --no-cache git

WORKDIR /kenza

# Use cached dependencies if possible
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

# Build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags="-w -s -X main.version=${APP_VERSION} -extldflags '-static'" -o worker/worker worker/*.go

# We need a Python-aware image 
FROM python:3.6-buster@sha256:aa2ba4a0cd6d9f056778dabd8f4103255d0bcb3974d9a38880d455780550a649

# Install docker (for Sagify's docker build step)
RUN apt-get update && \
    apt-get -y install apt-transport-https \
    ca-certificates \
    curl \
    gnupg2 \
    software-properties-common && \
    curl -fsSL https://download.docker.com/linux/debian/gpg | apt-key add - && \
    add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/debian $(lsb_release -cs) stable" && \
    apt-get update && \
    apt-cache policy docker-ce && \
    apt-get -y install docker-ce

# Install the AWS CLI
RUN pip3 install awscli --upgrade

COPY --from=builder /kenza/worker/worker /kenza/worker