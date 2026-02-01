package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"kakeibo-app/backend/internal/repository"

	"github.com/labstack/echo/v4"
)

// transaction_handler_test.go は TransactionHandler の HTTP ハンドラテストです。
// Echo のテストユーティリティを使ってリクエスト・レスポンスを検証します。

func TestGetTransactions_Empty(t *testing.T) {
	repo := repository.NewTransactionRepository()
	h := NewTransactionHandler(repo)
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/api/transactions", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.GetTransactions(c)
	if err != nil {
		t.Fatalf("GetTransactions: unexpected error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("GetTransactions: expected status 200, got %d", rec.Code)
	}

	var result []interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
		t.Fatalf("GetTransactions: invalid JSON: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("GetTransactions: expected empty array, got %d items", len(result))
	}
}

func TestCreateTransaction_Success(t *testing.T) {
	repo := repository.NewTransactionRepository()
	h := NewTransactionHandler(repo)
	e := echo.New()

	body := `{"date":"2025-01-15","type":"expense","category_id":1,"amount":1500,"memo":"昼食"}`
	req := httptest.NewRequest(http.MethodPost, "/api/transactions", bytes.NewBufferString(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.CreateTransaction(c)
	if err != nil {
		t.Fatalf("CreateTransaction: unexpected error: %v", err)
	}

	if rec.Code != http.StatusCreated {
		t.Errorf("CreateTransaction: expected status 201, got %d", rec.Code)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
		t.Fatalf("CreateTransaction: invalid JSON: %v", err)
	}
	if result["id"] == nil {
		t.Error("CreateTransaction: expected id in response")
	}
	if result["memo"] != "昼食" {
		t.Errorf("CreateTransaction: expected memo=昼食, got %v", result["memo"])
	}
}

func TestCreateTransaction_InvalidType(t *testing.T) {
	repo := repository.NewTransactionRepository()
	h := NewTransactionHandler(repo)
	e := echo.New()

	body := `{"date":"2025-01-15","type":"invalid","category_id":1,"amount":1500,"memo":"テスト"}`
	req := httptest.NewRequest(http.MethodPost, "/api/transactions", bytes.NewBufferString(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.CreateTransaction(c)
	if err != nil {
		t.Fatalf("CreateTransaction: unexpected error: %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("CreateTransaction: expected status 400 for invalid type, got %d", rec.Code)
	}
}

func TestCreateTransaction_InvalidDate(t *testing.T) {
	repo := repository.NewTransactionRepository()
	h := NewTransactionHandler(repo)
	e := echo.New()

	body := `{"date":"2025/01/15","type":"expense","category_id":1,"amount":1500,"memo":"テスト"}`
	req := httptest.NewRequest(http.MethodPost, "/api/transactions", bytes.NewBufferString(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.CreateTransaction(c)
	if err != nil {
		t.Fatalf("CreateTransaction: unexpected error: %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("CreateTransaction: expected status 400 for invalid date, got %d", rec.Code)
	}
}

func TestUpdateTransaction_Success(t *testing.T) {
	repo := repository.NewTransactionRepository()
	h := NewTransactionHandler(repo)
	e := echo.New()

	// 事前に1件作成
	createReq := httptest.NewRequest(http.MethodPost, "/api/transactions", bytes.NewBufferString(
		`{"date":"2025-01-15","type":"expense","category_id":1,"amount":1000,"memo":"元のメモ"}`))
	createReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(createReq, rec)
	_ = h.CreateTransaction(c)

	// 更新リクエスト
	updateBody := `{"date":"2025-01-20","type":"expense","category_id":2,"amount":2000,"memo":"更新後のメモ"}`
	updateReq := httptest.NewRequest(http.MethodPut, "/api/transactions/1", bytes.NewBufferString(updateBody))
	updateReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	updateRec := httptest.NewRecorder()
	updateC := e.NewContext(updateReq, updateRec)
	updateC.SetParamNames("id")
	updateC.SetParamValues("1")

	err := h.UpdateTransaction(updateC)
	if err != nil {
		t.Fatalf("UpdateTransaction: unexpected error: %v", err)
	}

	if updateRec.Code != http.StatusOK {
		t.Errorf("UpdateTransaction: expected status 200, got %d", updateRec.Code)
	}

	var updated map[string]interface{}
	if err := json.Unmarshal(updateRec.Body.Bytes(), &updated); err != nil {
		t.Fatalf("UpdateTransaction: invalid JSON: %v", err)
	}
	if updated["memo"] != "更新後のメモ" {
		t.Errorf("UpdateTransaction: expected memo=更新後のメモ, got %v", updated["memo"])
	}
}

func TestUpdateTransaction_InvalidId(t *testing.T) {
	repo := repository.NewTransactionRepository()
	h := NewTransactionHandler(repo)
	e := echo.New()

	body := `{"date":"2025-01-15","type":"expense","category_id":1,"amount":1500,"memo":"テスト"}`
	req := httptest.NewRequest(http.MethodPut, "/api/transactions/abc", bytes.NewBufferString(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc")

	err := h.UpdateTransaction(c)
	if err != nil {
		t.Fatalf("UpdateTransaction: unexpected error: %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("UpdateTransaction: expected status 400 for invalid id, got %d", rec.Code)
	}
}

func TestDeleteTransaction_Success(t *testing.T) {
	repo := repository.NewTransactionRepository()
	h := NewTransactionHandler(repo)
	e := echo.New()

	// 事前に1件作成
	createReq := httptest.NewRequest(http.MethodPost, "/api/transactions", bytes.NewBufferString(
		`{"date":"2025-01-15","type":"expense","category_id":1,"amount":1000,"memo":"削除対象"}`))
	createReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(createReq, rec)
	_ = h.CreateTransaction(c)

	// 削除
	delReq := httptest.NewRequest(http.MethodDelete, "/api/transactions/1", nil)
	delRec := httptest.NewRecorder()
	delC := e.NewContext(delReq, delRec)
	delC.SetParamNames("id")
	delC.SetParamValues("1")

	err := h.DeleteTransaction(delC)
	if err != nil {
		t.Fatalf("DeleteTransaction: unexpected error: %v", err)
	}

	if delRec.Code != http.StatusOK {
		t.Errorf("DeleteTransaction: expected status 200, got %d", delRec.Code)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(delRec.Body.Bytes(), &result); err != nil {
		t.Fatalf("DeleteTransaction: invalid JSON: %v", err)
	}
	if result["message"] == nil {
		t.Error("DeleteTransaction: expected message in response")
	}
}

func TestDeleteTransaction_InvalidId(t *testing.T) {
	repo := repository.NewTransactionRepository()
	h := NewTransactionHandler(repo)
	e := echo.New()

	req := httptest.NewRequest(http.MethodDelete, "/api/transactions/xyz", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("xyz")

	err := h.DeleteTransaction(c)
	if err != nil {
		t.Fatalf("DeleteTransaction: unexpected error: %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("DeleteTransaction: expected status 400 for invalid id, got %d", rec.Code)
	}
}
