FROM golang:1.11-alpine AS builder
WORKDIR /go/src/github.com/domgoer/k8s-admission-webhook-example
ADD . .
RUN CGO_ENABLED=0 go build -o admission-webhook

FROM alpine:3.8
LABEL maintainer="Domgoer <doumengcheng@iftech.io>"
EXPOSE 8080
ENTRYPOINT ["/admission-webhook"]
COPY --from=builder /go/src/github.com/domgoer/k8s-admission-webhook-example/admission-webhook /
