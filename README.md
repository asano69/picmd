# picmd

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/asano69/picmd)

<img src="frontend/public/favicon.svg" width="100" align="right" />

- Markdwonに画像を貼り付けるためのURLを生成する画像アップローダです。クリップボード、ファイルピッカーから簡単に画像をアップロードできます。
- 対応画像は、png, jpeg, webp。(gif, svgは対応予定)。複数枚の画像のバッチアップロード可です。
- 画像はpng透過画像の場合はpngに、それ以外の場合はwebpに変換されます。このとき、EXIF情報は削除されます。
- 画像はアップロードの後、サーバ側で圧縮され、UUIDv7が振られて保存されます。
- 画像の表示回数はviews属性に保存され、SQLコマンドで集計し、閲覧回数が極端に少ないものを削除できます。
- 画像のアップロードには認証が必要ですが、アップロードされた画像はURLを知っている人なら誰でも見ることができます。

## Screenshot
![](./.github/assets/sample-01.png)

## Work in progress
- Copied!ボタンをクリックのあと3秒後に戻す。
- 現在のアップロードの成功枚数/合計枚数を上部に比較的大きく表示させて、アップロード中にあと何枚未完了かわかるようにする。
- もとのファイル名とアップロードのあとのファイル名の対応関係をcsvでコピーできるようにする

## Plan
- 画像にうめこまれたEXIF情報がどのような場合に削除されずに保存されてしまうかの可能性を検討する
- 対応画像を増やす
    - SVG画像対応
- 画像のサイズをサーチパラメータなどで指定できるようにする

### Tech Stack
- バックエンドはGo+PocketBase **v0.39+**、frontendは、solid.js + **tailwind v4**で書かれています。

## 参考
- “cshum/imagor: libvipsを使用した、高速で安全な画像処理サーバーおよびGoライブラリ”. GitHub, [https://github.com/cshum/imagor?utm_source=chatgpt.com](https://github.com/cshum/imagor?utm_source=chatgpt.com), (Accessed 2026-07-13)
- “h2non/imaginary: Fast, simple, scalable, Docker-ready HTTP microservice for high-level image processing”. GitHub, [https://github.com/h2non/imaginary](https://github.com/h2non/imaginary), (Accessed 2026-07-13)
- “imgproxy/imgproxy: Fast and secure standalone server for resizing, processing, and converting images on the fly”. GitHub, [https://github.com/imgproxy/imgproxy?utm_source=chatgpt.com](https://github.com/imgproxy/imgproxy?utm_source=chatgpt.com), (Accessed 2026-07-13)
- “willnorris/imageproxy: A caching, resizing image proxy written in Go”. GitHub, [https://github.com/willnorris/imageproxy?utm_source=chatgpt.com](https://github.com/willnorris/imageproxy?utm_source=chatgpt.com), (Accessed 2026-07-13)
