"use client";

import { useState } from 'react';
import { HiEye, HiEyeOff } from 'react-icons/hi';

type InputProps = {
  label: string;
  type: 'email' | 'password' | 'text';
  placeholder: string;
  value: string; // valueプロパティを追加
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void; // onChangeプロパティを追加
};

export default function Input({ label, type, placeholder, value, onChange }: InputProps) {
  const [isPasswordVisible, setIsPasswordVisible] = useState(false);

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
          value={value} // valueをinput要素に設定
          onChange={onChange} // onChangeをinput要素に設定
          className="
            w-full
            border-0
            border-b-2
            border-[#C73BA4]
            bg-transparent
            py-2
            placeholder:text-gray-400
            focus:outline-none
            focus:ring-0
          "
          required // 入力を必須にする
        />
        {type === 'password' && (
          <button type="button" onClick={togglePasswordVisibility} className="absolute inset-y-0 right-0 flex items-center pr-3 text-[#C73BA4]">
            {isPasswordVisible ? <HiEye className="h-5 w-5" /> : <HiEyeOff className="h-5 w-5" />}
          </button>
        )}
      </div>
    </div>
  );
}