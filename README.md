# novel-web-publish(nwp)
これは小説サイト(なろう等)に自動アップロードできるコマンドラインツールです。

# 現在デプロイ先として対応中のサイト
* [小説家になろう](https://syosetu.com/)

# warning
* 現在、デプロイするとなろうのカテゴリーはその他になります。修正予定です。
* なろうのキーワードを設定する方法が現在ありません。修正予定です。

# Quick Start
```bash
/root/home/novel>nwp narou login 
Email>test@example.com
Password>
/root/home/novel>nwp new
Title>やっぱりベリーショートヘアがナンバーワン!!
/root/home/novel>ls
nwp.yml _summary.txt
/root/home/novel>cat nwp.yml
title: やっぱりベリーショートヘアがナンバーワン!!
deploys: ["narou"]
/root/home/novel>echo "超あらすじだけれど、俺は何の変哲もない高校二年生。友達がロング派なのでぶっ潰すぜ。ベリーショートヘアって素晴らしい。うぇーい">_summary.txt
/root/home/novel>nwp add "第一章パン咥えて走ったら交通事故で大変"
/root/home/novel>ls
nwp.yml _summary.txt 第一章パン咥えて走ったら交通事故で大変.txt
/root/home/novel>echo "遅刻を避けるべくパン咥えて走る俺。だが、そこにトラックが!!死に際に見たのは運転席にいるロング派の友人だった…。「クッソはめられた…異世界でなりあがって帰ってくるぜ」">>第一章パン咥えて走ったら交通事故で大変.txt
/root/home/novel>nwp deploy
Deploying to なろう...
Successd
```
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
/root/home/novel>nwp new
Title>やっぱりベリーショートヘアがナンバーワン!!
/root/home/novel>ls
nwp.yml _summary.txt
/root/home/novel>cat nwp.yml
title: やっぱりベリーショートヘアがナンバーワン!!
deploys: ["narou"]
/root/home/novel>echo "超あらすじだけれど、俺は何の変哲もない高校二年生。友達がロング派なのでぶっ潰すぜ。ベリーショートヘアって素晴らしい。うぇーい">_summary.txt
```
これで現在のワークディレクトリにプロジェクトが作成されます。  
nwp.ymlは設定ファイルで、_summary.txtにはあらすじを書きます。

### エピソードの追加

```bash
/root/home/novel>nwp add "第一章パン咥えて走ったら交通事故で大変"
/root/home/novel>ls
nwp.yml _summary.txt 第一章パン咥えて走ったら交通事故で大変.txt
/root/home/novel>echo "遅刻を避けるべくパン咥えて走る俺。だが、そこにトラックが!!死に際に見たのは運転席にいるロング派の友人だった…。「クッソはめられた…異世界でなりあがって帰ってくるぜ」">>第一章パン咥えて走ったら交通事故で大変.txt
```
nwp addでエピソードを追加することができます。  
<<エピソード名>>.txtが作られるので、そのファイルにエピソードを書きます。   
先頭にはyamlの設定がつくことに注意してください。

### デプロイ
```bash
/root/home/novel>nwp deploy
Deploying to なろう...
Successd
```
nwp deployとすると、nwp.yamlのdeploysに記されたデプロイ先にデプロイをおこないます。  
デプロイ先の小説投稿サイトと同期がとられて、差分の分だけ追加されたりあるいは消去されます。

