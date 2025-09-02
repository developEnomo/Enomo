import Link from 'next/link';
import { HiCog } from 'react-icons/hi';
import { IoChevronBack, IoSettingsOutline } from 'react-icons/io5';

type GroupPageHeaderProps = {
  groupId: string;
  groupName: string;
};

export default function GroupPageHeader({ groupId, groupName }: GroupPageHeaderProps) {
  return (
    <header className="fixed top-0 left-0 right-0 z-10 flex items-center justify-between w-full max-w-lg mx-auto p-4 bg-white">
      {/* 戻るボタン */}
      <Link href="/groups" className="text-4xl text-[#C73BA4]">
        <IoChevronBack />
      </Link>

      {/* グループ名 */}
      <h1 className="text-lg font-bold text-center text-black">
        {groupName}
      </h1>

      {/* 設定ボタン */}
      <Link href={`/groups/${groupId}/manage`} className="text-4xl text-[#C73BA4]">
        <HiCog />
      </Link>
    </header>
  );
}