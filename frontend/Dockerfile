# 家計簿 フロントエンド
# ラズパイ (ARM64) 向けビルド: docker buildx build --platform linux/arm64 -t kakeibo-frontend .

FROM node:20-alpine AS builder

WORKDIR /app

# 依存関係を先にコピー（キャッシュ効率化）
COPY package.json package-lock.json ./
RUN npm ci

# public が無い場合に備えて作成
RUN mkdir -p public

COPY . .

# ビルド時の環境変数（APIのURLは実行時に上書き可能）
ARG NEXT_PUBLIC_API_URL=http://localhost:8080
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL

RUN npm run build

# 実行用イメージ（standalone 出力を使用）
FROM node:20-alpine AS runner

WORKDIR /app

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000

ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

CMD ["node", "server.js"]
