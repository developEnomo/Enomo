"use client";

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import ModalHeader from '@/components/ModalHeader';
import UsernameChangeSection from '@/components/UsernameChangeSection';
import JoinGroupSection from '@/components/JoinGroupSection';
import CreateGroupSection from '@/components/CreateGroupSection';
import AccountManagementSection from '@/components/AccountManagementSection';
import Footer from '@/components/footer/Footer';

export default function SettingPage() {
  const router = useRouter();

  // --- 状態管理 ---
  const [username, setUsername] = useState('');
  const [initialUsername, setInitialUsername] = useState(''); // 更新失敗時に戻すための初期値
  const [groupId, setGroupId] = useState('');
  const [newGroupName, setNewGroupName] = useState('');
  const [updateFrequency, setUpdateFrequency] = useState('daily'); // この値は現在APIでは使用されません

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
      const response = await fetch('/api/v1/groups/make', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ name: newGroupName }),
      });
      if (response.status !== 201) { // 201 Created を期待
        const errorData = await response.json();
        throw new Error(errorData.error || 'グループの作成に失敗しました。');
      }
      const newGroup = await response.json();
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
        <Footer />
      </div>
    </div>
  );
}