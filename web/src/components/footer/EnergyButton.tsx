"use client";
import { useState } from "react";
import Image from "next/image";
import React from "react";

const EnergyButton = () => {
  // メニュー開閉状態
  const [menuOpen, setMenuOpen] = useState(false);
  // 選択されたエナジー（Face番号）
  const [selected, setSelected] = useState<number | null>(null);

  // メニュー開閉トグル
  const EnergyButtonClick = () => {
    setMenuOpen((prev) => !prev);
  };

  // 表示する画像のロジック
  const getEnergyImageSrc = () => {
    if (menuOpen) {
      return selected !== null
        ? `/Images/EnergyIcons/Energy_fill_${selected}.svg`
        : "/Images/EnergyIcons/Energy_empty_push.svg";
    }
    return selected !== null
      ? `/Images/EnergyIcons/Energy_fill_${selected}.svg`
      : "/Images/EnergyIcons/Energy_empty.svg";
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
            <button
              key={faceId}
              onClick={() => {
                setSelected(faceId);
                setMenuOpen(false); // 選択後にメニューを閉じる
              }}
            >
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