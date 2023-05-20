FROM golang:1.20-alpine AS Builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server .

FROM alpine:3.14 AS final

COPY --from=builder /build/server /bin/server

ENTRYPOINT ["/bin/server"]