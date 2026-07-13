# Overview

- picmdからpocketbase+soild.jsをベースにしたpicmd2に移行しようとしている。
- 画像のメタデータの管理にPocketbaseを使い、コンテンツの管理を容易にすることが移行の目的。
- バックエンドはGo+PocketBase **v0.39+**、frontendは、solid.js + **tailwind v4**で書かれています。

## Rules

- データベースのマイグレーションはPocketBaseのWEB UIから行うのでマイグレーションコードを作成する必要はまったくない。
- When fixing bugs, add a failing regression test first.
- All errors are user-facing, so messages should be clear.
- Keep functions small and focused.
- Module files should re-export what's needed, hide implementation details.
- Don't persist changes to the database during drilling. Use the cache.
- Don't use timezones: dates are naive for a reason. Due dates etc. are more like the dates in a journal entry than precise points in time.

# Work in Progress

## 前提の確認

- `compress/webp.go` のロジック(リサイズ→品質を二分探索でWebPエンコード)はほぼそのまま流用できます。

## Phase 0: 圧縮ロジックの移植先を作る

- `compress/webp.go` の中身を `internal/media`（または`internal/compress`）としてpicmd2側にコピー。
- picmd2の`go.mod`に `chai2010/webp` と `golang.org/x/image` を追加。
- ロジック自体はテストしやすい形（`io.Reader`を受けて`[]byte`を返す関数）のまま維持。CLAUDE.mdの「関数は小さく」に沿って、resize/encodeは既に分離されているのでそのまま活かせます。

## Phase 1: PocketBaseフックで圧縮を差し込む

- `internal/cmd/serve/serve.go` の `app.OnServe()` 付近、もしくは専用の`internal/hooks`パッケージに、`images`コレクションへのレコード作成をフックする処理を追加。
- 流れ: アップロードされたファイル → 一旦バイト列にデコード → resize/encode → 圧縮後のバイト列で`record`の`image`フィールドを差し替えてから保存。
- ここはPocketBase v0.39の file field の差し替えAPI（`filesystem.NewFileFromBytes`相当）の挙動を実装しながら確認していく必要があります。プランの中で一番「初めて触る」部分なので、最初は小さく「ファイルを受け取ってログに出すだけ」のフックを書いて動作を確認 → 圧縮処理を組み込む、の2段階で進めるのが安全そうです。
- `filename` / `filesize` フィールドは圧縮後のサイズ・拡張子で埋める。

## Phase 2: Dockerfile / ビルド周りの修正

- picmd2の`Dockerfile`を`CGO_ENABLED=1`に戻し、`gcc musl-dev`をGoビルドステージに追加。
- ローカル開発（NixOS）でもcgoが通る状態か確認（`nix develop`のシェルにgcc等が入っているか）。

## Phase 3: フロントエンド（Solid.js）のアップロード画面

- 旧`upload.html`/`upload.js`/`upload.css`のUX（ペースト・ドラッグ&ドロップ・プログレスバー・結果カード・Markdownコピー）を、`frontend/src/routes/Upload.jsx`のようなコンポーネントとして移植。
- アップロード自体は素の`fetch`/`XMLHttpRequest`ではなく、PocketBase JS SDKの`pb.collection('images').create(formData)`を使う形に統一。
- `main.jsx`のルーターに追加、`NavBar`にリンクを足す。

## Phase 4（確定版）: URL / access設計

**画像配信URL**
- 独自の`/i/<uuid>.ext`のようなカスタムルートは作らず、PocketBase標準の `/api/files/{collection}/{recordId}/{filename}` をそのまま使う。シンプルさ優先で、これ以上薄いレイヤーは増やさない。
- `images`コレクションの`image`フィールドは現状 `"protected": false` のままでOK（トークン不要で直リンク可能）。

**コレクションルール**
- `viewRule`: 現状マイグレーションでは`null`（＝superuser限定）になっているので、`""`（空文字＝誰でも可）に変更。これで「URLを知っていれば誰でも見れる」を実現。
- `listRule`: `null`のままにしておく。一覧APIまで公開すると全画像のURLが列挙できてしまうので、ここは締めたまま。
- `createRule`: 認証必須（superuserのみ）にする。アップロード自体はSPA経由＝ログイン後にしか叩けないので二重の防御になる。

**アクセスカウント（`access`フィールド）**
- ファイル配信はPocketBaseの内部ハンドラが処理するので、カウントを挟むには「レコードのCRUDフック」ではなく、ファイル配信自体をフックする仕組みが必要。PocketBase v0.39系にある `OnFileDownloadRequest`（名称はバージョンで多少変わるので実装時に確認）でリクエストを拾い、非同期・失敗しても配信は止めない形で`access`をインクリメントする方針。
- ここは「重要でない付随機能がクリティカルパスを壊さない」ことを優先し、カウント更新に失敗してもログに残すだけで200を返す、くらいのシンプルな実装で十分（アクセス数は参考値であり厳密なカウントである必要はない）。

**アップロード画面のアクセス制御**
- 現状すでに`main.jsx`の`AuthGate`がアプリ全体をログイン必須にしているので、新しく作る`Upload`ルートもそこに含めるだけで自然と「未ログインでは見れない」が満たされる。追加の認可ロジックは不要。


