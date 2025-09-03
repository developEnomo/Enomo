import SectionTitle from '@/components/SectionTitle';
import Button from '@/components/Button';

type Props = {
  isOwner: boolean;
  onLeave: () => void;
  onDelete: () => void;
};

const GroupDangerZoneSection = ({ isOwner, onLeave, onDelete }: Props) => {
  return (
    <div className="space-y-2">
      <div className="text-[#861F6D]">
        <SectionTitle>グループの脱退</SectionTitle>
      </div>
      <div className="flex flex-col items-center space-y-4 pt-2">
        <Button onClick={onLeave} className="w-full max-w-xs text-[#861F6D]">
          このグループを脱退する
        </Button>
        {isOwner && (
          <>
            <p className="text-sm text-red-600 text-center">
              グループの削除は作成者のみ実行できます
            </p>
            <Button onClick={onDelete} variant="danger" className="w-full max-w-xs">
              このグループを削除する
            </Button>
          </>
        )}
      </div>
    </div>
  );
};

export default GroupDangerZoneSection;