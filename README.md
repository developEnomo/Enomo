# Enomo (energy + mono)

## dockerの起動方法
---
DockerDesktop開きながら
```
docker-compose up --build
```
DockerDesktopみてコンテナが全部動いていたらOK（web-1, api-1, db-1）


## PRレビューの簡単なやり方
---
```
git fetch -a
git branch -a
```
をして存在するブランチを出す
出力例
```bash
ryooy@thinkpadryoya:~/nutmeg/Enomo$ git branch -a
  develop
  feat/koba/36-selection-algorithm
  feat/mitomen/27-LogoutandSession
  feat/mitomen/32-Make-Group
  feat/walt/29-register-page-component
  feat/walt/37-group-page-component　　
  remotes/origin/Feat/oniwa/00-add-icon-images     
  remotes/origin/HEAD -> origin/develop
  remotes/origin/develop
  remotes/origin/feat/koba/10-add-dockerfiles      
  remotes/origin/feat/koba/14-add-internal-directory
  remotes/origin/feat/koba/36-selection-algorithm  
  remotes/origin/feat/koba/6-add-server-dir        
  remotes/origin/feat/mitomen/15-Register-api      
  remotes/origin/feat/mitomen/27-LogoutandSession  
  remotes/origin/feat/mitomen/32-Make-Group  
```

レビューするとき、リモートの内容をローカルに落とさなきゃいけない  
以下のコマンドで落とす。
```
git checkout -b ブランチ名 origin\ブランチ名
```

例えば、issue36のPRが出ててPRのレビューするってなったら
```
git checkout -b feat/koba/36-selection-algorithm origin/feat/koba/36-selection-algorithm
```
でOK。

あとはフロントなら`localhost:3000`の指定されたエンドポイントにアクセスするとか。
バックなら`curl`使ってjson入れたときにちゃんとほしい応答が返ってくるかを確認するとか。

PR出す側が指定してあげてください。
