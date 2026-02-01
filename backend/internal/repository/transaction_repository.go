package repository

import (
	"fmt"
	"sync"
	"time"

	"kakeibo-app/backend/internal/domain"
)

// TransactionRepository は収支データの永続化を担当するリポジトリのインターフェースです。
// 最小限のAPIのためメモリ上に保持します（後でPostgreSQLへ拡張可能）。
type TransactionRepository interface {
	FindAll() ([]domain.Transaction, error)
	FindById(id int) (domain.Transaction, error)
	FindAllCategories() ([]domain.Category, error)
	FindCategoryById(id int) (domain.Category, error)
	Save(transaction *domain.Transaction) error
	Update(transaction *domain.Transaction) error
	Delete(id int) error
}

type transactionRepository struct {
	mu           sync.RWMutex
	transactions []domain.Transaction
	categories   []domain.Category
	nextID       int
}

// NewTransactionRepository はメモリベースのTransactionRepositoryを生成します。
func NewTransactionRepository() TransactionRepository {
	return &transactionRepository{
		transactions: []domain.Transaction{},
		categories: []domain.Category{
			{ID: 1, Name: "食費"},
			{ID: 2, Name: "交通費"},
			{ID: 3, Name: "住居費"},
			{ID: 4, Name: "光熱費"},
			{ID: 5, Name: "通信費"},
			{ID: 6, Name: "娯楽費"},
			{ID: 7, Name: "医療費"},
			{ID: 8, Name: "教育費"},
			{ID: 9, Name: "その他"},
			{ID: 10, Name: "給与"},
		},
		nextID: 1,
	}
}

func (r *transactionRepository) FindAll() ([]domain.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]domain.Transaction, len(r.transactions))
	copy(result, r.transactions)
	return result, nil
}

func (r *transactionRepository) FindAllCategories() ([]domain.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]domain.Category, len(r.categories))
	copy(result, r.categories)
	return result, nil
}

func (r *transactionRepository) FindById(id int) (domain.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, transaction := range r.transactions {
		if transaction.ID == id {
			return transaction, nil
		}
	}
	return domain.Transaction{}, fmt.Errorf("収支が見つかりません: %d", id)
}

func (r *transactionRepository) FindCategoryById(id int) (domain.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, category := range r.categories {
		if category.ID == id {
			return category, nil
		}
	}
	return domain.Category{}, fmt.Errorf("カテゴリが見つかりません: %d", id)
}

func (r *transactionRepository) Save(t *domain.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	t.ID = r.nextID
	t.CreatedAt = time.Now()
	r.nextID++
	r.transactions = append(r.transactions, *t)
	return nil
}

func (r *transactionRepository) Update(t *domain.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, transaction := range r.transactions {
		if transaction.ID == t.ID {
			t.CreatedAt = transaction.CreatedAt
			r.transactions[i] = *t
			return nil
		}
	}
	return fmt.Errorf("収支が見つかりません: %d", t.ID)
}

func (r *transactionRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, transaction := range r.transactions {
		if transaction.ID == id {
			r.transactions = append(r.transactions[:i], r.transactions[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("収支が見つかりません: %d", id)
}
