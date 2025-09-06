"use client";

import { useState, useEffect, use } from "react";
import { useRouter } from "next/navigation";

import ModalHeader from "@/components/ModalHeader";
import GroupIdSection from "@/components/GroupIdSection";
import GroupSettingsChangeSection from "@/components/GroupSettingsChangeSection";
import GroupDangerZoneSection from "@/components/GroupDangerZoneSection";
import Footer from "@/components/footer/Footer";

type PageProps = {
  params: Promise<{ groupId: string }>;
};

type GroupInfo = {
  id: string;
  name: string;
  updateFrequency: string;
  isOwner: boolean;
};

export default function GroupManagementPage({ params }: PageProps) {
  const { groupId } = use(params);
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

  //日時変換
  const hoursToFreq = (h?: number) => {
    if (!h) return "1d";
    if (h === 24) return "1d";
    if (h === 72) return "3d";
    if (h === 168) return "7d";
    return `${h}h`;
  };

  useEffect(() => {
    const load = async () => {
      setIsLoading(true);
      try {
        // 設定情報
        const res = await fetch(`/api/v1/groups/settings?group_id=${groupId}`, {
          credentials: "include",
        });
        if (!res.ok) throw new Error(String(res.status));
        const json = await res.json();

        // ★ ここがポイント：まず settings から、無ければ users/list から名前を解決
        let resolvedName: string | null = json.name ?? json.group_name ?? null;
        if (!resolvedName) {
          const listRes = await fetch(`/api/v1/users/list?limit=100&offset=0`, {
            credentials: "include",
          });
          if (listRes.ok) {
            const data = await listRes.json();
            const hit = Array.isArray(data?.groups)
              ? data.groups.find((g: any) => g.group_id === groupId)
              : null;
            if (hit?.name) resolvedName = hit.name as string;
          }
        }

        const next: GroupInfo = {
          id: json.id ?? groupId,
          name: resolvedName ?? "チーム",
          updateFrequency: json.update_frequency ?? "1d",
          isOwner: Boolean(json.is_owner ?? json.is_admin ?? false),
        };

        setGroupInfo(next);
        setEditedGroupName(next.name);
        setEditedFrequency(next.updateFrequency);
      } catch {
        const fallback: GroupInfo = {
          id: groupId,
          name: "チーム",
          updateFrequency: "1d",
          isOwner: false,
        };
        setGroupInfo(fallback);
        setEditedGroupName(fallback.name);
        setEditedFrequency(fallback.updateFrequency);
      } finally {
        setIsLoading(false);
      }
    };
    load();
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
