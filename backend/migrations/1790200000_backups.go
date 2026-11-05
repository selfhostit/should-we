package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// Nightly backups into pb_data/backups, keeping the last 7, like the other selfhostit
// apps. Change them in the admin UI under Settings > Backups; installs that set a
// schedule already are left alone.
func init() {
	m.Register(func(app core.App) error {
		settings := app.Settings()
		if settings.Backups.Cron != "" {
			return nil
		}
		settings.Backups.Cron = "0 3 * * *"
		settings.Backups.CronMaxKeep = 7
		return app.Save(settings)
	}, nil)
}
