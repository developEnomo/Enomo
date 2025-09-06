type UpdateCountdownProps = {
  // `daysLeft` から `updateIntervalInHours` に変更
  updateIntervalInHours: number;
};

export default function UpdateCountdown({ updateIntervalInHours }: UpdateCountdownProps) {
  // 24時間以下であれば時間で表示、それ以上は日数で表示
  const isWithin24Hours = updateIntervalInHours <= 24;
  const displayText = isWithin24Hours
    ? `${updateIntervalInHours} 時間`
    : `${Math.ceil(updateIntervalInHours / 24)} 日`;

  return (
    <p className="text-right text-[#861F6D] font-semibold">
      次の更新まで {displayText}
    </p>
  );
}