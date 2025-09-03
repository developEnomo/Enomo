import SectionTitle from '@/components/SectionTitle';
import Button from '@/components/Button';
import { FaCopy } from 'react-icons/fa';

type Props = {
  groupId: string;
};

const GroupIdSection = ({ groupId }: Props) => {
  const handleCopy = () => {
    navigator.clipboard.writeText(groupId)
      .then(() => {
        alert("グループIDをコピーしました！");
      })
      .catch(err => {
        console.error('コピーに失敗しました: ', err);
        alert("コピーに失敗しました。");
      });
  };

  return (
    <div className="space-y-2">
      <SectionTitle>グループID</SectionTitle>
      {/* このdivのクラスを変更 */}
      <div className="flex flex-col items-center gap-4 p-4 bg-white rounded-lg">
        <span className="text-lg font-mono text-[#861F6D] tracking-wider">{groupId}</span>
        <Button onClick={handleCopy} className="!px-4 !py-2 flex items-center gap-2 w-fit">
          <span>コピーする</span>
          <FaCopy />
        </Button>
      </div>
    </div>
  );
};

export default GroupIdSection;