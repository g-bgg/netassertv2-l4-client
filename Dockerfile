FROM golang:1.22-alpine AS builder
COPY . /build
WORKDIR /build

RUN go mod download && \
    CGO_ENABLED=0 GO111MODULE=on GOOS=linux go build -v -o /netassertv2-l4-client . && \
    ls -ltr /netassertv2-l4-client

FROM gcr.io/distroless/base:nonroot
COPY --from=builder /netassertv2-l4-client /usr/bin/netassertv2-l4-client

ENTRYPOINT [ "/usr/bin/netassertv2-l4-client" ]
