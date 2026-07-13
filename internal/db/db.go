// Package db wraps the embedded PocketBase datastore that persists card performance and
// review history. All persistence is issued through this package; no other package
// touches the datastore directly.
package db

import (
	"github.com/asano69/myapp/internal/errs"
	_ "github.com/asano69/myapp/migrations"

	"github.com/pocketbase/pocketbase"

	"os"
)

type Database struct{ app *pocketbase.PocketBase }

// OpenScratch creates a Database backed by a fresh, disposable PocketBase
// instance in its own temporary directory. Each call returns an
// independent, empty database with no effect on any other Database.
// PocketBase always needs a directory on disk, so this is myapp'
// equivalent of SQLite's ":memory:" mode.
func OpenScratch() (*Database, error) {
	dir, err := os.MkdirTemp("", "myapp-pocketbase-*")
	if err != nil {
		return nil, errs.Newf("create temporary PocketBase data directory: %v", err)
	}
	app := pocketbase.NewWithConfig(pocketbase.Config{DefaultDataDir: dir, HideStartBanner: true})
	if err := app.Bootstrap(); err != nil {
		return nil, errs.Newf("bootstrap PocketBase: %v", err)
	}
	return newDatabase(app)
}

// New wraps an already-bootstrapped PocketBase app and ensures the
// myapp schema exists in it. app is expected to be the single instance
// shared by the whole CLI (see cmd/myapp/main.go); its data directory is
// controlled by PocketBase's standard "--dir" flag, not by myapp itself.
func New(app *pocketbase.PocketBase) (*Database, error) {
	return newDatabase(app)
}

// newDatabase wraps app in a Database and applies any pending app-level
// schema migrations (see internal/migrations). System migrations
// (_collections, _params, ...) already ran inside app.Bootstrap(), so only
// the user-defined AppMigrations need to be applied here. Calling this on
// every startup (including in tests, via OpenScratch) is safe and
// idempotent — RunAppMigrations skips migrations already recorded in the
// _migrations table.
func newDatabase(app *pocketbase.PocketBase) (*Database, error) {
	if err := app.RunAppMigrations(); err != nil {
		return nil, errs.Newf("run migrations: %v", err)
	}

	db := &Database{app: app}

	return db, nil
}
