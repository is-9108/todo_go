package repository

import (
	"context"
	"database/sql"
	"fmt"

	"kakeibo-app/backend/internal/domain"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// postgresTransactionRepository は PostgreSQL 用の TransactionRepository 実装です。
type postgresTransactionRepository struct {
	db *sql.DB
}

// NewPostgresTransactionRepository は PostgreSQL に接続した TransactionRepository を返します。
// connString 例: "postgres://kakeibo:kakeibo@localhost:5432/kakeibo?sslmode=disable"
func NewPostgresTransactionRepository(connString string) (TransactionRepository, error) {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("DB接続失敗: %w", err)
	}

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("DB疎通確認失敗: %w", err)
	}

	return &postgresTransactionRepository{db: db}, nil
}

func (r *postgresTransactionRepository) FindAll() ([]domain.Transaction, error) {
	rows, err := r.db.QueryContext(context.Background(), `
		SELECT t.id, t.date, t.type, t.category_id, t.amount, t.memo, t.created_at,  c.id, c.name
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		ORDER BY t.date DESC, t.id DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("FindAll: %w", err)
	}
	defer rows.Close()

	var result []domain.Transaction
	for rows.Next() {
		var t domain.Transaction
		var catID sql.NullInt64
		var catName sql.NullString
		if err := rows.Scan(
			&t.ID, &t.Date, &t.Type, &t.CategoryId, &t.Amount, &t.Memo, &t.CreatedAt,
			&catID, &catName,
		); err != nil {
			return nil, fmt.Errorf("FindAll scan: %w", err)
		}
		if catID.Valid && catName.Valid {
			t.Category = domain.Category{ID: int(catID.Int64), Name: catName.String}
		}
		result = append(result, t)
	}
	return result, rows.Err()
}

func (r *postgresTransactionRepository) FindById(id int) (domain.Transaction, error) {
	var t domain.Transaction
	var catID sql.NullInt64
	var catName sql.NullString
	err := r.db.QueryRowContext(context.Background(), `
		SELECT t.id, t.date, t.type, t.category_id, t.amount, t.memo, t.created_at,  c.id, c.name
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.id = $1
	`, id).Scan(
		&t.ID, &t.Date, &t.Type, &t.CategoryId, &t.Amount, &t.Memo, &t.CreatedAt,
		&catID, &catName,
	)
	if err == sql.ErrNoRows {
		return domain.Transaction{}, fmt.Errorf("収支が見つかりません: %d", id)
	}
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("FindById: %w", err)
	}
	if catID.Valid && catName.Valid {
		t.Category = domain.Category{ID: int(catID.Int64), Name: catName.String}
	}
	return t, nil
}

func (r *postgresTransactionRepository) FindAllCategories() ([]domain.Category, error) {
	rows, err := r.db.QueryContext(context.Background(), `SELECT id, name FROM categories ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("FindAllCategories: %w", err)
	}
	defer rows.Close()

	var result []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, fmt.Errorf("FindAllCategories scan: %w", err)
		}
		result = append(result, c)
	}
	return result, rows.Err()
}

func (r *postgresTransactionRepository) FindCategoryById(id int) (domain.Category, error) {
	var c domain.Category
	err := r.db.QueryRowContext(context.Background(), `SELECT id, name FROM categories WHERE id = $1`, id).
		Scan(&c.ID, &c.Name)
	if err == sql.ErrNoRows {
		return domain.Category{}, fmt.Errorf("カテゴリが見つかりません: %d", id)
	}
	if err != nil {
		return domain.Category{}, fmt.Errorf("FindCategoryById: %w", err)
	}
	return c, nil
}

func (r *postgresTransactionRepository) Save(t *domain.Transaction) error {
	err := r.db.QueryRowContext(context.Background(), `
		INSERT INTO transactions (date, type, category_id, amount, memo)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`, t.Date, t.Type, t.CategoryId, t.Amount, t.Memo).Scan(&t.ID, &t.CreatedAt)
	if err != nil {
		return fmt.Errorf("Save: %w", err)
	}
	return nil
}

func (r *postgresTransactionRepository) Update(t *domain.Transaction) error {
	result, err := r.db.ExecContext(context.Background(), `
		UPDATE transactions
		SET date = $1, type = $2, category_id = $3, amount = $4, memo = $5
		WHERE id = $6
	`, t.Date, t.Type, t.CategoryId, t.Amount, t.Memo, t.ID)
	if err != nil {
		return fmt.Errorf("Update: %w", err)
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return fmt.Errorf("収支が見つかりません: %d", t.ID)
	}
	return nil
}

func (r *postgresTransactionRepository) Delete(id int) error {
	result, err := r.db.ExecContext(context.Background(), `DELETE FROM transactions WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("Delete: %w", err)
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return fmt.Errorf("収支が見つかりません: %d", id)
	}
	return nil
}
