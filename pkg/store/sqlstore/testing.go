package sqlstore

import (
	"fmt"
	"strings"
	"testing"

	_ "github.com/lib/pq"

	"github.com/opencars/alpr/pkg/config"
)

// TestDB returns special test connection and teardown function.
func TestDB(t *testing.T, conf *config.Database) (*Store, func(...string)) {
	t.Helper()

	store, err := New(conf)
	if err != nil {
		t.Fatal(err)
	}

	return store, func(tables ...string) {
		if len(tables) > 0 {
			_, err = store.db.Exec(fmt.Sprintf("TRUNCATE %s RESTART IDENTITY CASCADE", strings.Join(tables, ", ")))
			if err != nil {
				t.Fatal(err)
			}
		}

		if err := store.db.Close(); err != nil {
			t.Fatal(err)
		}
	}
}
