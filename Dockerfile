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
ARG VERSION

# git commit sha
ARG COMMIT

# build time 
ARG DATE

# compile 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -extldflags \"-static\" -X \"main.buildVersion=${VERSION}\" -X \"main.buildRef=${COMMIT}\" -X \"main.buildTime=${DATE}\"" \
    -a \
    -tags timetzdata \
    -o /bin/telegram-webhook-exporter .


# run 
FROM scratch


COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/telegram-webhook-exporter /bin/telegram-webhook-exporter

EXPOSE 8000

ENTRYPOINT [ "telegram-webhook-exporter" ]