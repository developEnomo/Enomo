"use client";

import { useState, useEffect } from 'react';
import AppHeader from '@/components/AppHeader';
import GroupList from '@/components/GroupList';
import EmptyGroupPlaceholder from '@/components/EmptyGroupPlaceholder';
import { Group } from '@/types/group';
import Footer from '@/components/footer/Footer';

export default function GroupListPage() {
  const [groups, setGroups] = useState<Group[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const initActiveGroup = async () => {
      try {
        await fetch('/api/v1/groups/now', { method: 'DELETE', credentials: 'include' });
      } catch {}
    };

    const fetchGroups = async () => {
      setIsLoading(true);
      try {
        // APIエンドポイントにリクエストを送信
        const response = await fetch('/api/v1/users/list?limit=50&offset=0', { credentials: 'include' });

        if (!response.ok) {
          throw new Error(`API error: ${response.status}`);
        }

        const data = await response.json();

        // APIのレスポンス（スネークケース）をフロントエンドの型（キャメルケース）に変換
        if (data.groups && Array.isArray(data.groups)) {
          const formattedGroups: Group[] = data.groups.map((group: any) => ({
            id: group.group_id,
            name: group.name,
            memberCount: group.member_count,
          }));
          setGroups(formattedGroups);
        } else {
          setGroups([]); // グループがない場合は空配列をセット
        }

      } catch (error) {
        console.error('Failed to fetch groups:', error);
        setGroups([]); // エラー発生時も空のリストを表示
      } finally {
        setIsLoading(false);
      }
    };

    initActiveGroup().finally(fetchGroups);
  }, []);

  return (
    <div className="min-h-screen bg-white">
      <AppHeader />
      <div className="pt-20">
        <h1 className="text-4xl font-bold text-center">Group List</h1>
      </div>
      <main className="pt-10 pb-40"> {/* ヘッダーの高さ分だけ余白を確保 */}
        {isLoading ? (
          <p className="text-center mt-10 font-bold text-[#C73BA4]">Loading...</p>
        ) : (
          groups.length > 0 ? (
            <GroupList groups={groups} />
          ) : (
            <EmptyGroupPlaceholder />
          )
        )}
      </main>
      <Footer/>
    </div>
  );
}