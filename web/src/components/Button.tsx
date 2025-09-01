// ボタンの見た目と機能を定義するファイル
type ButtonProps = {
  children: React.ReactNode;
  onClick?: () => void; // ボタンがクリックされた時の処理
};

export default function Button({ children, onClick }: ButtonProps) {
  return (
    <button
      onClick={onClick}
      className="rounded-lg bg-[#C73BA4] px-8 py-1.5 text-white font-bold hover:bg-pink-600 transition-colors"
    >
      {children}
    </button>
  );
}