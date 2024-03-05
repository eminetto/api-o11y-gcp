FROM golang:1.21-alpine AS builder


WORKDIR /workspace

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY auth/ auth/
COPY cmd/ cmd/
COPY feedback/ feedback/
COPY internal/ internal/
COPY ops/ ops/
COPY user/ user/
COPY vote/ vote/
COPY .env .env
COPY . /workspace/


# Build executable
# RUN CGO_ENABLED=1 go build -o api cmd/api/main.go
ENV CGO_ENABLED=1
RUN apk add --no-cache gcc musl-dev
RUN go build -ldflags='-s -w -extldflags "-static"' -o api ./cmd/api/main.go


FROM scratch
WORKDIR /
COPY --from=builder /workspace/api .
COPY --from=builder /workspace/.env .env
COPY --from=builder /workspace/ops/db/ ops/db/
EXPOSE 8081

ENTRYPOINT ["/api"]
