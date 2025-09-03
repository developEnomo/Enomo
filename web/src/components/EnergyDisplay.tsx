import Image from 'next/image';

// 各エナジーレベルに対応する画像パスを定義
const energyImageMap: { [key: number]: string } = {
  1: '/Images/EnergyIcons/Energy_full_1.svg',
  2: '/Images/EnergyIcons/Energy_full_2.svg',
  3: '/Images/EnergyIcons/Energy_full_3.svg',
  4: '/Images/EnergyIcons/Energy_full_4.svg',
};

type EnergyDisplayProps = {
  // メンバー全員のエナジーレベルを配列で受け取る (例: [1, 4, 2, 3, ...])
  energyLevels: number[];
};

export default function EnergyDisplay({ energyLevels }: EnergyDisplayProps) {
  return (
    <div className="w-full">
      <h2 className="text-lg font-bold text-center mb-4 text-gray-800">みんなのエナジー</h2>
      
      {/* 横スクロールを実現するコンテナ */}
      <div className="flex items-center space-x-4 p-2 border-2 border-gray-200 rounded-lg overflow-x-auto flex-nowrap">
        {energyLevels.map((level, index) => {
          // エナジーレベルが1〜4の範囲外の場合はデフォルト画像などを指定（ここでは念のため）
          const imageUrl = energyImageMap[level] || '/Images/EnergyIcons/Energy_empty.svg';
          
          return (
            <div key={index} className="flex-shrink-0">
              <Image
                src={imageUrl}
                alt={`エナジーレベル ${level}`}
                width={60} // 画像のサイズを適宜調整してください
                height={60}
                priority={index < 5} // 最初の数枚を優先的に読み込む
              />
            </div>
          );
        })}
      </div>
    </div>
  );
}