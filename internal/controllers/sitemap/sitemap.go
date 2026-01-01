package sitemap

import (
	"bytes"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sabloger/sitemap-generator/smg"
)

func Serve(ctx context.Context, app *fiber.App, host string) {
	app.Get("/sitemap.xml", func(c *fiber.Ctx) error {
		now := time.Now().UTC()
		sm := smg.NewSitemap(true)
		sm.SetHostname(host)
		sm.SetLastMod(&now)
		sm.SetCompress(false)
		err := sm.Add(&smg.SitemapLoc{
			Loc:        "/",
			LastMod:    &now,
			ChangeFreq: smg.Daily,
			Priority:   0.8,
		})
		if err != nil {
			return err
		}
		err = sm.Add(&smg.SitemapLoc{
			Loc:        "/login",
			LastMod:    &now,
			ChangeFreq: smg.Weekly,
			Priority:   0.6,
		})
		if err != nil {
			return err
		}
		sm.Finalize()

		var buf bytes.Buffer

		_, err = sm.WriteTo(&buf)
		if err != nil {
			return err
		}
		c.Set("Content-Type", "application/xml")
		return c.SendString(buf.String())
	})
}
