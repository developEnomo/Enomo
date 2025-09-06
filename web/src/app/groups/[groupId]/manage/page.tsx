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
    const load = async () => {
      setIsLoading(true);
      try {
        // セッション紐付けは fire-and-forget
        if (groupId) {
          fetch("/api/v1/groups/now", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify({ group_id: groupId }),
          }).catch(() => {});
        }

        // 設定取得（8秒で打ち切り）
        if (!groupId) {
          // groupId が未解決のままならフォールバック
          setGroupInfo({ id: "", name: "チーム", updateFrequencyKey: "1d", isOwner: false });
          return;
        }
        const ac = new AbortController();
        const tm = setTimeout(() => ac.abort(), 8000);
        const res = await fetch(`/api/v1/groups/settings?group_id=${groupId}`, {
          credentials: "include",
          signal: ac.signal,
        });
        clearTimeout(tm);

        if (!res.ok) throw new Error(`GET /groups/settings ${res.status}`);
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
      } catch (e) {
        console.error("[manage] load error:", e);
        // 失敗しても必ずフォールバックをセット
        setGroupInfo({ id: groupId ?? "", name: "チーム", updateFrequencyKey: "1d", isOwner: false });
      } finally {
        // どの経路でもローディングは必ず終了
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

  const handleLeaveGroup = async () => {
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
  };

  // ← ここを "isLoading のみ" に変更（!groupInfo に依存しない）
  if (isLoading) {
    return (
      <div className="fixed inset-0 bg-white bg-opacity-50 flex items-center justify-center">
        <div className="bg-white rounded-lg p-4">
          <p className="text-center mt-10 font-bold text-[#C73BA4]">Loading...</p>
        </div>
      </div>
    );
  }

  // ロード終了時に groupInfo が null のままは想定外なので安全側
  if (!groupInfo) {
    return (
      <div className="p-6">
        <p>データの読み込みに失敗しました。リロードしてください。</p>
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
            frequency={editedFrequencyKey}
            onFrequencyChange={(e) => setEditedFrequencyKey(e.target.value as FrequencyKey)}
            onSubmit={handleSave}
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
