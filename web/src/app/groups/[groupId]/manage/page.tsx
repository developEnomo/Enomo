"use client";

import { useState, useEffect } from "react";

import ModalHeader from "@/components/ModalHeader";
import GroupIdSection from "@/components/GroupIdSection";
import GroupSettingsChangeSection from "@/components/GroupSettingsChangeSection";
import GroupDangerZoneSection from "@/components/GroupDangerZoneSection";
import Footer from "@/components/footer/Footer";

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
  const { groupId } = params;

  useEffect(() => {
    if (!groupId) return;
    const body = JSON.stringify({ group_id: groupId });
    fetch("/api/v1/groups/now", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
      body,
    });
  }, [groupId]);

  const [groupInfo, setGroupInfo] = useState<GroupInfo | null>(null);
  const [editedGroupName, setEditedGroupName] = useState("");
  const [editedFrequency, setEditedFrequency] = useState("1d");
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchGroupInfo = async () => {
      setIsLoading(true);
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
          <p>読み込み中...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white min-h-screen">
      <div className="px-4 pb-8">
        <ModalHeader />

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
