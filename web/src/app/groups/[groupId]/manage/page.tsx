"use client";

import { useState, useEffect, use } from "react";
import { useRouter } from "next/navigation";

import ModalHeader from "@/components/ModalHeader";
import GroupIdSection from "@/components/GroupIdSection";
import GroupSettingsChangeSection from "@/components/GroupSettingsChangeSection";
import GroupDangerZoneSection from "@/components/GroupDangerZoneSection";
import Footer from "@/components/footer/Footer";

type PageProps = { params: Promise<{ groupId: string }> };

type FrequencyKey = "12h" | "1d" | "2d" | "3d" | "4d" | "5d" | "6d" | "7d";

type GroupInfo = {
  id: string;
  name: string;
  updateFrequencyKey: FrequencyKey;
  isOwner: boolean;
};

const keyToHours: Record<FrequencyKey, number> = {
  "12h": 12, "1d": 24, "2d": 48, "3d": 72, "4d": 96, "5d": 120, "6d": 144, "7d": 168,
};
const hoursToKey = (h?: number): FrequencyKey => {
  const m: Record<number, FrequencyKey> = { 12:"12h",24:"1d",48:"2d",72:"3d",96:"4d",120:"5d",144:"6d",168:"7d" };
  return m[h ?? 24] ?? "1d";
};

export default function GroupManagementPage({ params }: PageProps) {
  const { groupId } = use(params);
  const router = useRouter();

  const [groupInfo, setGroupInfo] = useState<GroupInfo | null>(null);
  const [editedGroupName, setEditedGroupName] = useState("");
  const [editedFrequencyKey, setEditedFrequencyKey] = useState<FrequencyKey>("1d");
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);

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
      await new Promise((resolve) => setTimeout(resolve, 500));

      setGroupInfo(dummyData);
      setEditedGroupName(dummyData.name);
      setEditedFrequency(dummyData.updateFrequency);
      setIsLoading(false);
    };

    fetchGroupInfo();
  }, [groupId]);

  // イベントハンドラ
  const handleUpdateSettings = () => {
    alert(`設定を更新: ${editedGroupName}`);
  };
  const handleLeaveGroup = () => {
    if (confirm("本当に脱退しますか？")) alert("脱退しました");
  };
  const handleDeleteGroup = () => {
    if (confirm("本当に削除しますか？")) alert("削除しました");
  };

  if (isLoading || !groupInfo) {
    return (
      <div className="fixed inset-0 bg-white bg-opacity-50 flex items-center justify-center">
        <div className="bg-white rounded-lg p-4">
          <p className="text-center mt-10 font-bold text-[#C73BA4]">
            Loading...
          </p>
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
            frequency={editedFrequencyKey}
            onFrequencyChange={(e) => setEditedFrequencyKey(e.target.value as FrequencyKey)}
            onSubmit={handleSave}
          />

          <GroupDangerZoneSection
            isOwner={groupInfo.isOwner}
            onLeave={handleLeave}
            onDelete={handleDelete}
          />
        </main>

        <Footer />
      </div>
    </div>
  );
}
