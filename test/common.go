package test

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/rafata1/feedies/config"
	"os"
	"path"
	"sync"
)

// integration ...
type integration struct {
	db   *sqlx.DB
	conf config.Config
}

var initOnce sync.Once

var globalConf config.Config
var globalDB *sqlx.DB

func NewIntegration() *integration {
	rootDir := findRootDir()

	conf := config.LoadTestConfig(rootDir)
	migrateUpForTesting(rootDir, conf.MySQL.DSN())

	db := conf.MySQL.MustConnect()

	globalConf = conf
	globalDB = db

	t := &integration{
		conf: globalConf,
		db:   globalDB,
	}
	t.truncate()
	return t
}

// Finish ...
func (t *integration) truncate() {
	query := `
	SELECT
		CONCAT('TRUNCATE ',table_name, ';') as query
	FROM
		information_schema.tables
	WHERE
		table_schema = ? AND table_name <> ?
	ORDER BY
		table_name DESC;
`
	var queries []string
	err := t.db.Select(&queries, query, t.conf.MySQL.Database, "schema_migrations")
	if err != nil {
		return
	}
	for _, q := range queries {
		t.db.MustExec(q)
	}
}

func findRootDir() string {
	workdir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	directory := workdir
	for {
		files, err := os.ReadDir(directory)
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if file.Name() == "go.mod" {
				return directory
			}
		}
		directory = path.Dir(directory)
	}
}

func migrateUpForTesting(rootDir string, dsn string) {
	sourceURL := fmt.Sprintf("file://%s", path.Join(rootDir, "migration"))
	databaseURL := fmt.Sprintf("mysql://%s", dsn)

	fmt.Println("SourceURL:", sourceURL)
	fmt.Println("DatabaseURL:", databaseURL)

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		panic(err)
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		fmt.Println("No change in migration")
		return
	}
	if err != nil {
		panic(err)
	}
}
