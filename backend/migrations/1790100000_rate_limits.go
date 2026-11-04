package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// Rate limits, which PocketBase ships turned off: on signing in to the admin
// UI, so its password can't be guessed by trying many, and on starting a
// decision, which needs no account, so nobody can fill the disk with them.
// They count per address; behind a reverse proxy every visitor shares the
// proxy's address, so the limits then hold for everyone together. Change them
// in the admin UI under Settings > Rate limits; installs that turned them on
// already are left alone.
func init() {
	m.Register(func(app core.App) error {
		settings := app.Settings()
		if settings.RateLimits.Enabled {
			return nil
		}
		settings.RateLimits.Enabled = true
		settings.RateLimits.Rules = []core.RateLimitRule{
			{Label: "*:auth", MaxRequests: 10, Duration: 60},
			{Label: "POST /api/sw/decisions", MaxRequests: 20, Duration: 60},
		}
		return app.Save(settings)
	}, nil)
}
