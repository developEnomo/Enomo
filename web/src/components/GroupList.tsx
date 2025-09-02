import { Group } from '@/types/group';
import GroupListItem from './GroupListItem';

type GroupListProps = {
  groups: Group[];
};

export default function GroupList({ groups }: GroupListProps) {
  return (
    <div className="w-full max-w-md mx-auto">
      {groups.map((group) => (
        <GroupListItem key={group.id} group={group} />
      ))}
    </div>
  );
}