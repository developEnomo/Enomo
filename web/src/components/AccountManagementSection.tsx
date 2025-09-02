import SectionTitle from '@/components/SectionTitle';
import Button from '@/components/Button';

type Props = {
  onLogout: () => void;
  onDeleteAccount: () => void;
};

const AccountManagementSection = ({ onLogout, onDeleteAccount }: Props) => {
  return (
    <div className="space-y-4">
      <SectionTitle>アカウント管理</SectionTitle>
      <div className="space-y-2 text-center">
        <Button onClick={onLogout} variant="third">
          ログアウトする
        </Button>
        <Button onClick={onDeleteAccount} variant="danger">
          アカウントを削除する
        </Button>
      </div>
    </div>
  );
};

export default AccountManagementSection;