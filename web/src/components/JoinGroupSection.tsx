import SectionTitle from '@/components/SectionTitle';
import Button from '@/components/Button';

type Props = {
  groupId: string;
  onGroupIdChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  onSubmit: () => void;
};

const JoinGroupSection = ({ groupId, onGroupIdChange, onSubmit }: Props) => {
  return (
    <div className="space-y-4">
      <SectionTitle>グループ参加</SectionTitle>
      <div className="flex flex-col items-center space-y-4"> 
        <input
                type="text"
                placeholder="参加するグループIDを入力（○桁）"
                value={groupId}
                onChange={onGroupIdChange}
                className="w-4/5 border-0 border-b-2 border-[#C73BA4] bg-transparent py-1.5 text-black placeholder:text-gray-400 focus:outline-none focus:ring-0 text-center"
            />
            <div className="pt-2 text-center">
                <Button onClick={onSubmit}>グループに参加する</Button>
            </div>
        </div>
    </div>
  );
};

export default JoinGroupSection;