FROM golang:1.21 AS builder

WORKDIR /src

COPY * .

RUN go mod tidy

# Build executable
RUN go build -o /src/api cmd/api/main.go

FROM scratch
WORKDIR /src
COPY --from=builder /src/api ./
EXPOSE 8081
CMD ["/src/api"]
