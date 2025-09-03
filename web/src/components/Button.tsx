"use client";
import React from "react";

type ButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: "primary" | "secondary" | "third" | "danger";
};

export default function Button({ children, variant = "primary", className = "", ...rest }: ButtonProps) {
  const base = "rounded-lg px-8 py-1.5 font-bold transition-colors";
  const variants = {
    primary: "bg-[#C73BA4] text-white hover:bg-[#F66FD4]",
    secondary: "bg-[#F66FD4] text-white hover:bg-[#A45B7D]",
    third: "bg-[#861F6D] text-white hover:bg-[#861F6D]",
    danger: "border-2 border-red-500 text-red-500 hover:bg-red-50",
  } as const;
  return (
    <button {...rest} className={`${base} ${variants[variant]} ${className}`}>
      {children}
    </button>
  );
}