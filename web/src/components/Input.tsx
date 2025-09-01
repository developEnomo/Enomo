"use client"; // stateを使うため、ファイルの先頭にこの一行を追加

import { useState } from 'react';
import { HiEye, HiEyeOff } from 'react-icons/hi';

type InputProps = {
  label: string;
  type: 'email' | 'password' | 'text';
  placeholder: string;
};

export default function Input({ label, type, placeholder }: InputProps) {
  const [isPasswordVisible, setIsPasswordVisible] = useState(false);

  // パスワードの表示・非表示を切り替える
  const togglePasswordVisibility = () => {
    setIsPasswordVisible(!isPasswordVisible);
  };
  
  const inputType = type === 'password' && isPasswordVisible ? 'text' : type;

  return (
    <div className="w-full">
      <label className="text-sm font-medium text-[#861F6D]">{label}</label>
      <div className="relative">
        <input
          type={inputType}
          placeholder={placeholder}
          className="
            w-full
            border-0                  /* 1. 上左右の枠線をリセット */
            border-b-2                /* 2. 下線を2pxの太さに */
            border-[#C73BA4]           /* 3. 下線の色を常にピンクに */
            bg-transparent            /* 背景色を透明に */
            py-2                      
            placeholder:text-[rgb(156,163,175)] /* 4. プレースホルダーの色をRGBで指定 */
            focus:outline-none        /* フォーカス時の外枠を削除 */
            focus:ring-0              /* フォーカス時のリングを削除 */
          "
        />
        {type === 'password' && (
          // 目のアイコンをクリックすると表示が切り替わる
          <button type="button" onClick={togglePasswordVisibility} className="absolute inset-y-0 right-0 flex items-center pr-3 text-[#C73BA4]">
            {isPasswordVisible ? <HiEye className="h-5 w-5" /> : <HiEyeOff className="h-5 w-5" />}
          </button>
        )}
      </div>
    </div>
  );
}