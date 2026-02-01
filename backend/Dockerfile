# 家計簿API バックエンド
# ラズパイ (ARM64) 向けビルド: docker buildx build --platform linux/arm64 -t kakeibo-backend .

FROM golang:1.24-alpine AS builder

WORKDIR /app

# 依存関係を先にコピー（キャッシュ効率化）
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o server ./cmd/server

# 実行用イメージ
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Tokyo

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
