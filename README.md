# 家計簿アプリ 仕様書

## 1. プロジェクト概要

| 項目 | 内容 |
|------|------|
| アプリ名 | 家計簿アプリ |
| 概要 | 収支（収入・支出）を登録・編集・削除し、カテゴリ別のグラフで可視化するWebアプリケーション |
| 技術スタック | バックエンド: Go + Echo / フロントエンド: Next.js (App Router) + TypeScript + Tailwind CSS |
| データベース | メモリ（デフォルト）または PostgreSQL |

---

## 2. 機能仕様

### 2.1 画面一覧

| 画面 | パス | 説明 |
|------|------|------|
| グラフ | `/` | カテゴリ別の収入・支出を棒グラフで表示。収入合計・支出合計をサマリー表示 |
| 登録 | `/register` | 新規収支の登録フォーム |
| 編集 | `/transactions` | 登録済み収支の一覧表示・編集・削除 |

### 2.2 ヘッダーメニュー

全画面共通で表示。3つの画面（グラフ・登録・編集）へ遷移するナビゲーションボタンを提供する。

### 2.3 機能詳細

#### 収支の登録

- 日付（必須、YYYY-MM-DD形式）
- 種別（必須、収入/支出）
- カテゴリ（必須、定義済みカテゴリから選択）
- 金額（必須、1以上）
- メモ（任意）

#### 収支の編集

- 一覧から対象を選択し、インラインで編集可能
- 編集内容: 日付、種別、カテゴリ、金額、メモ

#### 収支の削除

- 一覧から削除ボタンで削除
- 削除前に確認ダイアログを表示

#### グラフ表示

- カテゴリごとに収入・支出を集計
- 棒グラフで内訳を表示（収入: 緑、支出: 赤）
- 収入合計・支出合計をサマリー表示

---

## 3. システム構成

```
KakeiboApp/
├── backend/                 # Go API サーバー
│   ├── cmd/server/          # エントリーポイント
│   ├── internal/
│   │   ├── domain/          # ドメインモデル
│   │   ├── handler/         # HTTPハンドラ
│   │   └── repository/      # データアクセス層
│   └── migrations/          # DB初期化SQL
├── frontend/                # Next.js フロントエンド
│   ├── app/                 # App Router
│   │   ├── components/      # 共通コンポーネント
│   │   ├── register/        # 登録画面
│   │   └── transactions/    # 編集画面
│   └── lib/                 # API クライアント
├── docker-compose.yml       # PostgreSQL コンテナ
└── README.md                # 本仕様書
```

### 3.1 アーキテクチャ

- **バックエンド**: クリーンアーキテクチャを意識（domain / handler / repository）
- **フロントエンド**: App Router、REST API でバックエンドと通信
- **データ永続化**: インターフェースによる抽象化。メモリ実装とPostgreSQL実装を切り替え可能

---

## 4. API 仕様

### 4.1 ベースURL

- 開発環境: `http://localhost:8080`

### 4.2 エンドポイント一覧

| メソッド | パス | 説明 |
|----------|------|------|
| GET | /api/health | ヘルスチェック |
| GET | /api/categories | カテゴリ一覧取得 |
| GET | /api/transactions | 収支一覧取得 |
| POST | /api/transactions | 収支登録 |
| PUT | /api/transactions/:id | 収支更新 |
| DELETE | /api/transactions/:id | 収支削除 |

### 4.3 リクエスト・レスポンス

#### 収支登録 POST /api/transactions

**リクエスト**

```json
{
  "date": "2025-01-31",
  "type": "expense",
  "category_id": 1,
  "amount": 1500,
  "memo": "昼食"
}
```

| フィールド | 型 | 必須 | 説明 |
|------------|-----|------|------|
| date | string | ○ | YYYY-MM-DD形式 |
| type | string | ○ | "income" または "expense" |
| category_id | number | ○ | カテゴリID（1〜10） |
| amount | number | ○ | 金額（円） |
| memo | string | - | メモ |

**レスポンス（201 Created）**

登録された収支オブジェクト（id, created_at 付き）

#### 収支更新 PUT /api/transactions/:id

**リクエスト**: 登録と同様のJSON形式

#### 収支削除 DELETE /api/transactions/:id

**レスポンス（200 OK）**

```json
{
  "message": "収支が削除されました"
}
```

### 4.4 エラーレスポンス

すべてのエラーは次の形式で返却。

```json
{
  "error": "エラーメッセージ"
}
```

| HTTPステータス | 説明 |
|----------------|------|
| 400 Bad Request | バリデーションエラー（日付形式不正、type不正など） |
| 500 Internal Server Error | サーバーエラー、データ未検出時 |

---

## 5. データモデル

### 5.1 収支（Transaction）

| フィールド | 型 | 説明 |
|------------|-----|------|
| id | number | 一意ID（自動採番） |
| date | string | 取引日（ISO 8601形式） |
| type | string | "income" / "expense" |
| category_id | number | カテゴリID |
| category | object | カテゴリ詳細（id, name） |
| amount | number | 金額（支出は負の値で保持） |
| memo | string | メモ |
| created_at | string | 登録日時（ISO 8601形式） |

### 5.2 カテゴリ（Category）

| ID | 名称 |
|----|------|
| 1 | 食費 |
| 2 | 交通費 |
| 3 | 住居費 |
| 4 | 光熱費 |
| 5 | 通信費 |
| 6 | 娯楽費 |
| 7 | 医療費 |
| 8 | 教育費 |
| 9 | その他 |
| 10 | 給与 |

