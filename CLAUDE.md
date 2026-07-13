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

- 複数画像のアップロード機能
- 画像配信URLの再検討
    - この画像アップローダから別のアップローダに自然に移行できるように、画像URLはPocketBaseに依存しない表現のほうが安心できる。
    - もとのpicmdのように、カスタムルートで画像表示できるようにしたい。
- 画像にうめこまれたメタデータから個人情報を削除する方法の検討

# 画像配信URLの再検討

UUIDv7をレコード作成時にフックで払い出して、それを鍵に配信URLを組み立てるというのは、PocketBase依存を薄くするという目的にちょうど合っています。既存の `RegisterImageCompression` と同じ形のフックが書けるので、実装コストも小さいです。

以下、具体的な計画です。

## 1. UUID付与フック

`internal/hooks/images.go` に圧縮フックとは別関数として追加するのがシンプルだと思います（関心事を分ける、CLAUDE.mdの「Keep functions small and focused」にも合う）。

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

## 2. 一意性の担保

`uuid` フィールドに一意インデックスを貼っておくべきです。CLAUDE.mdのルール通り、これはコードでマイグレーションを書くのではなく、PocketBaseの管理UIからインデックスを追加してください（自動でスナップショットmigrationが生成されます）。

## 3. 配信ルート `/i/{uuid}`


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

---

まとめると、変更点は「フックを1つ追加」「ルートを1つ追加」「フロント1行修正」「管理UIでインデックス追加」の4点で、既存の設計（`internal/hooks`, `internal/cmd/serve`の構成）にそのまま乗る形なので、シンプルさは保てると思います。この方向で進めて良さそうですか?
 

