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
<<<<<<< HEAD
    if (!groupId) return;
    fetch("/api/v1/groups/now", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
      body: JSON.stringify({ group_id: groupId }),
    }).catch(() => {});
  }, [groupId]);

  useEffect(() => {
    const load = async () => {
      setIsLoading(true);
      try {
        const res = await fetch(`/api/v1/groups/settings?group_id=${groupId}`, { credentials: "include" });
        if (!res.ok) throw new Error(String(res.status));
        const j = await res.json();

        const name = String(j.group_name ?? j.name ?? j.groupName ?? "チーム");
        const hours = Number(j.refresh_interval_hours ?? j.hours ?? 24);

        const info: GroupInfo = {
          id: groupId,
          name,
          updateFrequencyKey: hoursToKey(hours),
          isOwner: true, // 必要なら別APIで厳密化
        };
        setGroupInfo(info);
        setEditedGroupName(info.name);
        setEditedFrequencyKey(info.updateFrequencyKey);
      } catch {
        const fallback: GroupInfo = { id: groupId, name: "チーム", updateFrequencyKey: "1d", isOwner: false };
        setGroupInfo(fallback);
        setEditedGroupName(fallback.name);
        setEditedFrequencyKey(fallback.updateFrequencyKey);
      } finally {
        setIsLoading(false);
      }
    };
    load();
  }, [groupId]);

  const handleSave = async () => {
    if (!groupInfo) return;
    setIsSaving(true);
    try {
      const payload = {
        group_id: groupInfo.id,
        group_name: editedGroupName,
        refresh_interval_hours: keyToHours[editedFrequencyKey],
      };
      const res = await fetch("/api/v1/groups/settings/update", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify(payload),
      });
      const text = await res.text();
      if (!res.ok) {
        alert(`更新に失敗しました（${res.status}）\n${text}`);
        return;
      }
      setGroupInfo(prev => prev ? { ...prev, name: editedGroupName, updateFrequencyKey: editedFrequencyKey } : prev);
      alert("保存しました");
      router.refresh();
    } catch {
      alert("通信エラーが発生しました");
    } finally {
      setIsSaving(false);
    }
  };

  const handleLeave = async () => {
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

  const handleDelete = async () => {
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
        alert(`削除に失敗しました（${res.status}）\n${msg}`);
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
=======
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
>>>>>>> 8fe444288cc0557209819415a4bb30920ed8a2d1
  };

  if (isLoading || !groupInfo) {
    return (
      <div className="fixed inset-0 bg-white bg-opacity-50 flex items-center justify-center">
<<<<<<< HEAD
        <div className="bg-white rounded-lg p-4"><p>読み込み中...</p></div>
=======
        <div className="bg-white rounded-lg p-4">
          <p className="text-center mt-10 font-bold text-[#C73BA4]">
            Loading...
          </p>
        </div>
>>>>>>> 8fe444288cc0557209819415a4bb30920ed8a2d1
      </div>
    );
  }

  return (
    <div className="bg-white min-h-screen">
      <div className="px-4 pb-8">
<<<<<<< HEAD
        <ModalHeader />
=======
        {/* 1. settingページと同じModalHeaderを使用 */}
        <ModalHeader />

        {/* 2. 中央揃えのタイトルを追加 */}
>>>>>>> 8fe444288cc0557209819415a4bb30920ed8a2d1
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
