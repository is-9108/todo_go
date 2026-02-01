// main.go は家計簿APIサーバーのエントリーポイントです。
// Echoサーバーを起動し、CORSを設定してフロントエンドからのリクエストを受け付けます。
// 環境変数 DATABASE_URL が設定されている場合は PostgreSQL を使用します。
package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"kakeibo-app/backend/internal/domain"
	"kakeibo-app/backend/internal/handler"
	"kakeibo-app/backend/internal/repository"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	_ = godotenv.Load()

	e := echo.New()

	// CORS: WiFi と VPN の両方のオリジンを許可
	// CORS_ORIGINS 例: "http://192.168.1.100:3000,http://10.0.0.5:3000"
	corsOrigins := strings.Split(os.Getenv("CORS_ORIGINS"), ",")
	for i := range corsOrigins {
		corsOrigins[i] = strings.TrimSpace(corsOrigins[i])
	}
	if len(corsOrigins) == 1 && corsOrigins[0] == "" {
		corsOrigins = []string{"http://localhost:3000"} // 開発用デフォルト
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     corsOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	var repo repository.TransactionRepository
	useMemory := os.Getenv("DATABASE_URL") == ""
	if useMemory {
		repo = repository.NewTransactionRepository()
		log.Println("メモリストアを使用しています（DATABASE_URL 未設定）")
	} else {
		var err error
		repo, err = repository.NewPostgresTransactionRepository(os.Getenv("DATABASE_URL"))
		if err != nil {
			log.Fatalf("PostgreSQL接続失敗: %v", err)
		}
		log.Println("PostgreSQL に接続しました")
	}

	th := handler.NewTransactionHandler(repo)

	e.GET("/api/categories", th.GetCategories)
	e.GET("/api/transactions", th.GetTransactions)
	e.POST("/api/transactions", th.CreateTransaction)
	e.PUT("/api/transactions/:id", th.UpdateTransaction)
	e.DELETE("/api/transactions/:id", th.DeleteTransaction)
	e.GET("/api/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// メモリストア時のみサンプルデータを投入
	if useMemory {
		sampleCategory := domain.Category{ID: 1, Name: "食費"}
		sample := domain.Transaction{
			Date:       parseDate("2025-01-15"),
			Type:       "expense",
			CategoryId: sampleCategory.ID,
			Amount:     -1500,
			Memo:       "サンプル：昼食",
			Category:   sampleCategory,
		}
		_ = repo.Save(&sample)
	}

	addr := ":8080"
	log.Printf("サーバー起動: http://localhost%s", addr)
	if err := e.Start(addr); err != nil {
		log.Fatal(err)
	}
}

func parseDate(s string) time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return t
}
