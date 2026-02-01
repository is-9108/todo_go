"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

/**
 * Header は3つの画面（メイン・登録・編集）へ遷移するナビゲーションメニューを提供します。
 */
export default function Header() {
  const pathname = usePathname();

  const navItems = [
    { href: "/", label: "グラフ" },
    { href: "/register", label: "登録" },
    { href: "/transactions", label: "編集" },
  ];

  return (
    <header className="border-b border-slate-200 bg-white shadow-sm">
      <nav className="mx-auto flex max-w-4xl items-center gap-1 px-4 py-3">
        <h1 className="mr-6 text-lg font-bold text-slate-800">家計簿</h1>
        <div className="flex gap-1">
          {navItems.map((item) => {
            const isActive = pathname === item.href;
            return (
              <Link
                key={item.href}
                href={item.href}
                className={`rounded-lg px-4 py-2 text-sm font-medium transition ${
                  isActive
                    ? "bg-slate-800 text-white"
                    : "text-slate-600 hover:bg-slate-100 hover:text-slate-800"
                }`}
              >
                {item.label}
              </Link>
            );
          })}
        </div>
      </nav>
    </header>
  );
}
