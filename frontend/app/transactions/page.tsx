"use client";

import { useEffect, useState } from "react";
import {
  getTransactions,
  updateTransaction,
  deleteTransaction,
  getCategories,
  type Transaction,
  type CreateTransactionRequest,
  type UpdateTransactionRequest,
  type Category,
} from "@/lib/api";

/**
 * 編集画面: 登録済み収支の一覧表示・編集・削除を提供します。
 */
export default function TransactionsPage() {
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [editForm, setEditForm] = useState<UpdateTransactionRequest | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [actionError, setActionError] = useState<string | null>(null);

  const fetchTransactions = async () => {
    try {
      setError(null);
      const data = await getTransactions();
      setTransactions(Array.isArray(data) ? data : []);
    } catch (e) {
      setError(e instanceof Error ? e.message : "データの取得に失敗しました");
    } finally {
      setLoading(false);
    }
  };

  const fetchCategories = async () => {
    try {
      const data = await getCategories();
      setCategories(Array.isArray(data) ? data : []);
    } catch {
      // カテゴリ取得失敗は編集に影響
    }
  };

  useEffect(() => {
    fetchTransactions();
    fetchCategories();
  }, []);

  const startEdit = (t: Transaction) => {
    setEditingId(t.id);
    setEditForm({
      date: new Date(t.date).toISOString().slice(0, 10),
      type: t.type,
      category_id: t.category_id,
      amount: Math.abs(t.amount),
      memo: t.memo,
    });
    setActionError(null);
  };

  const cancelEdit = () => {
    setEditingId(null);
    setEditForm(null);
    setActionError(null);
  };

  const handleUpdate = async () => {
    if (!editingId || !editForm) return;
    setSubmitting(true);
    setActionError(null);
    try {
      await updateTransaction(editingId, editForm);
      await fetchTransactions();
      cancelEdit();
    } catch (e) {
      setActionError(e instanceof Error ? e.message : "更新に失敗しました");
    } finally {
      setSubmitting(false);
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm("この収支を削除しますか？")) return;
    setSubmitting(true);
    setActionError(null);
    try {
      await deleteTransaction(id);
      await fetchTransactions();
      if (editingId === id) cancelEdit();
    } catch (e) {
      setActionError(e instanceof Error ? e.message : "削除に失敗しました");
    } finally {
      setSubmitting(false);
    }
  };

  const formatAmount = (amount: number) => {
    const abs = Math.abs(amount);
    return `${amount >= 0 ? "+" : "-"}¥${abs.toLocaleString()}`;
  };

  if (loading) {
    return <p className="text-slate-500">読み込み中...</p>;
  }

  if (error) {
    return (
      <p className="text-red-600">
        {error}
        <br />
        <span className="text-sm">
          バックエンドが http://localhost:8080 で起動しているか確認してください。
        </span>
      </p>
    );
  }

  return (
    <section className="rounded-lg bg-white p-6 shadow">
      <h2 className="mb-6 text-xl font-semibold text-slate-700">収支一覧・編集</h2>

      {actionError && (
        <p className="mb-4 text-sm text-red-600">{actionError}</p>
      )}

      {(Array.isArray(transactions) ? transactions : []).length === 0 ? (
        <p className="text-slate-500">データがありません。登録画面から追加してください。</p>
      ) : (
        <div className="overflow-x-auto">
          <table className="w-full text-left">
            <thead>
              <tr className="border-b border-slate-200">
                <th className="py-2 pr-4 text-sm font-medium text-slate-600">日付</th>
                <th className="py-2 pr-4 text-sm font-medium text-slate-600">種別</th>
                <th className="py-2 pr-4 text-sm font-medium text-slate-600">カテゴリ</th>
                <th className="py-2 pr-4 text-sm font-medium text-slate-600">金額</th>
                <th className="py-2 pr-4 text-sm font-medium text-slate-600">メモ</th>
                <th className="py-2 text-sm font-medium text-slate-600">操作</th>
              </tr>
            </thead>
            <tbody>
              {(Array.isArray(transactions) ? transactions : []).map((t) => (
                <tr key={t.id} className="border-b border-slate-100">
                  {editingId === t.id && editForm ? (
                    <>
                      <td className="py-3 pr-4">
                        <input
                          type="date"
                          value={editForm.date}
                          onChange={(e) =>
                            setEditForm((prev) =>
                              prev ? { ...prev, date: e.target.value } : null
                            )
                          }
                          className="w-full rounded border border-slate-300 px-2 py-1 text-sm"
                        />
                      </td>
                      <td className="py-3 pr-4">
                        <select
                          value={editForm.type}
                          onChange={(e) =>
                            setEditForm((prev) =>
                              prev
                                ? { ...prev, type: e.target.value as "income" | "expense" }
                                : null
                            )
                          }
                          className="rounded border border-slate-300 px-2 py-1 text-sm"
                        >
                          <option value="expense">支出</option>
                          <option value="income">収入</option>
                        </select>
                      </td>
                      <td className="py-3 pr-4">
                        <select
                          value={editForm.category_id}
                          onChange={(e) =>
                            setEditForm((prev) =>
                              prev
                                ? {
                                    ...prev,
                                    category_id: parseInt(e.target.value, 10) || 0,
                                  }
                                : null
                            )
                          }
                          className="rounded border border-slate-300 px-2 py-1 text-sm"
                        >
                          {categories.map((c) => (
                            <option key={c.id} value={c.id}>
                              {c.name}
                            </option>
                          ))}
                        </select>
                      </td>
                      <td className="py-3 pr-4">
                        <input
                          type="number"
                          min="1"
                          value={editForm.amount || ""}
                          onChange={(e) =>
                            setEditForm((prev) =>
                              prev
                                ? {
                                    ...prev,
                                    amount: parseInt(e.target.value, 10) || 0,
                                  }
                                : null
                            )
                          }
                          className="w-24 rounded border border-slate-300 px-2 py-1 text-sm"
                        />
                      </td>
                      <td className="py-3 pr-4">
                        <input
                          type="text"
                          value={editForm.memo}
                          onChange={(e) =>
                            setEditForm((prev) =>
                              prev ? { ...prev, memo: e.target.value } : null
                            )
                          }
                          className="w-full rounded border border-slate-300 px-2 py-1 text-sm"
                        />
                      </td>
                      <td className="py-3">
                        <button
                          onClick={handleUpdate}
                          disabled={submitting}
                          className="mr-2 text-sm text-blue-600 hover:underline disabled:opacity-50"
                        >
                          保存
                        </button>
                        <button
                          onClick={cancelEdit}
                          disabled={submitting}
                          className="mr-2 text-sm text-slate-600 hover:underline disabled:opacity-50"
                        >
                          キャンセル
                        </button>
                      </td>
                    </>
                  ) : (
                    <>
                      <td className="py-3 pr-4">
                        {new Date(t.date).toLocaleDateString("ja-JP")}
                      </td>
                      <td className="py-3 pr-4">
                        <span
                          className={
                            t.type === "income"
                              ? "text-emerald-600"
                              : "text-rose-600"
                          }
                        >
                          {t.type === "income" ? "収入" : "支出"}
                        </span>
                      </td>
                      <td className="py-3 pr-4">{t.category?.name ?? ""}</td>
                      <td
                        className={`py-3 pr-4 font-medium ${
                          t.amount >= 0 ? "text-emerald-600" : "text-rose-600"
                        }`}
                      >
                        {formatAmount(t.amount)}
                      </td>
                      <td className="py-3 text-slate-600">{t.memo}</td>
                      <td className="py-3">
                        <button
                          onClick={() => startEdit(t)}
                          className="mr-3 text-sm text-blue-600 hover:underline"
                        >
                          編集
                        </button>
                        <button
                          onClick={() => handleDelete(t.id)}
                          disabled={submitting}
                          className="text-sm text-red-600 hover:underline disabled:opacity-50"
                        >
                          削除
                        </button>
                      </td>
                    </>
                  )}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </section>
  );
}
