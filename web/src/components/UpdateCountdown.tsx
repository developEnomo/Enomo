type UpdateCountdownProps = {
  daysLeft: number;
};

export default function UpdateCountdown({ daysLeft }: UpdateCountdownProps) {
  return (
    <p className="text-right text-[#861F6D] font-semibold">
      次の更新まで {daysLeft} 日
    </p>
  );
}