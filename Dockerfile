FROM --platform=$BUILDPLATFORM golang:1.20 AS builder

ARG TARGETARCH

RUN apt-get update && apt-get install -y gcc-aarch64-linux-gnu

WORKDIR /app

COPY . /app/

RUN if [ "$TARGETARCH" = "arm64" ]; then make dist-arm64 ; else make dist ; fi

FROM ubuntu:20.04

RUN apt-get update && apt-get install -y ca-certificates

COPY --from=builder /app/out/* /app/
COPY ./public /app/public

WORKDIR /app

ENTRYPOINT [ "./gws" ]
