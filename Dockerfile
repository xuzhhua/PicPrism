# Stage 1: 构建前端
FROM node:22-alpine AS frontend
WORKDIR /app/web
COPY web/package*.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

# Stage 2: 构建 Go 后端
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/web/dist ./web/dist
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o picprism ./cmd/picprism

# Stage 3: 最终镜像
FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/picprism ./

VOLUME ["/data"]
EXPOSE 8080
ENV PICPRISM_DATA_DIR=/data
ENV PICPRISM_PORT=8080

ENTRYPOINT ["./picprism"]
