package domain

import "time"

// Transaction は収支データを表すドメインモデルです。
// 家計簿の1件の収入または支出を保持します。
type Transaction struct {
	ID         int       `json:"id"`
	Date       time.Time `json:"date"`
	Type       string    `json:"type"` // "income" または "expense"
	CategoryId int       `json:"category_id"`
	Amount     int       `json:"amount"`
	Memo       string    `json:"memo"`
	CreatedAt  time.Time `json:"created_at"`
	Category   Category  `json:"category"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// CreateTransactionRequest は新規収支登録時のリクエストボディです。
type CreateTransactionRequest struct {
	Date       string `json:"date"` // "2006-01-02" 形式
	Type       string `json:"type"` // "income" または "expense"
	CategoryId int    `json:"category_id"`
	Amount     int    `json:"amount"`
	Memo       string `json:"memo"`
}

// UpdateTransactionRequest は収支更新時のリクエストボディです。
type UpdateTransactionRequest struct {
	Date       string `json:"date"` // "2006-01-02" 形式
	Type       string `json:"type"` // "income" または "expense"
	CategoryId int    `json:"category_id"`
	Amount     int    `json:"amount"`
	Memo       string `json:"memo"`
}