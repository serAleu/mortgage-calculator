# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
# Копируем сначала только файлы модулей для лучшего кэширования слоев Docker
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/server

# Final stage - используем "scratch" для минимального размера
FROM scratch
COPY --from=builder /server /server
COPY config.yml ./
EXPOSE 8282
CMD ["/server"]