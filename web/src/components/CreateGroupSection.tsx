import SectionTitle from '@/components/SectionTitle';
import Button from '@/components/Button';

type Props = {
  groupName: string;
  onGroupNameChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  frequency: string;
  onFrequencyChange: (e: React.ChangeEvent<HTMLSelectElement>) => void;
  onSubmit: () => void;
};

const CreateGroupSection = ({ groupName, onGroupNameChange, frequency, onFrequencyChange, onSubmit }: Props) => {
  const days = Array.from({ length: 10 }, (_, i) => i + 1);

  return (
    <div>
    <SectionTitle>新規グループ作成</SectionTitle>
    <div className="p-4 border border-[#861F6D] rounded-lg space-y-4">
      <input
        type="text"
        placeholder="グループ名を入力"
        value={groupName}
        onChange={onGroupNameChange}
        className="w-full border-0 border-b-2 border-[#C73BA4] bg-transparent py-1.5 text-black placeholder:text-gray-400 focus:outline-none focus:ring-0"
      />
      <div>
        <label className="text-sm text-[#861F6D]">プレイリスト更新頻度</label>
        <div className="flex items-center space-x-2 mt-1">
          <select
            value={frequency}
            onChange={onFrequencyChange}
            className="
              {/*appearance-none*/} /* デフォルトの矢印を非表示にする */
              w-20 /* 幅を少し広げると見栄えが良くなります */
              bg-white
              border-2 border-[#861F6D]
              rounded-md py-2 pl-3 pr-8
              focus:outline-none focus:ring-2 focus:ring-[#C73BA4]
              text-black text-center
              bg-no-repeat bg-right
              /* ↓ URLエンコードしたSVGで矢印を指定 */
              bg-[url('data:image/svg+xml,%3csvg_xmlns=%22http://www.w3.org/2000/svg%22_fill=%22none%22_viewBox=%220_0_20_20%22%3e%3cpath_stroke=%22%23C73BA4%22_stroke-linecap=%22round%22_stroke-linejoin=%22round%22_stroke-width=%221.5%22_d=%22M6_8l4_4_4-4%22/%3e%3c/svg%3e')]
            "
          >
            {days.map((day) => (
              <option key={day} value={day}>
                {day}
              </option>
            ))}
          </select>
          <span className="text-black">日ごとにプレイリスト更新</span>
        </div>
      </div>
    </div>
    <div className="pt-2 text-center">
        <Button onClick={onSubmit} className="w-fit">
          グループを作成する
        </Button>
      </div>
    </div>
  );
};

export default CreateGroupSection;