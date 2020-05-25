FROM golang:1.14.3 as builder

WORKDIR /src/
COPY . /src/

ARG VERSION=v0.0.1-default

ENV VERSION=$VERSION
ENV GO111MODULE=on

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags "-X main.Version=${VERSION}" \
    -mod vendor -o ./eventmaker ./cmd

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /src/eventmaker /

ENTRYPOINT ["/eventmaker"]
