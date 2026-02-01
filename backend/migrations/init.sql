-- 家計簿DB初期スキーマ
-- PostgreSQL コンテナ初回起動時に自動実行されます。

-- カテゴリテーブル
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);

-- 収支テーブル
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('income', 'expense')),
    category_id INTEGER NOT NULL REFERENCES categories(id),
    amount INTEGER NOT NULL,
    memo TEXT DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions(date);
CREATE INDEX IF NOT EXISTS idx_transactions_category_id ON transactions(category_id);

-- 初期カテゴリデータ
INSERT INTO categories (id, name) VALUES
    (1, '食費'),
    (2, '交通費'),
    (3, '住居費'),
    (4, '光熱費'),
    (5, '通信費'),
    (6, '娯楽費'),
    (7, '医療費'),
    (8, '教育費'),
    (9, 'その他'),
    (10, '給与')
ON CONFLICT (id) DO NOTHING;

SELECT setval('categories_id_seq', (SELECT MAX(id) FROM categories));
