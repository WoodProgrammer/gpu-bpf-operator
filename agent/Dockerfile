FROM golang:1.25 AS builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace
# Copy the go source
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o agent main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/templates ./templates
COPY --from=builder /workspace/agent .
USER 65532:65532
ENTRYPOINT ["./agent"]