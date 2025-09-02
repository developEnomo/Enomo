import { FaBars } from 'react-icons/fa';
import Link from 'next/link';

export default function AppHeader() {
  return (
    <header className="fixed top-0 left-0 right-0 bg-white">
      <div className="max-w-md mx-auto flex items-center justify-between p-4">
        {/* 左側のスペース確保用 (アイコンなどがない場合) */}
        <div className="w-8"></div>
        
        
        {/* 設定へのリンクを持つハンバーガーメニュー */}
        <Link href="/setting" className="text-4xl text-[#C73BA4]">
          <FaBars />
        </Link>
      </div>
    </header>
  );
}