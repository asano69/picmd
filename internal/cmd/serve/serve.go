package serve

import (
	"fmt"
	"github.com/asano69/picmd2/internal/assets"
	"github.com/asano69/picmd2/internal/config"
	"github.com/asano69/picmd2/internal/hooks"
	"io/fs"
	"net/http"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"

	"github.com/sirupsen/logrus"
)

func Run(app *pocketbase.PocketBase, cfg *config.Config) error {
	hooks.RegisterImageCompression(app)
	hooks.RegisterUUIDAssignment(app)
	hooks.RegisterViewCounter(app)

	assetsFS, err := fs.Sub(assets.FS, "assets")
	if err != nil {
		return fmt.Errorf("sub assets fs: %w", err)
	}
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		e.Router.GET("/assets/{path...}", apis.Static(assetsFS, false))

		serveShell := func(re *core.RequestEvent) error {
			re.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
			http.ServeFileFS(re.Response, re.Request, assets.FS, "index.html")
			return nil
		}

		e.Router.GET("/", serveShell)

		e.Router.GET("/favicon.svg", func(re *core.RequestEvent) error {
			re.Response.Header().Set("Content-Type", "image/svg+xml")
			http.ServeFileFS(re.Response, re.Request, assets.FS, "favicon.svg")
			return nil
		})

		// /i/{uuid} is the stable public link handed out for pasting into
		// Markdown. It redirects to PocketBase's native file URL, so this
		// is the only place that needs to change if storage ever moves
		// away from PocketBase.
		e.Router.GET("/i/{uuid}", func(re *core.RequestEvent) error {
			record, err := app.FindFirstRecordByFilter(
				"images", "uuid = {:uuid}",
				dbx.Params{"uuid": re.Request.PathValue("uuid")},
			)
			if err != nil {
				return apis.NewNotFoundError("image not found", err)
			}
			target := fmt.Sprintf("/api/files/images/%s/%s", record.Id, record.GetString("filename"))
			http.Redirect(re.Response, re.Request, target, http.StatusFound)
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
