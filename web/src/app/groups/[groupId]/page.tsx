"use client";
import { useEffect, useState, use } from "react";
import GroupPageHeader from '@/components/GroupPageHeader';
import MusicPlayer from '@/components/MusicPlayer';
import UpdateCountdown from '@/components/UpdateCountdown';
import EnergyDisplay from '@/components/EnergyDisplay';
import Footer from '@/components/footer/Footer';

// ページに渡されるパラメータの型定義
type PageProps = {
  params: Promise<{ groupId: string }>;
};

type GroupData = {
  groupName: string;
  spotifyTrackId: string;
  // `daysUntilUpdate` を `updateIntervalInHours` に変更
  updateIntervalInHours: number;
  memberEnergyLevels: number[];
};

export default function GroupDetailPage({ params }: PageProps) {
  const { groupId } = use(params);

  const [data, setData] = useState<GroupData | null>(null);
  const [loading, setLoading] = useState(true);

  // アクティブグループ登録
  useEffect(() => {
    if (!groupId) return;
    const body = JSON.stringify({ group_id: groupId });
    fetch("/api/v1/groups/now", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
      body: body,
    }).catch((err) => {
      console.error("Failed to set active group:", err);
    });
  }, [groupId]);

  // データ取得
  useEffect(() => {
    async function load() {
      if (!groupId) return;
      setLoading(true);
      try {
        const [settingsRes, listRes, recoRes] = await Promise.all([
          fetch(`/api/v1/groups/settings?group_id=${groupId}`, { credentials: 'include' }),
          fetch(`/api/v1/groups/list?group_id=${groupId}&limit=50`, { credentials: 'include' }),
          fetch(`/api/v1/groups/${groupId}/recommendations?limit=1`, { credentials: 'include' })
        ]);

        if (!settingsRes.ok || !listRes.ok || !recoRes.ok) {
          throw new Error('Failed to fetch group data');
        }

        const settingsData = await settingsRes.json();
        const listData = await listRes.json();
        const recoData = await recoRes.json();
        
        const groupName = settingsData.group_name || "チーム";
        // 時間をそのまま取得
        const updateIntervalInHours = settingsData.refresh_interval_hours || 24;

        const memberEnergyLevels = listData.members.map((member: any) => member.energy_value);
        
        const spotifyTrackId = recoData.tracks?.[0]?.id ?? "003vvx7Niy0yvhvHt4a68B"; // フォールバック用のID

        setData({
          groupName,
          spotifyTrackId,
          updateIntervalInHours, // 日数ではなく時間をセット
          memberEnergyLevels,
        });

      } catch (error) {
        console.error("Failed to load group data:", error);
        setData(null);
      } finally {
        setLoading(false);
      }
    }
    load();
  }, [groupId]);

  if (loading) {
    return <div className="p-4 text-center mt-20 font-bold text-[#C73BA4]">読み込み中...</div>;
  }
  
  if (!data) {
    return <div className="p-4 text-center mt-20 font-bold text-red-500">グループ情報の取得に失敗しました。</div>;
  }

  return (
    <div className="relative flex flex-col items-center min-h-screen bg-white text-gray-800 pt-20 pb-40">
      <GroupPageHeader groupId={groupId} groupName={data.groupName} />

      <main className="w-full max-w-lg px-6 flex flex-col items-center gap-8">
        <h2 className="text-2xl font-bold">Energy Song</h2>
        <MusicPlayer trackId={data.spotifyTrackId} />
        <div className="w-full flex justify-end">
           {/* プロパティ名を変更して時間データを渡す */}
           <UpdateCountdown updateIntervalInHours={data.updateIntervalInHours} />
        </div>
        <EnergyDisplay energyLevels={data.memberEnergyLevels} />
      </main>

      <Footer />
    </div>
  );
}