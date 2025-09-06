"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";

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
  const router = useRouter();

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

  const handleLeaveGroup = async() => {
    if (!confirm("本当に脱退しますか？")) return;
    setIsLoading(true);
    try {
      const res = await fetch("/api/v1/groups/leave", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ group_id: groupId }),
      });
      if (!res.ok) {
        const msg = await res.text();
        alert(`脱退に失敗しました（${res.status}）\n${msg}`);
        setIsLoading(false);
        return;
      }
      await fetch("/api/v1/groups/now", { method: "DELETE", credentials: "include" });
      router.push("/groups");
      router.refresh();
    } catch {
      alert("通信エラーが発生しました");
      setIsLoading(false);
    }
  };

  const handleDeleteGroup = async () => {
    if (!confirm("本当に削除しますか？この操作は元に戻せません。")) return;
    setIsLoading(true);
    try {
      const res = await fetch("/api/v1/groups/delete", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ group_id: groupId }),
      });
      if (!res.ok) {
        const msg = await res.text();
        const human =
          res.status === 403 ? "削除権限がありません。" : `削除に失敗しました（${res.status}）`;
        alert(`${human}\n${msg}`);
        setIsLoading(false);
        return;
      }
      await fetch("/api/v1/groups/now", { method: "DELETE", credentials: "include" });
      router.push("/groups");
      router.refresh();
    } catch {
      alert("通信エラーが発生しました");
      setIsLoading(false);
    }
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
