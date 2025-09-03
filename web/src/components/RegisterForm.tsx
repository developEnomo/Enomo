"use client";

import { useState } from 'react';
import Input from './Input';
import Button from './Button';

export default function RegisterForm() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [username, setUsername] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  // フォーム送信時の処理
  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsLoading(true);

    console.log("[register] submit start", { email, username });

    try {
      const response = await fetch("/api/v1/users/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          email,
          password,
          display_name: username,
          energy_value: 3,
        }),
      });

      const rawBody = await response.clone().text();
      console.log("[register] status", response.status);
      console.log("[register] body", rawBody);

      if (!response.ok) throw new Error(rawBody || `HTTP ${response.status}`);

      // 必要ならJSONを読む
      let user: any = {};
      try { user = await response.json(); } catch {}
      console.log("[register] parsed", user);

      window.location.assign("/login");
    } catch (err) {
      console.error("[register] error", err);
      alert(err instanceof Error ? err.message : String(err));
    } finally {
      setIsLoading(false);
      console.log("[register] submit end");
    }
  };


  return (
    <div className="w-full max-w-sm px-4">
      {/* ヘッダー */}
      <div className="text-center mt-20 mb-20">
        <h1 className="mt-10 mb-10 text-4xl font-bold text-black">Sign Up</h1>
        <p className="mt-4 font-bold text-[#861F6D]">新しいアカウントを作成</p>
      </div>

      {/* フォーム */}
      <form onSubmit={handleSubmit} className="mt-15 mb-15 space-y-8">
        <Input
          label="メールアドレス"
          type="email"
          placeholder="メールアドレスを入力"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />
        <Input
          label="パスワード"
          type="password"
          placeholder="パスワードを入力"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        <div>
          <Input
            label="ユーザーネーム"
            type="text"
            placeholder="ユーザーネームを入力"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
          <p className="mt-2 text-xs text-[#861F6D]">ユーザーネームは後から変更できます</p>
        </div>
        
        {/* 新規登録ボタン */}
        <div className="pt-8 flex justify-center">
          <Button type="submit" disabled={isLoading} aria-busy={isLoading}>
            {isLoading ? '登録中...' : '新規登録する'}
          </Button>
        </div>
      </form>
    </div>
  );
}