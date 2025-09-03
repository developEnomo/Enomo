"use client";

import { useState, useEffect, use } from "react";
import { useRouter } from 'next/navigation';

// 必要なコンポーネントをインポート
import ModalHeader from "@/components/ModalHeader";
import AppLogo from "@/components/AppLogo";
import GroupIdSection from "@/components/GroupIdSection";
import GroupSettingsChangeSection from "@/components/GroupSettingsChangeSection";
import GroupDangerZoneSection from "@/components/GroupDangerZoneSection";
import Footer from "@/components/footer/Footer";

// 型定義
type PageProps = {
  params: { groupId: string };
};
type GroupInfo = {
  id: string;
  name: string;
  updateFrequency: string;
  isOwner: boolean;
};

export default function GroupManagementPage({ params }: PageProps) {
  const router = useRouter();
  const { groupId } = use(params);

  // 状態管理
  const [groupInfo, setGroupInfo] = useState<GroupInfo | null>(null);
  const [editedGroupName, setEditedGroupName] = useState("");
  const [editedFrequency, setEditedFrequency] = useState("1d");
  const [isLoading, setIsLoading] = useState(true);

  // データ取得
  useEffect(() => {
    const fetchGroupInfo = async () => {
      setIsLoading(true);
      // ダミーデータ
      const dummyData: GroupInfo = {
        id: groupId,
        name: "ひよこさんチーム",
        updateFrequency: "1d",
        isOwner: true,
      };
      await new Promise(resolve => setTimeout(resolve, 500));

      setGroupInfo(dummyData);
      setEditedGroupName(dummyData.name);
      setEditedFrequency(dummyData.updateFrequency);
      setIsLoading(false);
    };

    fetchGroupInfo();
  }, [groupId]);

  // イベントハンドラ
  const handleUpdateSettings = () => { alert(`設定を更新: ${editedGroupName}`); };
  const handleLeaveGroup = () => { if (confirm("本当に脱退しますか？")) alert("脱退しました"); };
  const handleDeleteGroup = () => { if (confirm("本当に削除しますか？")) alert("削除しました"); };

  if (isLoading || !groupInfo) {
    return (
      <div className="fixed inset-0 bg-white bg-opacity-50 flex items-center justify-center">
        <div className="bg-white rounded-lg p-4">
          <p>読み込み中...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white min-h-screen">
      <div className="px-4 pb-8">
        
        {/* 1. settingページと同じModalHeaderを使用 */}
        <ModalHeader />
        
        {/* 2. 中央揃えのタイトルを追加 */}
        <div className="text-center">
          <h1 className="text-2xl font-bold text-black">{groupInfo.name}</h1>
        </div>

        <main className="space-y-10 pt-8">
          <GroupIdSection groupId={groupInfo.id} />

          <GroupSettingsChangeSection
            groupName={editedGroupName}
            onNameChange={(e) => setEditedGroupName(e.target.value)}
            frequency={editedFrequency}
            onFrequencyChange={(e) => setEditedFrequency(e.target.value)}
            onSubmit={handleUpdateSettings}
          />
          
          <GroupDangerZoneSection
            isOwner={groupInfo.isOwner}
            onLeave={handleLeaveGroup}
            onDelete={handleDeleteGroup}
          />
        </main>
        <Footer />
      </div>
    </div>
  );
}