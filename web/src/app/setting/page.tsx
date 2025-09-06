"use client";

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import ModalHeader from '@/components/ModalHeader';
import UsernameChangeSection from '@/components/UsernameChangeSection';
import JoinGroupSection from '@/components/JoinGroupSection';
import CreateGroupSection from '@/components/CreateGroupSection';
import AccountManagementSection from '@/components/AccountManagementSection';
import Footer from '@/components/footer/Footer';

// Helper function to convert frequency string to hours
const frequencyToHours = (freq: string): number => {
  switch (freq) {
    case '12h': return 12;
    case '1d': return 24;
    case '2d': return 48;
    case '3d': return 72;
    case '4d': return 96;
    case '5d': return 120;
    case '6d': return 144;
    case '7d': return 168;
    default: return 24; // Default to 1 day (24 hours)
  }
};

export default function SettingPage() {
  const router = useRouter();

  // --- 状態管理 ---
  const [username, setUsername] = useState('');
  const [initialUsername, setInitialUsername] = useState(''); // 更新失敗時に戻すための初期値
  const [groupId, setGroupId] = useState('');
  const [newGroupName, setNewGroupName] = useState('');
  const [updateFrequency, setUpdateFrequency] = useState('1d'); // API仕様に合わせて初期値を'1d'に変更

  // --- 初期データ取得 ---
  useEffect(() => {
    const fetchUserData = async () => {
      try {
        const response = await fetch('/api/v1/users/me', { credentials: 'include' });
        if (!response.ok) {
          // 認証エラーなどの場合はログインページにリダイレクト
          router.push('/login');
          return;
        }
        const userData = await response.json();
        setUsername(userData.display_name);
        setInitialUsername(userData.display_name);
      } catch (error) {
        console.error('Failed to fetch user data:', error);
        alert('ユーザー情報の取得に失敗しました。');
      }
    };
    fetchUserData();
  }, [router]);

  // --- イベントハンドラ ---
  const handleUsernameChange = async () => {
    try {
      const response = await fetch('/api/v1/users/rename', {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ display_name: username }),
      });
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'ユーザー名の変更に失敗しました。');
      }
      const updatedUser = await response.json();
      setInitialUsername(updatedUser.display_name);
      alert('ユーザー名を変更しました。');
    } catch (error) {
      console.error('Error updating username:', error);
      alert(error instanceof Error ? error.message : 'エラーが発生しました。');
      setUsername(initialUsername); // 失敗したら元の名前に戻す
    }
  };

  const handleJoinGroup = async () => {
    if (!groupId) {
      alert('グループIDを入力してください。');
      return;
    }
    try {
      const response = await fetch('/api/v1/groups/add', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ group_id: groupId }),
      });
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'グループへの参加に失敗しました。');
      }
      alert(`グループID「${groupId}」に参加しました。`);
      setGroupId('');
      router.push('/groups'); // グループ一覧ページに遷移
    } catch (error) {
      console.error('Error joining group:', error);
      alert(error instanceof Error ? error.message : 'エラーが発生しました。');
    }
  };

  const handleCreateGroup = async () => {
    if (!newGroupName) {
      alert('グループ名を入力してください。');
      return;
    }
    try {
      // 1. グループを作成
      const createResponse = await fetch('/api/v1/groups/make', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ name: newGroupName }),
      });

      if (createResponse.status !== 201) { // 201 Created を期待
        const errorData = await createResponse.json();
        throw new Error(errorData.error || 'グループの作成に失敗しました。');
      }
      const newGroup = await createResponse.json();
      
      // 2. 作成したグループの更新頻度を設定
      const hours = frequencyToHours(updateFrequency);
      const updateResponse = await fetch('/api/v1/groups/settings/update', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include',
          body: JSON.stringify({
              group_id: newGroup.id,
              refresh_interval_hours: hours,
          }),
      });

      if (!updateResponse.ok) {
        // 更新頻度の設定に失敗してもグループ作成は成功しているので、その旨を伝える
        const errorData = await updateResponse.json();
        throw new Error(`グループは作成されましたが、更新頻度の設定に失敗しました: ${errorData.error}`);
      }

      alert(`グループ「${newGroup.name}」を作成しました。`);
      setNewGroupName('');
      router.push('/groups'); // グループ一覧ページに遷移

    } catch (error) {
      console.error('Error creating group:', error);
      alert(error instanceof Error ? error.message : 'エラーが発生しました。');
    }
  };

  const handleLogout = async () => {
    if (confirm('本当にログアウトしますか？')) {
      try {
        const response = await fetch('/api/v1/users/logout', { method: 'POST', credentials: 'include' });
        if (!response.ok) throw new Error('ログアウトに失敗しました。');
        alert('ログアウトしました。');
        router.push('/login');
      } catch (error) {
        console.error('Logout error:', error);
        alert(error instanceof Error ? error.message : 'エラーが発生しました。');
      }
    }
  };

  const handleDeleteAccount = async () => {
    if (confirm('アカウントを削除すると元に戻せません。本当に削除しますか？')) {
      try {
        const response = await fetch('/api/v1/users/delete', { method: 'POST', credentials: 'include' });
        if (!response.ok) {
          const errorData = await response.json();
          throw new Error(errorData.error || 'アカウントの削除に失敗しました。');
        }
        alert('アカウントを削除しました。');
        router.push('/');
      } catch (error) {
        console.error('Delete account error:', error);
        alert(error instanceof Error ? error.message : 'エラーが発生しました。');
      }
    }
  };

  return (
    <div className="bg-white min-h-screen pb-40">
      <div className="max-w-md mx-auto px-4">
        <ModalHeader />
        <div className="pt-10">
          <h1 className="text-4xl text-black font-bold text-center">Setting</h1>
        </div>
        <main className="space-y-10 py-4">
          <UsernameChangeSection
            username={username}
            onUsernameChange={(e) => setUsername(e.target.value)}
            onSubmit={handleUsernameChange}
          />
          <JoinGroupSection
            groupId={groupId}
            onGroupIdChange={(e) => setGroupId(e.target.value)}
            onSubmit={handleJoinGroup}
          />
          <CreateGroupSection
            groupName={newGroupName}
            onGroupNameChange={(e) => setNewGroupName(e.target.value)}
            frequency={updateFrequency}
            onFrequencyChange={(e) => setUpdateFrequency(e.target.value)}
            onSubmit={handleCreateGroup}
          />
          <AccountManagementSection
            onLogout={handleLogout}
            onDeleteAccount={handleDeleteAccount}
          />
        </main>
    
      </div>
          <Footer />
    </div>
  );
}