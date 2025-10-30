FROM golang:1.25 AS builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace
# Copy the go source
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -o agent .

FROM ubuntu:25.04

WORKDIR /
COPY --from=builder /workspace/templates ./templates
COPY --from=builder /workspace/agent .
RUN apt-get update && apt-get install -y \
    bpftrace \
    linux-headers-generic \
    clang \
    llvm \
    libelf-dev \
    libbpf-dev \
    zlib1g-dev \
    cmake \
    make \
    g++ \
    libclang-dev \
    llvm-dev\
    libelf-dev\
    zlib1g-dev && rm -rf /var/lib/apt/lists/*

ENTRYPOINT ["./agent"]