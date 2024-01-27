
FROM quay.io/projectquay/golang:1.20 as builder
WORKDIR /go/src/app
COPY . .
ARG TARGETARCH
RUN make build TARGETARCH=$TARGETARCH

FROM golang:latest
WORKDIR /
COPY --from=builder /go/src/app/ibot .
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["./publisher-slack-bot", "go"]
