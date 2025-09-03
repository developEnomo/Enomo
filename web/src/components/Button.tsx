type ButtonProps = {
  children: React.ReactNode;
  onClick?: () => void;
  variant?: 'primary' | 'danger';
  type?: 'button' | 'submit' | 'reset';
  className?: string;
};

export default function Button({
  children,
  onClick,
  variant = 'primary',
  type = 'button',
  className = '',
}: ButtonProps) {
  // ↓ ここから 'w-full' を削除します
  const baseClasses = 'rounded-lg px-8 py-1.5 font-bold transition-colors';

  const variantClasses = {
    primary: 'bg-[#C73BA4] text-white hover:bg-[#F66FD4]',
    secondary: 'bg-[#F66FD4] text-white hover:bg-[#A45B7D]',
    third: 'bg-[#861F6D] text-white hover:bg-[#861F6D]',
    danger: 'border-2 border-red-500 text-red-500 hover:bg-red-50',
  };

  return (
    <button
      type={type}
      onClick={onClick}
      className={`${baseClasses} ${variantClasses[variant]} ${className}`}
    >
      {children}
    </button>
  );
}