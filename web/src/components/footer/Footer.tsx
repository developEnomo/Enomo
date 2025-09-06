import EnergyButton from "./EnergyButton";
const Footer = () => {
  return (
    <div className="fixed bottom-0 w-full h-[60px] z-10">
      {/* 背景レイヤー */}
      <div className="absolute bottom-0 w-full h-[60px] bg-[#C73BA4] z-0" />

      {/* ボタンレイヤー */}
      <div className="fixed bottom-0 w-full flex flex-col-reverse items-center justify-center py-4">
        <EnergyButton />
      </div>
    </div>
  );
};
export default Footer;
