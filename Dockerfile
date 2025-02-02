FROM golang:1-bookworm AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download -x

COPY . ./

# Set release var to current commit hash.
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.release=`git rev-parse --short=8 HEAD`'" -o /bin/server ./cmd/server

FROM gcr.io/distroless/base-debian11

WORKDIR /app

COPY --from=builder /bin/server .

CMD ["./server"]
