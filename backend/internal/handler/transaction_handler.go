package handler

import (
	"net/http"
	"strconv"
	"time"

	"kakeibo-app/backend/internal/domain"
	"kakeibo-app/backend/internal/repository"

	"github.com/labstack/echo/v4"
)

// TransactionHandler は収支関連のHTTPリクエストを処理するハンドラです。
type TransactionHandler struct {
	repo repository.TransactionRepository
}

// NewTransactionHandler はTransactionHandlerを生成します。
func NewTransactionHandler(repo repository.TransactionRepository) *TransactionHandler {
	return &TransactionHandler{repo: repo}
}

// GetCategories は全カテゴリを取得するGET /api/categoriesのハンドラです。
func (h *TransactionHandler) GetCategories(c echo.Context) error {
	categories, err := h.repo.FindAllCategories()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "カテゴリの取得に失敗しました: " + err.Error(),
		})
	}
	return c.JSON(http.StatusOK, categories)
}

// GetTransactions は全収支データを取得するGET /api/transactionsのハンドラです。
func (h *TransactionHandler) GetTransactions(c echo.Context) error {
	transactions, err := h.repo.FindAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "収支データの取得に失敗しました: " + err.Error(),
		})
	}
	return c.JSON(http.StatusOK, transactions)
}

// CreateTransaction は新規収支を登録するPOST /api/transactionsのハンドラです。
func (h *TransactionHandler) CreateTransaction(c echo.Context) error {
	var req domain.CreateTransactionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "リクエストボディの解析に失敗しました: " + err.Error(),
		})
	}

	if req.Type != "income" && req.Type != "expense" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "typeは income または expense を指定してください",
		})
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "dateは YYYY-MM-DD 形式で指定してください",
		})
	}

	amount := req.Amount
	if req.Type == "expense" && amount > 0 {
		amount = -amount // 支出は負の値で統一
	} else if req.Type == "income" && amount < 0 {
		amount = -amount // 収入は正の値で統一
	}

	categoryId := req.CategoryId
	category, err := h.repo.FindCategoryById(categoryId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "カテゴリの取得に失敗しました: " + err.Error(),
		})
	}

	transaction := domain.Transaction{
		Date:       date,
		Type:       req.Type,
		CategoryId: categoryId,
		Amount:     amount,
		Memo:       req.Memo,
		Category:   category,
	}

	if err := h.repo.Save(&transaction); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "収支の保存に失敗しました: " + err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, transaction)
}

// UpdateTransaction は収支を更新するPUT /api/transactions/{id}のハンドラです。
func (h *TransactionHandler) UpdateTransaction(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "idは整数で指定してください",
		})
	}

	var req domain.UpdateTransactionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "リクエストボディの解析に失敗しました: " + err.Error(),
		})
	}

	if req.Type != "income" && req.Type != "expense" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "typeは income または expense を指定してください",
		})
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "dateは YYYY-MM-DD 形式で指定してください",
		})
	}

	amount := req.Amount
	if req.Type == "expense" && amount > 0 {
		amount = -amount // 支出は負の値で統一
	} else if req.Type == "income" && amount < 0 {
		amount = -amount // 収入は正の値で統一
	}

	categoryId := req.CategoryId
	category, err := h.repo.FindCategoryById(categoryId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "カテゴリの取得に失敗しました: " + err.Error(),
		})
	}

	transaction := domain.Transaction{
		ID:         id,
		Date:       date,
		Type:       req.Type,
		CategoryId: categoryId,
		Amount:     amount,
		Memo:       req.Memo,
		Category:   category,
	}

	if err := h.repo.Update(&transaction); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "収支の更新に失敗しました: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, transaction)
}

// DeleteTransaction は収支を削除するDELETE /api/transactions/{id}のハンドラです。
func (h *TransactionHandler) DeleteTransaction(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "idは整数で指定してください",
		})
	}
	if err := h.repo.Delete(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "収支の削除に失敗しました: " + err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "収支が削除されました",
	})
}
