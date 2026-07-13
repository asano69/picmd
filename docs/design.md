- 旧実装はXHRで正確な進捗%を出していましたが、PocketBase SDKは内部で`fetch`を使っており、素の`fetch`にはアップロード進捗イベントがありません。シンプルさ優先の方針に沿って、進捗バーは「アップロード中…」の単純な表示に簡略化します


## access設計

**アクセスカウント（`views`フィールド）**
- ファイル配信はPocketBaseの内部ハンドラが処理するので、カウントを挟むには「レコードのCRUDフック」ではなく、ファイル配信自体をフックする仕組みが必要。PocketBase v0.39系にある `OnFileDownloadRequest`でリクエストを拾い、非同期・失敗しても配信は止めない形で`views`をインクリメントする方針。
- ここは「重要でない付随機能がクリティカルパスを壊さない」ことを優先し、カウント更新に失敗してもログに残すだけで200を返す、くらいのシンプルな実装で十分（アクセス数は参考値であり厳密なカウントである必要はない）。

**アップロード画面のアクセス制御**
- 現状すでに`main.jsx`の`AuthGate`がアプリ全体をログイン必須にしている。`Upload`ルートもそこに含めるだけで自然と「未ログインでは見れない」が満たされる。

## 画像配信URLの設計

- UUIDv7をレコード作成時にフックで払い出して、それを鍵に配信URLを組み立てるというのは、PocketBase依存を薄くするという目的にちょうど合っています。
- 既存の `RegisterImageCompression` と同じ形のフックが書けるので、実装コストも小さいです。


### 1. UUID付与フック

`internal/hooks/images.go` に圧縮フックとは別関数として追加するのがシンプルだと思います

```go
// RegisterUUIDAssignment stamps each new "images" record with a UUIDv7,
// used to build a stable public URL independent of PocketBase's own
// collection/record-id scheme (see the "/i/{uuid}" route in serve.go).
func RegisterUUIDAssignment(app core.App) {
	app.OnRecordCreateRequest("images").BindFunc(func(e *core.RecordRequestEvent) error {
		id, err := uuid.NewV7()
		if err != nil {
			return fmt.Errorf("generate uuid: %w", err)
		}
		e.Record.Set("uuid", id.String())
		return e.Next()
	})
}
```

同じイベントに複数の `BindFunc` を登録できるので、圧縮フックとは完全に独立させて `serve.go` 側で両方 `Register...` すればOKです。


### 2. 配信ルート `/i/{uuid}`

リダイレクト方式を使う

```go
e.Router.GET("/i/{uuid}", func(re *core.RequestEvent) error {
	record, err := app.FindFirstRecordByFilter(
		"images", "uuid = {:uuid}",
		dbx.Params{"uuid": re.Request.PathValue("uuid")},
	)
	if err != nil {
		return apis.NewNotFoundError("image not found", err)
	}
	// Redirect to PocketBase's native file URL. The published "/i/{uuid}"
	// link stays stable even if the storage backend changes later.
	target := fmt.Sprintf("/api/files/images/%s/%s", record.Id, record.GetString("filename"))
	http.Redirect(re.Response, re.Request, target, http.StatusFound)
	return nil
})
```

- 既存の `OnFileDownloadRequest`（閲覧数カウント）をそのまま使い回せる。
- レンジリクエスト、Content-Type、ETagなど全部PocketBaseの既存実装に乗っかれるので実装量が最小。
- デメリット: ブラウザの最終URLは結局PocketBaseの `/api/files/...` になる（＝完全に隠蔽はできない）。でも「Markdownに貼るURL」自体は `/i/{uuid}` のまま不変なので、目的（＝別のアップローダに移行してもリンク切れしない）は満たせます。



