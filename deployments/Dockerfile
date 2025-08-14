# Build stage
FROM golang:1.25 as builder

WORKDIR /build

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/sinkhole_detox/main.go

# App container
FROM scratch

WORKDIR /

COPY --from=builder /build/main /main
USER nonroot
CMD [ "/main" ]
