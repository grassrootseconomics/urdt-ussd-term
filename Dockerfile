FROM golang:1.23.3-bookworm AS build

ENV CGO_ENABLED=1

ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG BUILD=dev

WORKDIR /build
COPY . .
RUN apt update && apt install libgdbm-dev

WORKDIR /build
RUN echo "Building on $BUILDPLATFORM, building for $TARGETPLATFORM"
RUN go mod download
RUN go build -tags logtrace -o ussd-term -ldflags="-X main.build=${BUILD} -s -w" cmd/main.go

FROM debian:bookworm-slim

ENV DEBIAN_FRONTEND=noninteractive

RUN apt update && apt install libgdbm-dev ca-certificates -y
RUN apt-get clean && rm -rf /var/lib/apt/lists/*

WORKDIR /service

COPY --from=build /build/ussd-term .
COPY --from=build /build/LICENSE .
COPY --from=build /build/.env.example .
RUN mv .env.example .env

CMD ["./ussd-term"]