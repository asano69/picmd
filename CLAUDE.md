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
 

