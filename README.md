# picmd

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/asano69/picmd)

<img src="frontend/public/favicon.svg" width="100" align="right" />

- Markdwonに画像を貼り付けるためのURLを生成する画像アップローダ。
- クリップボード、ファイルピッカーから簡単に画像をアップロードできる。複数枚対応。
- 画像はアップロードの後、サーバ側で圧縮されて保存される。
- 画像のメタデータの管理にPocketbaseを使い、コンテンツの管理を容易にする。


## Screenshot
![](./.github/assets/sample-01.png)

## Work in progress


## Plan
- 画像にうめこまれたメタデータから個人情報を削除する方法の検討
- SVG画像対応
- 画像のサイズをサーチパラメータなどで指定できるようにする

### Tech Stack
- バックエンドはGo+PocketBase **v0.39+**、frontendは、solid.js + **tailwind v4**で書かれています。

