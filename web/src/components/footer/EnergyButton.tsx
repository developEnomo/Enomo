"use client";
import { useState, useEffect } from "react";
import Image from "next/image";
import React from "react";

const EnergyButton = () => {
  const [menuOpen, setMenuOpen] = useState(false);
  const [selected, setSelected] = useState<number | null>(null);

  // 初回レンダー時に localStorage から選択状態を読み込む
  useEffect(() => {
    const stored = localStorage.getItem("selectedFaceId");
    if (stored !== null) {
      setSelected(Number(stored));
    }
  }, []);

  // 選択されたら localStorage に保存
  const handleSelect = async (faceId: number) => {
    setSelected(faceId);

    // APIに反映
    const response = await fetch('/api/v1/users/energy', {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include', // Cookieのセッションを送る
      body: JSON.stringify({ energy_value: faceId }),
    });

    if (!response.ok) {
      alert('エナジー更新に失敗しました');
      return;
    }

    localStorage.setItem('selectedFaceId', String(faceId));
    setMenuOpen(false);
  };

  const EnergyButtonClick = () => {
    setMenuOpen((prev) => !prev);
  };

  const getEnergyImageSrc = () => {
    if (menuOpen) {
      return selected !== null
        ? `/Images/EnergyIcons/Energy_full_${selected}.png`
        : "/Images/EnergyIcons/Energy_empty_push.png";
    }
    return selected !== null
      ? `/Images/EnergyIcons/Energy_full_${selected}.png`
      : "/Images/EnergyIcons/Energy_empty.png";
  };

  return (
    <>
      <button onClick={EnergyButtonClick}>
        <Image
          src={getEnergyImageSrc()}
          alt="energy"
          width={100}
          height={100}
          className="hover:scale-110 transition-transform active:scale-95 transition-transform duration-150"
        />
      </button>

      {menuOpen && (
        <div className="flex justify-center items-center border-2 border-[#C73BA4] px-4 py-2 m-4 bg-white rounded-lg gap-x-4">
          {[1, 2, 3, 4].map((faceId) => (
            <button key={faceId} onClick={() => handleSelect(faceId)}>
              <Image
                src={`/Images/EnergyIcons/Face${faceId}.svg`}
                width={50}
                height={50}
                alt={`face${faceId}`}
                className={`hover:scale-110 transition-transform ${
                  selected === faceId ? "ring-2 ring-[#C73BA4] rounded-full" : ""
                }`}
              />
            </button>
          ))}
        </div>
      )}
    </>
  );
};

export default EnergyButton;