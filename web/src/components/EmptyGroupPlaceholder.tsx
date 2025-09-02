import Image from 'next/image';

export default function EmptyGroupPlaceholder() {
  return (
    <div className="text-center text-gray-500 mt-20 flex flex-col items-center">
    {/* ロゴ画像 */}
      <Image 
        src="/Images/Iogos/logo_gray.svg"
        alt="Enomo Icon Gray"
        width={70}
        height={95}
        className="mb-8"
      />
      <p className="font-bold">まだグループに加入していません！</p>
      <p>画面右上の設定から、新規グループの作成や</p>
      <p>既存グループへの加入が行えます</p>
    </div>
  );
}