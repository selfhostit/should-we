package migrations_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/tests"
	_ "github.com/selfhostit/should-we/migrations"
)

func TestNightlyBackups(t *testing.T) {
	a, err := tests.NewTestApp()
	if err != nil {
		t.Fatal(err)
	}
	defer a.Cleanup()
	b := a.Settings().Backups
	if b.Cron != "0 3 * * *" || b.CronMaxKeep != 7 {
		t.Errorf("backups = %q, keep %d; want nightly at 3, keep 7", b.Cron, b.CronMaxKeep)
	}
}
