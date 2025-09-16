# Stage 1: Build frontend with pnpm
FROM node:22-alpine AS frontend-builder

WORKDIR /app/frontend

# Install pnpm
RUN corepack enable

COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN pnpm dependencies

COPY frontend/ .

# env vars for fallback
ENV VITE_MONTO_BS=150
ENV VITE_MONTO_USD=1

RUN pnpm run build

# Stage 2: Build backend with Go 1.24.x
FROM golang:1.24-alpine AS backend-builder

WORKDIR /app

# Install git (needed for some Go modules) and ca-certificates for HTTPS
RUN apk add --no-cache git ca-certificates

# Copy Go modules manifests and download dependencies
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend source code
COPY backend/ ./

# Copy frontend build output from frontend-builder into backend embed folder
COPY --from=frontend-builder /app/frontend/dist ./cmd/app/dist

# Build the backend binary
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -o bin/app ./cmd/app/main.go

# Stage 3: Minimal runtime image
FROM alpine:latest

WORKDIR /app

RUN addgroup -S app && adduser -S app -G app

# Install ca-certificates to avoid HTTPS errors at runtime
RUN apk add --no-cache ca-certificates

# Copy the backend binary
COPY --from=backend-builder /app/bin/app ./app

# Copy migrations folder
COPY --from=backend-builder /app/migrations ./migrations

EXPOSE 8080
USER app

CMD ["./app"]
