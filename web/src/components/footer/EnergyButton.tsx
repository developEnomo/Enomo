"use client";
import { useState } from "react";
import Image from "next/image";
import React from "react";

const EnergyButton = () => {
  //メニュー開閉
  const [menuOpen, setMenuOpen] = useState(false);
  //エナジー選択済みかどうか
  const [selected, setSelected] = useState<number | null>(null);
  //ボタンクリックごとにメニュー開閉
  const EnergyButtonClick = () => {
    setMenuOpen((prev) => !prev);
  };
  return (
    <>
      <button onClick={EnergyButtonClick}>
        <Image
          src={
            menuOpen
              ? "/Images/EnergyIcons/Energy_empty_push.svg"
              : selected !== null
              ? `/Images/EnergyIcons/Energy_fill_${selected}.svg`
              : "/Images/EnergyIcons/Energy_empty.svg"
          }
          alt="energy"
          width={100}
          height={100}
        />
      </button>
      {menuOpen && (
        <div className="flex justify-center items-center border-2 border-[#C73BA4] px-4 py-2 m-4 bg-white rounded-lg gap-x-4">
          <button>
            <Image
              src="/Images/EnergyIcons/Face1.svg"
              width={50}
              height={50}
              alt="face1"
            />
          </button>
          <button>
            <Image
              src="/Images/EnergyIcons/Face2.svg"
              width={50}
              height={50}
              alt="face2"
            />
          </button>
          <button>
            <Image
              src="/Images/EnergyIcons/Face3.svg"
              width={50}
              height={50}
              alt="face3"
            />
          </button>
          <button>
            <Image
              src="/Images/EnergyIcons/Face4.svg"
              width={50}
              height={50}
              alt="face4"
            />
          </button>
        </div>
      )}
    </>
  );
};
export default EnergyButton;
