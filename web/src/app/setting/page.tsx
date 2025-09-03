"use client";

import { useState } from 'react';
import ModalHeader from '@/components/ModalHeader';
import UsernameChangeSection from '@/components/UsernameChangeSection';
import JoinGroupSection from '@/components/JoinGroupSection';
import CreateGroupSection from '@/components/CreateGroupSection';
import AccountManagementSection from '@/components/AccountManagementSection';
import AppLogo from '@/components/AppLogo';
import Footer from '@/components/footer/Footer';

export default function SettingPage() {
  // --- 状態管理 ---
  const [username, setUsername] = useState('えねも えね男'); // 初期値はAPIから取得
  const [groupId, setGroupId] = useState('');
  const [newGroupName, setNewGroupName] = useState('');
  const [updateFrequency, setUpdateFrequency] = useState('daily');

  // --- イベントハンドラ ---
  const handleUsernameChange = () => {
    // TODO: ユーザー名変更APIを呼び出す
    alert(`ユーザー名を「${username}」に変更します。`);
  };

  const handleJoinGroup = () => {
    // TODO: グループ参加APIを呼び出す
    if (!groupId) {
      alert('グループIDを入力してください。');
      return;
    }
    alert(`グループID「${groupId}」に参加します。`);
  };

  const handleCreateGroup = () => {
    // TODO: グループ作成APIを呼び出す
    if (!newGroupName) {
      alert('グループ名を入力してください。');
      return;
    }
    alert(`グループ名「${newGroupName}」で作成します。\n更新頻度: ${updateFrequency}`);
  };

  const handleLogout = () => {
    // TODO: ログアウト処理
    if (confirm('本当にログアウトしますか？')) {
      alert('ログアウトしました。');
      // 例: ログインページにリダイレクト
      // router.push('/login');
    }
  };

  const handleDeleteAccount = () => {
    // TODO: アカウント削除処理
    if (confirm('アカウントを削除すると元に戻せません。本当に削除しますか？')) {
      alert('アカウントを削除しました。');
      // 例: トップページなどにリダイレクト
      // router.push('/');
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