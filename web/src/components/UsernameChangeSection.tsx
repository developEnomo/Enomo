import SectionTitle from '@/components/SectionTitle';
import Button from '@/components/Button';

type Props = {
  username: string;
  onUsernameChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  onSubmit: () => void;
};

const UsernameChangeSection = ({ username, onUsernameChange, onSubmit }: Props) => {
  return (
    <div className="space-y-4">
      <SectionTitle>ユーザーネーム変更</SectionTitle>
      <div className="flex flex-col items-center space-y-4"> 
        <input
        type="text"
        value={username}
        onChange={onUsernameChange}
        className="w-1/2 border-0 border-b-2 border-[#C73BA4] bg-transparent py-1.5 text-black placeholder:text-gray-400 focus:outline-none focus:ring-0 text-center"
        />
        <Button onClick={onSubmit}>
          変更する
        </Button>
      </div>
    </div>
  );
};

export default UsernameChangeSection;