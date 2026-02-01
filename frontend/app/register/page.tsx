"use client";

import { useEffect, useState } from "react";
import {
  createTransaction,
  getCategories,
  type Category,
  type CreateTransactionRequest,
} from "@/lib/api";
import { useRouter } from "next/navigation";

/**
 * 登録画面: 収支の新規登録フォームを提供します。
 */
export default function RegisterPage() {
  const router = useRouter();
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [form, setForm] = useState<CreateTransactionRequest>({
    date: new Date().toISOString().slice(0, 10),
    type: "expense",
    category_id: 0,
    amount: 0,
    memo: "",
  });
  const [submitting, setSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState<string | null>(null);

  useEffect(() => {
    const fetchCategories = async () => {
      try {
        setError(null);
        const data = await getCategories();
        setCategories(Array.isArray(data) ? data : []);
      } catch (e) {
        setError(e instanceof Error ? e.message : "カテゴリの取得に失敗しました");
      } finally {
        setLoading(false);
      }
    };
    fetchCategories();
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitError(null);
    setSubmitting(true);
    try {
      await createTransaction(form);
      setForm({
        date: new Date().toISOString().slice(0, 10),
        type: "expense",
        category_id: 0,
        amount: 0,
        memo: "",
      });
      router.push("/");
    } catch (e) {
      setSubmitError(e instanceof Error ? e.message : "登録に失敗しました");
    } finally {
      setSubmitting(false);
    }
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
      <h2 className="mb-6 text-xl font-semibold text-slate-700">新規登録</h2>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <div>
            <label className="mb-1 block text-sm text-slate-600">日付</label>
            <input
              type="date"
              required
              value={form.date}
              onChange={(e) =>
                setForm((prev) => ({ ...prev, date: e.target.value }))
              }
              className="w-full rounded border border-slate-300 px-3 py-2 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
            />
          </div>
          <div>
            <label className="mb-1 block text-sm text-slate-600">種別</label>
            <select
              value={form.type}
              onChange={(e) =>
                setForm((prev) => ({
                  ...prev,
                  type: e.target.value as "income" | "expense",
                }))
              }
              className="w-full rounded border border-slate-300 px-3 py-2 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
            >
              <option value="expense">支出</option>
              <option value="income">収入</option>
            </select>
          </div>
          <div>
            <label className="mb-1 block text-sm text-slate-600">カテゴリ</label>
            <select
              value={form.category_id}
              onChange={(e) =>
                setForm((prev) => ({
                  ...prev,
                  category_id: parseInt(e.target.value, 10) || 0,
                }))
              }
              className="w-full rounded border border-slate-300 px-3 py-2 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
            >
              <option value={0}>選択してください</option>
              {categories.map((c) => (
                <option key={c.id} value={c.id}>
                  {c.name}
                </option>
              ))}
            </select>
          </div>
          <div>
            <label className="mb-1 block text-sm text-slate-600">金額（円）</label>
            <input
              type="number"
              required
              min="1"
              value={form.amount || ""}
              onChange={(e) =>
                setForm((prev) => ({
                  ...prev,
                  amount: parseInt(e.target.value, 10) || 0,
                }))
              }
              className="w-full rounded border border-slate-300 px-3 py-2 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
            />
          </div>
        </div>
        <div>
          <label className="mb-1 block text-sm text-slate-600">メモ</label>
          <input
            type="text"
            placeholder="任意"
            value={form.memo}
            onChange={(e) =>
              setForm((prev) => ({ ...prev, memo: e.target.value }))
            }
            className="w-full rounded border border-slate-300 px-3 py-2 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
          />
        </div>
        {submitError && (
          <p className="text-sm text-red-600">{submitError}</p>
        )}
        <div className="flex gap-2">
          <button
            type="submit"
            disabled={submitting}
            className="rounded bg-blue-600 px-4 py-2 font-medium text-white transition hover:bg-blue-700 disabled:opacity-50"
          >
            {submitting ? "登録中..." : "登録"}
          </button>
          <button
            type="button"
            onClick={() => router.push("/")}
            className="rounded border border-slate-300 px-4 py-2 font-medium text-slate-600 transition hover:bg-slate-50"
          >
            キャンセル
          </button>
        </div>
      </form>
    </section>
  );
}
