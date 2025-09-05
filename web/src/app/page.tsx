"use client";
import { useEffect } from "react";
import { useRouter } from "next/navigation";
import Image from "next/image";
import Footer from "@/components/footer/Footer";

export default function Home() {
  const router = useRouter();

  useEffect(() => {
    const timer = setTimeout(() => {
      router.push("/login"); // 遷移先のパスを指定
    }, 3000); // 3秒後

    return () => clearTimeout(timer); // クリーンアップ
  }, [router]);

  return (
    <div className="flex flex-col items-center justify-center min-h-screen gap-y-10">
      <Image
        src="/Images/Iogos/Logo_Enomo.svg"
        alt="Enomo"
        width={328}
        height={108}
      />
      <Image
        src="/Images/Iogos/Logo_Enomo_Icon.svg"
        alt="Enomo"
        width={80}
        height={106}
      />
      <p className="text-[#C73BA4] font-bold">Loading...</p>
    </div>
  );
}
