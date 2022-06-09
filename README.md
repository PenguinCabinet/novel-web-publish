# 📚novel-web-publish(nwp)
これは小説サイト(なろう等)に自動アップロードできるコマンドラインツールです。

# ⬇ダウンロード
[Windows](https://github.com/PenguinCabinet/novel-web-publish/releases/download/v0.0.12/nwp_windows_x86_64.zip)   
[Windows(ARM)](https://github.com/PenguinCabinet/novel-web-publish/releases/download/v0.0.12/nwp_windows_armv6.zip)   
[Linux](https://github.com/PenguinCabinet/novel-web-publish/releases/download/v0.0.12/nwp_linux_x86_64.tar.gz)

# 📒現在デプロイ先として対応中のサイト
* [小説家になろう](https://syosetu.com/)

# ⚠Warning
現在アルファー版で問題点があります。
* 現在、デプロイするとなろうの小説カテゴリーはその他になります。修正予定です。
* なろうのキーワードを設定する方法が現在ありません。修正予定です。

# 🚀Quick Start
```bash
/root/home/novel>nwp narou login 
Email>test@example.com
Password>
/root/home/novel>nwp new  #プロジェクトをカレントディレクトリに作成
Title>やっぱりベリーショートヘアがナンバーワン!!
/root/home/novel>ls
nwp.yml _summary.txt
/root/home/novel>cat nwp.yml
title: やっぱりベリーショートヘアがナンバーワン!!
deploys:
- narou
/root/home/novel>micro _summary.txt #お好きなエディタを使って_summary.txtにあらすじを書いてください
/root/home/novel>nwp add パン咥えて走ったらロング派の友人にダンプで突っ込まれ、異世界転生した件 #エピソードを作成
/root/home/novel>ls
nwp.yml _summary.txt パン咥えて走ったらロング派の友人にダンプで突っ込まれ、異世界転生した件.txt
/root/home/novel>cat パン咥えて走ったらロング派の友人にダンプで突っ込まれ、異世界転生した件.txt
---
title: パン咥えて走ったらロング派の友人にダンプで突っ込まれ、異世界転生した件
index: 1

---
/root/home/novel>micro パン咥えて走ったらロング派の友人にダンプで突っ込まれ、異世界転生した件.txt #メタデータの下にエピソードの本文を書いてください
/root/home/novel>nwp deploy
Deploying to なろう...
Successd
```

詳しいプロジェクトのファイル構成はexamplesにある、[サンプルプロジェクト「やっぱりベリーショートヘアがナンバーワン!!」](./examples/%E3%82%84%E3%81%A3%E3%81%B1%E3%82%8A%E3%83%99%E3%83%AA%E3%83%BC%E3%82%B7%E3%83%A7%E3%83%BC%E3%83%88%E3%83%98%E3%82%A2%E3%81%8C%E3%83%8A%E3%83%B3%E3%83%90%E3%83%BC%E3%83%AF%E3%83%B3!!/)をご覧ください。  
サンプルプロジェクト「やっぱりベリーショートヘアがナンバーワン」をなろうにデプロイした場合が、[こちらです](https://ncode.syosetu.com/n3082hr/)


## 説明
### なろうへログイン
```bash
/root/home/novel>nwp narou login 
Email>test@example.com
Password>
```
なろうアカウントのパスワード、メールアドレスを入力します。
これでsecretsフォルダにAPI Keyが保存されたnarou.jsonが作られます。  
セキュリティの関係上、secretsは外部に流出しないよう注意してください。

### プロジェクトの作成
```bash
/root/home/novel>nwp new  #プロジェクトをカレントディレクトリに作成
Title>やっぱりベリーショートヘアがナンバーワン!!
/root/home/novel>ls
nwp.yml _summary.txt
/root/home/novel>cat nwp.yml
title: やっぱりベリーショートヘアがナンバーワン!!
deploys:
- narou
/root/home/novel>micro _summary.txt #お好きなエディタを使って_summary.txtにあらすじを書いてください
```
これで現在のカレントディレクトリディレクトリにプロジェクトが作成されます。  
nwp.ymlは設定ファイルで、_summary.txtにはあらすじを書きます。  
文字コードはUTF-8です。

### エピソードの追加

```bash
/root/home/novel>nwp add パン咥えて走ったらロング派の友人にダンプで突っ込まれ、異世界転生した件 #エピソードを作成
/root/home/novel>ls
nwp.yml _summary.txt パン咥えて走ったらロング派の友人にダンプで突っ込まれ、異世界転生した件.txt
/root/home/novel>micro パン咥えて走ったらロング派の友人にダンプで突っ込まれ、異世界転生した件.txt #メタデータの下にエピソードの本文を書いてください
```
nwp addでエピソードを追加することができます。  
<<エピソード名>>.txtが作られるので、そのファイルにエピソードを書きます。   
先頭にはyamlの設定がつくことに注意してください。  
文字コードはUTF-8です。

### デプロイ
```bash
/root/home/novel>nwp deploy
Deploying to なろう...
Successd
```
nwp deployとすると、nwp.yamlのdeploysに記されたデプロイ先にデプロイがおこないます。  
デプロイ先の小説投稿サイトと同期がとられて、差分の分だけ追加されたりあるいは消去されます。
    
ちなみにデプロイ先の規定を満たしていないとエラーを吐きます  
エラー例
```bash
Deploying to なろう...
Error:エピソード「パン咥えて走ったらロング派の友人にダンプで突っ込まれ、異世界転生した件」の本文が200文字以下です。エピソードの本文は200文字より大きく70000文字より小さくなければなりません。
Failed
```
この場合、なろうでは一エピソードあたり200文字以上でなければならないため、デプロイに失敗しました。


