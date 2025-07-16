# Stage 1: Build the application in a temporary environment
FROM golang:1.24-alpine AS builder

# Update packages to patch vulnerabilities in the build environment
RUN apk update && apk upgrade

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build a static, self-contained binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -o clima-cep-server ./cmd/server

# Stage 2: Create a minimal and secure final image
FROM gcr.io/distroless/static-debian11

WORKDIR /root/

# Copy only the compiled binary and the config file from the builder stage
COPY --from=builder /app/clima-cep-server .
COPY --from=builder /app/configs/configs.yml ./configs/configs.yml

EXPOSE 8080

CMD ["./clima-cep-server"]