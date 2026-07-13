- 旧実装はXHRで正確な進捗%を出していましたが、PocketBase SDKは内部で`fetch`を使っており、素の`fetch`にはアップロード進捗イベントがありません。シンプルさ優先の方針に沿って、進捗バーは「アップロード中…」の単純な表示に簡略化します
- pb.files.getURL(record, record.image) はPocketBase標準の /api/files/{collection}/{recordId}/{filename} を組み立てます。

## URL / access設計
**画像配信URL**
- PocketBase標準の `/api/files/{collection}/{recordId}/{filename}` をそのまま使っている。
- `images`コレクションの`image`フィールドは`"protected": false` のままでOK（トークン不要で直リンク可能）。

**アクセスカウント（`access`フィールド）**
- ファイル配信はPocketBaseの内部ハンドラが処理するので、カウントを挟むには「レコードのCRUDフック」ではなく、ファイル配信自体をフックする仕組みが必要。PocketBase v0.39系にある `OnFileDownloadRequest`でリクエストを拾い、非同期・失敗しても配信は止めない形で`access`をインクリメントする方針。
- ここは「重要でない付随機能がクリティカルパスを壊さない」ことを優先し、カウント更新に失敗してもログに残すだけで200を返す、くらいのシンプルな実装で十分（アクセス数は参考値であり厳密なカウントである必要はない）。

**アップロード画面のアクセス制御**
- 現状すでに`main.jsx`の`AuthGate`がアプリ全体をログイン必須にしている。`Upload`ルートもそこに含めるだけで自然と「未ログインでは見れない」が満たされる。
