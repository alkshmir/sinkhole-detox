# Build stage
FROM golang:1.25 as builder

WORKDIR /build

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o main ./cmd/sinkhole_detox/main.go
RUN echo "nobody:x:65534:65534:Nobody:/:" > /etc/passwd

# App container
FROM scratch

WORKDIR /

COPY --from=builder /build/main /main
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /etc/passwd /etc/passwd
USER nobody
CMD [ "/main" ]
