FROM golang:alpine AS build

ADD . /build
WORKDIR /build
RUN go build -o server .

FROM golang:alpine

RUN apk add --update sudo bash

RUN go install tailscale.com/cmd/tailscale@v1.22.2
RUN go install tailscale.com/cmd/tailscaled@v1.22.2

COPY --from=build /build/server /run/server
COPY --from=build /build/entrypoint.sh /run/entrypoint.sh

# Default to empty $TAILSCALE_KEY.
ARG TAILSCALE_KEY=""
RUN mkdir -p run/tailscale-storage
ENTRYPOINT bash /run/entrypoint.sh