FROM golang:1.19-alpine AS builder
RUN apk add upx

WORKDIR /go/src/github.com/infinimesh/tn_fake_feeder

COPY go.mod go.sum main.go ./
ADD pkg pkg

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o faker .
RUN upx ./faker

RUN apk add -U --no-cache ca-certificates

FROM scratch

COPY --from=builder /go/src/github.com/infinimesh/tn_fake_feeder/faker /faker
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ADD track.db /

LABEL org.opencontainers.image.source https://github.com/infinimesh/tn_fake_feeder

ENTRYPOINT [ "/faker" ]