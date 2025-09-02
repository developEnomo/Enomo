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
    e.preventDefault(); // ページの再読み込みを防ぐ
    setIsLoading(true);

    // TODO: ここにGoバックエンドへのAPIリクエストを実装します
    console.log('Submitting:', { email, password, username });
    
    // ダミーの待機時間
    await new Promise(resolve => setTimeout(resolve, 1000)); 
    
    alert(`アカウントを登録しました！\nEmail: ${email}\nUsername: ${username}`);
    setIsLoading(false);
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
          <Button>
            {isLoading ? '登録中...' : '新規登録する'}
          </Button>
        </div>
      </form>
    </div>
  );
}