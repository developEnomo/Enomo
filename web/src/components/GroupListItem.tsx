import Link from 'next/link';
import { Group } from '@/types/group';
import { FaChevronRight } from 'react-icons/fa';

type GroupListItemProps = {
  group: Group;
};

export default function GroupListItem({ group }: GroupListItemProps) {
  return (
    <Link href={`/groups/${group.id}`} className="flex items-center justify-between w-full p-4 my-2">
      <div className="flex items-center">
        {/* TODO: グループアイコンを後で追加 */}
        <span className="text-lg text-black font-semibold">{group.name}</span>
      </div>
      <div className="flex items-center space-x-2 text-black">
        <span>{group.memberCount}</span>
        <div className="ml-3 text-[#C73BA4]">
          <FaChevronRight />
        </div>
      </div>
    </Link>
  );
}