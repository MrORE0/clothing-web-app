# ---- Build stage ----
# Stage 1: Build the Go app
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

# Install build tools
RUN apk add --no-cache git

# Copy go.mod and go.sum first (cache dependencies)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the project
COPY . .

# Build the binary
RUN go build -o scraper .

# Stage 2: Runtime container
FROM alpine:3.20

WORKDIR /app

# Install Chromium + necessary libs for chromedp
RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ca-certificates \
    ttf-freefont \
    dumb-init

# Set the path for chromedp to find Chromium
ENV CHROME_PATH=/usr/bin/chromium-browser
ENV PATH=$PATH:/usr/bin

# Copy binary
COPY --from=builder /app/scraper .

ENTRYPOINT ["dumb-init", "./scraper"]
