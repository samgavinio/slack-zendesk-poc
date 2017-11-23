package operation

import (
	"errors"
	"flag"
	"log"

	"github.com/go-gormigrate/gormigrate"

	"github.com/zendesk/slack-poc/config"
)

func Migrate(
	cfg config.Config,
	versions []*gormigrate.Migration,
	options *gormigrate.Options,
) {
	var (
		version string
		migrateUp bool
		migrateDown bool
	)

	flag.StringVar(&version, "version", "", "Version to migrate.")
	flag.BoolVar(&migrateUp, "up", false, "Migrate version up.")
	flag.BoolVar(&migrateDown, "down", false, "Migrate version down.")

	flag.Parse()
	db := connect(cfg)

	var (
		err error
		vs []*gormigrate.Migration
		v *gormigrate.Migration
	)

	if migrateUp || migrateDown {
		v, err = getVersion(version, versions)
		if err != nil {
			log.Fatalf("Could not migrate: %v", err)
		}
		vs = append(vs, v)
	} else {
		vs = versions
	}

	m := gormigrate.New(db, options, vs)
	if migrateDown {
		if err = m.RollbackMigration(v); err != nil {
			log.Fatalf("Could not migrate: %v", err)
		}
	} else {
		if err = m.Migrate(); err != nil {
			log.Fatalf("Could not migrate: %v", err)
		}
	}

	log.Printf("Migration completed successfully.")
}

func getVersion(version string, versions []*gormigrate.Migration) (*gormigrate.Migration, error) {
	for i := range versions {
		if versions[i].ID == version {
			return versions[i], nil
		}
	}

	return nil, errors.New("cannot find version " + version)
}
