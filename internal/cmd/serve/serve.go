// Package serve implements the "serve" command, which runs a single HTTP server
// that hosts the index page and all drill sessions defined in the config file.
package serve

import (
	"fmt"
	"github.com/asano69/myapp/internal/assets"
	"github.com/asano69/myapp/internal/config"
	"io/fs"
	"net/http"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"

	"github.com/sirupsen/logrus"
)

// Run opens the database and collection once, registers all drill routes, then
// starts listening. The database and collection are shared across all sessions.
func Run(app *pocketbase.PocketBase, cfg *config.Config) error {

	// assetsFS exposes just the "assets/" subdirectory that Vite's default
	// (unprefixed) base writes hashed JS/CSS bundles into, so they're served
	// at the conventional /assets/... URL instead of /static/assets/....
	assetsFS, err := fs.Sub(assets.FS, "assets")
	if err != nil {
		return fmt.Errorf("sub assets fs: %w", err)
	}
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		// GET /api/sessions reloads the collection from disk on every request
		// so decks/cards added or removed since startup are reflected.
		e.Router.GET("/assets/{path...}", apis.Static(assetsFS, false))
		// Solid Router decides which screen to render client-side, so both
		// /drill and / serve the same static shell. This shell is left
		// unauthenticated on purpose: it's an empty HTML/JS bundle with no
		// data in it. Every route that actually returns collection data is
		// guarded above with RequireSuperuserAuth, so an unauthenticated
		// visitor only ever sees the login screen the SPA renders client-side.
		serveShell := func(re *core.RequestEvent) error {
			re.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
			http.ServeFileFS(re.Response, re.Request, assets.FS, "index.html")
			return nil
		}

		e.Router.GET("/", serveShell)

		// Vite's public/ directory (favicon.svg etc.) is copied to the root
		// of the build output, so it's served directly rather than under
		// /assets/.
		e.Router.GET("/favicon.svg", func(re *core.RequestEvent) error {
			re.Response.Header().Set("Content-Type", "image/svg+xml")
			http.ServeFileFS(re.Response, re.Request, assets.FS, "favicon.svg")
			return nil
		})

		return e.Next()
	})

	logrus.WithField("addr", addr).Info("listening")
	return apis.Serve(app, apis.ServeConfig{
		HttpAddr:        addr,
		ShowStartBanner: false,
	})
}
