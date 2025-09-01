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
  return(
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
  

  );
  
};
export default EnergyButton;
