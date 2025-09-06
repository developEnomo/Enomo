"use client";
import { useEffect, useState } from "react";
import GroupPageHeader from '@/components/GroupPageHeader';
import MusicPlayer from '@/components/MusicPlayer';
import UpdateCountdown from '@/components/UpdateCountdown';
import EnergyDisplay from '@/components/EnergyDisplay';
import Footer from '@/components/footer/Footer';

// ページに渡されるパラメータの型定義
type PageProps = {
  params: { groupId: string };
};

type GroupData = {
  groupName: string;
  spotifyTrackId: string;
  daysUntilUpdate: number;
  memberEnergyLevels: number[];
};

export default function GroupDetailPage({ params }: PageProps) {
  const { groupId } = params;

  const [data, setData] = useState<GroupData | null>(null);
  const [loading, setLoading] = useState(true);

  // アクティブグループ登録
  useEffect(() => {
    const body = JSON.stringify({ group_id: groupId });
    fetch("/api/v1/groups/now", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
      body,
    }).catch(() => {});
  }, [groupId]);

  // データ取得（必要に応じてBEのAPIに差し替え）
  useEffect(() => {
    async function load() {
      setLoading(true);
      try {
        // TODO: 実APIに置き換える
        // const res = await fetch(`/api/v1/groups/${groupId}`, { credentials: "include" });
        // if (!res.ok) throw new Error(String(res.status));
        // const json = await res.json();

        const json: GroupData = {
          groupName: "ひよこさんチーム",
          spotifyTrackId: "003vvx7Niy0yvhvHt4a68B",
          daysUntilUpdate: 3,
          memberEnergyLevels: [4, 2, 1, 3, 4, 3, 2, 1, 4, 3],
        };
        setData(json);
      } catch {
        setData(null);
      } finally {
        setLoading(false);
      }
    }
    load();
  }, [groupId]);

  if (loading || !data) {
    return <div className="p-4">読み込み中...</div>;
  }

  return (
    <div className="relative flex flex-col items-center min-h-screen bg-white text-gray-800 pt-20 pb-40">
      <GroupPageHeader groupId={groupId} groupName={data.groupName} />

      <main className="w-full max-w-lg px-6 flex flex-col items-center gap-8">
        <h2 className="text-2xl font-bold">Energy Song</h2>
        <MusicPlayer trackId={data.spotifyTrackId} />
        <div className="w-full flex justify-end">
          <UpdateCountdown daysLeft={data.daysUntilUpdate} />
        </div>
        <EnergyDisplay energyLevels={data.memberEnergyLevels} />
      </main>

      <Footer />
    </div>
  );
}