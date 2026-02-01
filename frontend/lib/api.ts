/**
 * API クライアント: バックエンド（Go + Echo）との通信を担当します。
 * 環境変数 NEXT_PUBLIC_API_URL でAPIのベースURLを指定します。
 */

export type Transaction = {
  id: number;
  date: string;
  type: "income" | "expense";
  category_id: number;
  category: Category;
  amount: number;
  memo: string;
  created_at: string;
};

export type Category = {
  id: number;
  name: string;
};

export type CreateTransactionRequest = {
  date: string;
  type: "income" | "expense";
  category_id: number;
  amount: number;
  memo: string;
};

export type UpdateTransactionRequest = {
  date: string;
  type: "income" | "expense";
  category_id: number;
  amount: number;
  memo: string;
};

const API_BASE = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export async function getTransactions(): Promise<Transaction[]> {
  const res = await fetch(`${API_BASE}/api/transactions`);
  if (!res.ok) {
    throw new Error(`収支データの取得に失敗しました: ${res.status}`);
  }
  return res.json();
}

export async function getCategories(): Promise<Category[]> {
  const res = await fetch(`${API_BASE}/api/categories`);
  if (!res.ok) {
    throw new Error(`カテゴリの取得に失敗しました: ${res.status}`);
  }
  return res.json();
}

export async function createTransaction(
  data: CreateTransactionRequest
): Promise<Transaction> {
  const res = await fetch(`${API_BASE}/api/transactions`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(
      err.error ?? `収支の登録に失敗しました: ${res.status}`
    );
  }
  return res.json();
}

export async function updateTransaction(
  id: number,
  data: UpdateTransactionRequest
): Promise<Transaction> {
  const res = await fetch(`${API_BASE}/api/transactions/${id}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(
      err.error ?? `収支の更新に失敗しました: ${res.status}`
    );
  }
  return res.json();
}

export async function deleteTransaction(
  id: number
): Promise<null> {
  const res = await fetch(`${API_BASE}/api/transactions/${id}`, {
    method: "DELETE",
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(
      err.error ?? `収支の削除に失敗しました: ${res.status}`
    );
  }
  return null;
}