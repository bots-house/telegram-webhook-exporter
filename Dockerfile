# build static binary
FROM golang:1.16.0-alpine3.13 as builder 

# hadolint ignore=DL3018
RUN apk --no-cache add  \
    ca-certificates \
    git 

WORKDIR /go/src/github.com/bots-house/telegram-webhook-exporter

# download dependencies 
COPY go.mod go.sum ./
RUN go mod download 

COPY . .

# git tag 
ARG BUILD_VERSION

# git commit sha
ARG BUILD_REF

# build time 
ARG BUILD_TIME

# compile 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
      -ldflags="-w -s -extldflags \"-static\" -X \"main.buildVersion=${BUILD_VERSION}\" -X \"main.buildRef=${BUILD_REF}\" -X \"main.buildTime=${BUILD_TIME}\"" \
      -a \
      -tags timetzdata \
      -o /bin/telegram-webhook-exporter .


# run 
FROM scratch


COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/telegram-webhook-exporter /bin/telegram-webhook-exporter

EXPOSE 8000

# Reference: https://github.com/opencontainers/image-spec/blob/master/annotations.md
LABEL org.opencontainers.image.source="https://github.com/bots-house/telegram-webhook-exporter"

ENTRYPOINT [ "telegram-webhook-exporter" ]