### 5.3 DBスキーマ（PostgreSQL）

- **categories**: id (SERIAL), name (VARCHAR)
- **transactions**: id (SERIAL), date (DATE), type (VARCHAR), category_id (FK), amount (INTEGER), memo (TEXT), created_at (TIMESTAMPTZ)

---

## 6. 環境・セットアップ

### 6.1 必要環境

- Go 1.21+
- Node.js 18+
- Docker / Docker Compose（PostgreSQL使用時）

### 6.2 起動手順

#### メモリストアで起動（簡易）

```bash
# バックエンド
cd backend
go mod tidy
go run ./cmd/server

# フロントエンド（別ターミナル）
cd frontend
npm install
npm run dev
```

- バックエンド: http://localhost:8080
- フロントエンド: http://localhost:3000

#### PostgreSQL で起動

```bash
# 1. PostgreSQL 起動
docker compose up -d

# 2. backend/.env を作成
# DATABASE_URL=postgres://kakeibo:kakeibo@localhost:5432/kakeibo?sslmode=disable

# 3. バックエンド起動
cd backend
go run ./cmd/server
```

### 6.3 環境変数

| 変数 | 用途 | 例 |
|------|------|-----|
| DATABASE_URL | DB接続（バックエンド） | postgres://kakeibo:kakeibo@localhost:5432/kakeibo?sslmode=disable |
| CORS_ORIGINS | 許可するフロントエンドオリジン（カンマ区切り） | http://192.168.1.100:3000,http://10.0.0.5:3000 |
| NEXT_PUBLIC_API_URL | APIのベースURL（ビルド時フォールバック） | http://localhost:8080 |

### 6.4 ラズパイ（ARM64）向け Docker デプロイ

#### 方法A: GitHub Actions + Self-hosted Runner（推奨・自宅ネットワーク向け）

自宅のラズパイに Runner を登録し、push 時に自動デプロイ。

**1. ラズパイの事前準備**

```bash
# Docker, Docker Compose, Git をインストール
sudo apt update && sudo apt install -y docker.io docker-compose git
sudo usermod -aG docker $USER  # 要ログアウト
```

**2. Self-hosted Runner の登録**

1. GitHub リポジトリ → Settings → Actions → Runners → New self-hosted runner
2. OS: Linux, Architecture: ARM64 を選択
3. 表示されるコマンドをラズパイで実行:

```bash
mkdir actions-runner && cd actions-runner
curl -o actions-runner-linux-arm64-*.tar.gz -L https://github.com/actions/runner/releases/...
tar xzf ./actions-runner-linux-arm64-*.tar.gz
./config.sh --url https://github.com/ユーザー名/KakeiboApp --token 発行されたトークン
./run.sh  # または sudo ./svc.sh install && sudo ./svc.sh start でサービス化
```

**3. GitHub Secrets（任意・推奨）**

Settings → Secrets and variables → Actions で追加:

| 名前 | 例 |
|------|-----|
| CORS_ORIGINS | http://192.168.1.100:3000,http://10.0.0.5:3000 |
| NEXT_PUBLIC_API_URL | http://192.168.1.100:8080 |

**4. デプロイ実行**

- main ブランチへ push で自動実行
- または Actions タブ → Deploy to Raspberry Pi → Run workflow

#### 方法B: 手動デプロイ

```bash
# プロジェクトをクローン
git clone https://github.com/ユーザー名/KakeiboApp.git && cd KakeiboApp

# .env を作成
cp .env.example .env
# CORS_ORIGINS に Wi-Fi と VPN の両方のフロントエンドURLを設定

# ビルド＆起動
docker compose -f docker-compose.prod.yml build
docker compose -f docker-compose.prod.yml up -d
```

- フロントエンド: ポート 3000
- バックエンド: ポート 8080
- PostgreSQL: コンテナ内のみ（外部公開なし）

**Wi-Fi / VPN 両対応**: フロントエンドはアクセス元ホストを自動検出してAPIに接続。バックエンドの CORS は `CORS_ORIGINS` 環境変数で複数オリジンを指定可能（カンマ区切り）。

---

## 7. 非機能要件

### 7.1 CORS

- 許可オリジン: 環境変数 `CORS_ORIGINS`（カンマ区切り）で指定。未設定時は `http://localhost:3000`
- 例（Wi-Fi+VPN）: `CORS_ORIGINS=http://192.168.1.100:3000,http://10.0.0.5:3000`
- 許可メソッド: GET, POST, PUT, DELETE, OPTIONS

### 7.2 データ永続化

- **DATABASE_URL 未設定**: メモリストア。サーバー再起動でデータ消去
- **DATABASE_URL 設定**: PostgreSQL。データ永続化

### 7.3 テスト

- バックエンド: `go test ./internal/...` でリポジトリ・ハンドラのテストを実行

### 7.4 GitHub Actions

| ワークフロー | トリガー | 内容 |
|-------------|---------|------|
| CI | push / PR | Go テスト、Next.js リント・ビルド、Docker ビルド確認 |
| Deploy | push / 手動 | ラズパイ上の Self-hosted Runner で docker compose 実行 |

**Deploy の前提**: ラズパイに Self-hosted Runner を登録すること（README 6.4 参照）。SSH 不要・自宅ネットワークから GitHub へのアウトバウンド接続のみで動作。
