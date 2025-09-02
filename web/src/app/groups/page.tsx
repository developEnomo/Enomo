"use client";

import { useState, useEffect } from 'react';
import AppHeader from '@/components/AppHeader';
import GroupList from '@/components/GroupList';
import EmptyGroupPlaceholder from '@/components/EmptyGroupPlaceholder';
import { Group } from '@/types/group';

// ---開発用のモックデータ---
// グループがある場合のデータ
const mockGroups: Group[] = [
  { id: '1', name: 'ひよこさんチーム', memberCount: 3 },
  { id: '2', name: 'ひよこさんチーム', memberCount: 3 },
  { id: '3', name: 'ひよこさんチーム', memberCount: 3 },
  { id: '4', name: 'ひよこさんチーム', memberCount: 3 },
  { id: '5', name: 'ひよこさんチーム', memberCount: 3 },
];

// グループがない場合のデータ
const emptyGroups: Group[] = [];
// ---ここまで---

export default function GroupListPage() {
  const [groups, setGroups] = useState<Group[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // TODO: ここでGoバックエンドからデータをフェッチする
    // 今はモックデータでシミュレーション
    const fetchGroups = async () => {
      setIsLoading(true);
      // 2秒待ってからデータをセット
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      // グループがある場合とない場合を切り替えてテストできます
      setGroups(mockGroups); 
      //setGroups(emptyGroups); 
      
      setIsLoading(false);
    };

    fetchGroups();
  }, []);

  return (
    <div className="min-h-screen bg-white">
      <AppHeader />
      <div className="pt-20">
        <h1 className="text-4xl font-bold text-center">Group List</h1>
      </div>
      <main className="pt-10"> {/* ヘッダーの高さ分だけ余白を確保 */}
        {isLoading ? (
          <p className="text-center mt-10">読み込み中...</p>
        ) : (
          groups.length > 0 ? (
            <GroupList groups={groups} />
          ) : (
            <EmptyGroupPlaceholder />
          )
        )}
      </main>
    </div>
  );
}