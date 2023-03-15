# Step 1 - compile code binary
FROM golang:1.20.2-alpine AS builder

# make image public
LABEL org.opencontainers.image.source=https://github.com/cguertin14/ddns
LABEL org.opencontainers.image.description="ddns container image"
LABEL org.opencontainers.image.licenses=Apache
LABEL maintainer="Charles Guertin <charlesguertin@live.ca>"

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT=""

ENV CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH} \
    GOARM=${TARGETVARIANT}

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./
RUN go build -o ./ddns .


# Step 2 - import necessary files to run program.
FROM gcr.io/distroless/base-debian11:nonroot
COPY --from=builder /app/ddns /ddns
ENTRYPOINT ["/ddns"]
