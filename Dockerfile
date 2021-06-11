FROM golang:latest AS build
ARG GO_OS="linux"
ARG GO_ARCH="amd64"
WORKDIR /build/
COPY . .

# Build binary output
RUN go mod download && \
    GOOS=${GO_OS} GOARCH=${GO_ARCH} CGO_ENABLED=0 go build -v -o satotx main.go

FROM alpine:latest
COPY --from=build /build/satotx satotx
ENV LISTEN 0.0.0.0:8000
EXPOSE 8000
CMD ["./satotx"]
