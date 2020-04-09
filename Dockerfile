# Build the manager binary
FROM golang:1.13 as builder


WORKDIR /workspace

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy in the go src
COPY pkg/    pkg/
COPY cmd/    cmd/
COPY notification/ notification/
COPY util/ util/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o manager ./cmd/manager

# Copy the controller-manager into a thin image
FROM scratch
WORKDIR /
COPY --from=builder /workspace/manager .
ENTRYPOINT ["/manager"]
