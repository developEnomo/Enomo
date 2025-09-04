"use client";

import { useState } from 'react';
import Input from './Input';
import Button from './Button';
import Link from 'next/link';

export default function LoginForm() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleLogin = async () => {
    if (isLoading) return;
    setIsLoading(true);
    try {
      const response = await fetch("/api/v1/users/login", {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        alert(`ログイン失敗: ${errorData.message}`);
        return;
      }

      //const data = await response.json();
      //alert(`ログイン成功: ${data.message}`);
      window.location.assign("/groups");
    } catch (error) {
      console.error(error);
      alert('通信エラーが発生しました');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="w-full max-w-sm">
      {/* 見出し */}
      <div className="text-center mt-20 mb-20">
        <h1 className="text-4xl mt-10 mb-10 font-bold text-black">Log In</h1>
        <p className="mt-10 text-[#861F6D] font-bold">既存のアカウントでログイン</p>
      </div>

      {/* 入力フォーム */}
      <div className="mt-15 mb-15 space-y-6">
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
      </div>

      {/* ログインボタンとリンク */}
      <div className="mt-10 flex flex-col items-center space-y-4">
        <div className="mt-8 mb-8">
          <Button onClick={handleLogin} disabled={isLoading} aria-busy={isLoading}>
            {isLoading ? 'ログイン中...' : 'ログインする'}
          </Button>
        </div>
        <Link href="/register" className="text-sm text-[#861F6D] mt-2 mb-2 underline hover:text-[#F66FD4]">
          新規登録はこちら
        </Link>
      </div>
    </div>
  );
}