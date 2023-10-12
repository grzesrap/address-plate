## Build
FROM golang:1.20-bullseye AS build

RUN apt-get update && apt-get install -y ca-certificates openssl dumb-init

ARG cert_location=/usr/local/share/ca-certificates

# Get certificate
RUN openssl s_client -showcerts -connect github.com:443 </dev/null 2>/dev/null | openssl x509 -outform PEM > ${cert_location}/github.crt
RUN openssl s_client -showcerts -connect proxy.golang.org:443 </dev/null 2>/dev/null | openssl x509 -outform PEM > ${cert_location}/proxy.golang.crt
RUN update-ca-certificates

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download -x

COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /address-plate

## Deploy
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /usr/bin/dumb-init /usr/bin/dumb-init
COPY --from=build /address-plate /address-plate

USER nonroot:nonroot

EXPOSE 8080
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/address-plate"]
