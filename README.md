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

## 参考
- “cshum/imagor: libvipsを使用した、高速で安全な画像処理サーバーおよびGoライブラリ”. GitHub, [https://github.com/cshum/imagor?utm_source=chatgpt.com](https://github.com/cshum/imagor?utm_source=chatgpt.com), (Accessed 2026-07-13)
- “h2non/imaginary: Fast, simple, scalable, Docker-ready HTTP microservice for high-level image processing”. GitHub, [https://github.com/h2non/imaginary](https://github.com/h2non/imaginary), (Accessed 2026-07-13)
- “imgproxy/imgproxy: Fast and secure standalone server for resizing, processing, and converting images on the fly”. GitHub, [https://github.com/imgproxy/imgproxy?utm_source=chatgpt.com](https://github.com/imgproxy/imgproxy?utm_source=chatgpt.com), (Accessed 2026-07-13)
- “willnorris/imageproxy: A caching, resizing image proxy written in Go”. GitHub, [https://github.com/willnorris/imageproxy?utm_source=chatgpt.com](https://github.com/willnorris/imageproxy?utm_source=chatgpt.com), (Accessed 2026-07-13)
