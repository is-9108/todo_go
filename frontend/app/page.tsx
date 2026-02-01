"use client";

import { useEffect, useState } from "react";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import { getTransactions, type Transaction } from "@/lib/api";

/**
 * メイン画面: カテゴリごとに収入と支出を分けてグラフ表示します。
 */
export default function Home() {
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchData = async () => {
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
    fetchData();
  }, []);

  // カテゴリ別に収入・支出を集計
  const chartData = (() => {
    const byCategory: Record<
      string,
      { category: string; 収入: number; 支出: number }
    > = {};

    for (const t of (Array.isArray(transactions) ? transactions : [])) {
      const name = t.category?.name ?? "その他";
      if (!byCategory[name]) {
        byCategory[name] = { category: name, 収入: 0, 支出: 0 };
      }
      if (t.type === "income" || t.amount > 0) {
        byCategory[name].収入 += Math.abs(t.amount);
      } else {
        byCategory[name].支出 += Math.abs(t.amount);
      }
    }

    return Object.values(byCategory).sort(
      (a, b) => b.収入 + b.支出 - (a.収入 + a.支出)
    );
  })();

  const formatYAxis = (value: number) => {
    if (value >= 10000) return `${value / 10000}万`;
    return value.toString();
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

  const list = Array.isArray(transactions) ? transactions : [];
  if (list.length === 0) {
    return (
      <section className="rounded-lg bg-white p-8 shadow">
        <h2 className="mb-4 text-xl font-semibold text-slate-700">
          カテゴリ別 収支グラフ
        </h2>
        <p className="text-slate-500">
          データがありません。登録画面から収支を追加してください。
        </p>
      </section>
    );
  }

  return (
    <section className="space-y-8">
      <h2 className="text-xl font-semibold text-slate-700">
        カテゴリ別 収支グラフ
      </h2>

      <div className="rounded-lg bg-white p-6 shadow">
        <div className="mb-4 flex items-center justify-between">
          <h3 className="text-lg font-medium text-slate-600">
            収入・支出の内訳（円）
          </h3>
        </div>
        <div className="h-[400px]">
          <ResponsiveContainer width="100%" height="100%">
            <BarChart data={chartData} margin={{ top: 20, right: 30, left: 20, bottom: 60 }}>
              <CartesianGrid strokeDasharray="3 3" stroke="#e2e8f0" />
              <XAxis
                dataKey="category"
                tick={{ fontSize: 12 }}
                angle={-45}
                textAnchor="end"
                height={80}
              />
              <YAxis tickFormatter={formatYAxis} tick={{ fontSize: 12 }} />
              <Tooltip
                formatter={(value, name) => [
                  `¥${(value ?? 0).toLocaleString()}`,
                  name,
                ]}
                labelFormatter={(label) => `カテゴリ: ${label}`}
              />
              <Legend />
              <Bar dataKey="収入" fill="#10b981" radius={[4, 4, 0, 0]} />
              <Bar dataKey="支出" fill="#ef4444" radius={[4, 4, 0, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </div>
      </div>

      {/* サマリー */}
      <div className="grid gap-4 sm:grid-cols-2">
        <div className="rounded-lg bg-emerald-50 p-4 shadow">
          <p className="text-sm font-medium text-emerald-700">収入合計</p>
          <p className="text-2xl font-bold text-emerald-600">
            ¥
            {list
              .filter((t) => t.type === "income" || t.amount > 0)
              .reduce((sum, t) => sum + Math.abs(t.amount), 0)
              .toLocaleString()}
          </p>
        </div>
        <div className="rounded-lg bg-rose-50 p-4 shadow">
          <p className="text-sm font-medium text-rose-700">支出合計</p>
          <p className="text-2xl font-bold text-rose-600">
            ¥
            {list
              .filter((t) => t.type === "expense" || t.amount < 0)
              .reduce((sum, t) => sum + Math.abs(t.amount), 0)
              .toLocaleString()}
          </p>
        </div>
      </div>
    </section>
  );
}
