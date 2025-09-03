import GroupPageHeader from '@/components/GroupPageHeader';
import MusicPlayer from '@/components/MusicPlayer';
import UpdateCountdown from '@/components/UpdateCountdown';
import EnergyDisplay from '@/components/EnergyDisplay';
import Footer from '@/components/footer/Footer';

// ページに渡されるパラメータの型定義
type PageProps = {
  params: { groupId: string };
};

// サーバーサイドでデータを取得する非同期関数
// TODO: 実際にはここでGoバックエンドのAPIをfetchします
async function getGroupData(groupId: string) {
  // --- ダミーデータ ---
  // APIから取得するデータの例
  const dummyData = {
    groupName: 'ひよこさんチーム',
    spotifyTrackId: '003vvx7Niy0yvhvHt4a68B', // 米津玄師 - LOSER のID
    daysUntilUpdate: 3,
    // 👇 メンバー10人分のエナジーレベル配列に変更
    memberEnergyLevels: [4, 2, 1, 3, 4, 3, 2, 1, 4, 3],
  };
  // --- ここまで ---

  // API連携の実装例:
  // const res = await fetch(`http://localhost:8080/api/groups/${groupId}`);
  // if (!res.ok) {
  //   throw new Error('Failed to fetch group data');
  // }
  // const data = await res.json();
  // return data;

  return dummyData;
}

export default async function GroupPlayPage({ params }: PageProps) {
  // URLの[groupId]を取得
  const { groupId } = params;

  // データを取得
  const groupData = await getGroupData(groupId);

  return (
    <div className="relative flex flex-col items-center min-h-screen bg-white text-gray-800 pt-20 pb-40">
      {/* ヘッダー */}
      <GroupPageHeader groupId={groupId} groupName={groupData.groupName} />

      <main className="w-full max-w-lg px-6 flex flex-col items-center gap-8">
        {/* ページタイトル */}
        <h2 className="text-2xl font-bold">Energy Song</h2>

        {/* 音楽プレーヤー */}
        <MusicPlayer trackId={groupData.spotifyTrackId} />

        {/* 更新カウントダウン */}
        <div className="w-full flex justify-end">
          <UpdateCountdown daysLeft={groupData.daysUntilUpdate} />
        </div>

        {/* エネルギー表示 (渡すpropを変更) */}
        <EnergyDisplay energyLevels={groupData.memberEnergyLevels} />
      </main>

      {/* フッターは別で作成されているとのことなので、
        ここにそのコンポーネントを配置してください。
        例: <Footer /> 
      */}
      <Footer />
    </div>
  );
}