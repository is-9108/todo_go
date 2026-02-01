package repository

import (
	"testing"
	"time"

	"kakeibo-app/backend/internal/domain"
)

// transaction_repository_test.go は TransactionRepository の単体テストです。
// メモリベースのリポジトリの CRUD 操作を検証します。

func TestTransactionRepository_FindAll(t *testing.T) {
	repo := NewTransactionRepository()

	// 初期状態は空
	all, err := repo.FindAll()
	if err != nil {
		t.Fatalf("FindAll: unexpected error: %v", err)
	}
	if len(all) != 0 {
		t.Errorf("FindAll: expected 0 transactions, got %d", len(all))
	}

	// 1件保存後
	tx := &domain.Transaction{
		Date:       time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
		Type:       "expense",
		CategoryId: 1,
		Amount:     -1000,
		Memo:       "テスト",
	}
	if err := repo.Save(tx); err != nil {
		t.Fatalf("Save: unexpected error: %v", err)
	}

	all, err = repo.FindAll()
	if err != nil {
		t.Fatalf("FindAll: unexpected error: %v", err)
	}
	if len(all) != 1 {
		t.Errorf("FindAll: expected 1 transaction, got %d", len(all))
	}
	if all[0].ID != 1 || all[0].Memo != "テスト" {
		t.Errorf("FindAll: unexpected data: %+v", all[0])
	}
}

func TestTransactionRepository_FindById(t *testing.T) {
	repo := NewTransactionRepository()

	// 存在しないID
	_, err := repo.FindById(999)
	if err == nil {
		t.Error("FindById: expected error for non-existent ID")
	}

	// 保存して取得
	tx := &domain.Transaction{
		Date:       time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
		Type:       "income",
		CategoryId: 10,
		Amount:     50000,
		Memo:       "給与",
	}
	if err := repo.Save(tx); err != nil {
		t.Fatalf("Save: unexpected error: %v", err)
	}

	found, err := repo.FindById(1)
	if err != nil {
		t.Fatalf("FindById: unexpected error: %v", err)
	}
	if found.ID != 1 || found.Memo != "給与" || found.Amount != 50000 {
		t.Errorf("FindById: unexpected data: %+v", found)
	}
}

func TestTransactionRepository_FindCategoryById(t *testing.T) {
	repo := NewTransactionRepository()

	// 存在するカテゴリ
	cat, err := repo.FindCategoryById(1)
	if err != nil {
		t.Fatalf("FindCategoryById: unexpected error: %v", err)
	}
	if cat.ID != 1 || cat.Name != "食費" {
		t.Errorf("FindCategoryById: expected 食費, got %+v", cat)
	}

	// 存在しないカテゴリ
	_, err = repo.FindCategoryById(999)
	if err == nil {
		t.Error("FindCategoryById: expected error for non-existent category")
	}
}

func TestTransactionRepository_Save(t *testing.T) {
	repo := NewTransactionRepository()

	tx := &domain.Transaction{
		Date:       time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC),
		Type:       "expense",
		CategoryId: 2,
		Amount:     500,
		Memo:       "電車",
	}

	if err := repo.Save(tx); err != nil {
		t.Fatalf("Save: unexpected error: %v", err)
	}

	if tx.ID != 1 {
		t.Errorf("Save: expected ID=1, got %d", tx.ID)
	}
	if tx.CreatedAt.IsZero() {
		t.Error("Save: CreatedAt should be set")
	}

	all, _ := repo.FindAll()
	if len(all) != 1 {
		t.Errorf("Save: expected 1 transaction after save, got %d", len(all))
	}
}

func TestTransactionRepository_Update(t *testing.T) {
	repo := NewTransactionRepository()

	tx := &domain.Transaction{
		Date:       time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
		Type:       "expense",
		CategoryId: 1,
		Amount:     -1000,
		Memo:       "元のメモ",
	}
	if err := repo.Save(tx); err != nil {
		t.Fatalf("Save: unexpected error: %v", err)
	}
	originalCreatedAt := tx.CreatedAt

	// 更新
	tx.Memo = "更新後のメモ"
	tx.Amount = -1500
	if err := repo.Update(tx); err != nil {
		t.Fatalf("Update: unexpected error: %v", err)
	}

	updated, err := repo.FindById(1)
	if err != nil {
		t.Fatalf("FindById: unexpected error: %v", err)
	}
	if updated.Memo != "更新後のメモ" {
		t.Errorf("Update: expected Memo=更新後のメモ, got %s", updated.Memo)
	}
	if updated.Amount != -1500 {
		t.Errorf("Update: expected Amount=-1500, got %d", updated.Amount)
	}
	if !updated.CreatedAt.Equal(originalCreatedAt) {
		t.Error("Update: CreatedAt should be preserved")
	}

	// 存在しないIDで更新
	notFound := &domain.Transaction{ID: 999, Memo: "存在しない"}
	if err := repo.Update(notFound); err == nil {
		t.Error("Update: expected error for non-existent ID")
	}
}

func TestTransactionRepository_Delete(t *testing.T) {
	repo := NewTransactionRepository()

	tx := &domain.Transaction{
		Date:       time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
		Type:       "expense",
		CategoryId: 1,
		Amount:     -1000,
		Memo:       "削除対象",
	}
	if err := repo.Save(tx); err != nil {
		t.Fatalf("Save: unexpected error: %v", err)
	}

	if err := repo.Delete(1); err != nil {
		t.Fatalf("Delete: unexpected error: %v", err)
	}

	_, err := repo.FindById(1)
	if err == nil {
		t.Error("Delete: transaction should not exist after delete")
	}

	all, _ := repo.FindAll()
	if len(all) != 0 {
		t.Errorf("Delete: expected 0 transactions, got %d", len(all))
	}

	// 存在しないIDで削除
	if err := repo.Delete(999); err == nil {
		t.Error("Delete: expected error for non-existent ID")
	}
}
