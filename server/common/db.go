package common

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type DBLogger interface {
	boil.ContextExecutor
	sqlx.Queryer
}

type QueryLogger struct {
	Db     sqlx.DB
	Logger *log.Logger
}

var _ DBLogger = QueryLogger{}

func NewDbLogger(db *sqlx.DB) DBLogger {
	return QueryLogger{
		Db:     *db,
		Logger: log.New(os.Stdout, "sqlx-logger: ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix),
	}
}

func dbLogger(logger *log.Logger, query string, args ...interface{}) {
	logger.Printf("query: '%s'", query)
	logger.Printf(" args: %+v", args...)
}

func (q QueryLogger) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	dbLogger(q.Logger, query, args)
	return q.Db.Queryx(query, args...)
}

func (q QueryLogger) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	dbLogger(q.Logger, query, args)
	return q.Db.QueryRowx(query, args...)
}

func (q QueryLogger) Exec(query string, args ...interface{}) (sql.Result, error) {
	dbLogger(q.Logger, query, args)
	return q.Db.Exec(query, args...)
}

func (q QueryLogger) Query(query string, args ...interface{}) (*sql.Rows, error) {
	dbLogger(q.Logger, query, args)
	return q.Db.Query(query, args...)
}

func (q QueryLogger) QueryRow(query string, args ...interface{}) *sql.Row {
	dbLogger(q.Logger, query, args)
	return q.Db.QueryRow(query, args...)
}

func (q QueryLogger) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	dbLogger(q.Logger, query, args)
	return q.Db.ExecContext(ctx, query, args...)
}

func (q QueryLogger) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	dbLogger(q.Logger, query, args)
	return q.Db.QueryContext(ctx, query, args...)
}

func (q QueryLogger) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	dbLogger(q.Logger, query, args)
	return q.Db.QueryRowContext(ctx, query, args...)
}

// ------- DB UTILS

func DbUrl(dbName, user, password, host, port string, ssl bool) string {
	sslStr := "enable"
	if !ssl {
		sslStr = "disable"
	}

	// postgres url fails if no password is set but colon is added
	if password != "" {
		password = fmt.Sprintf(":%s", password)
	}

	return fmt.Sprintf(
		"postgres://%s%s@%s:%s/%s?sslmode=%s",
		user, password,
		host, port,
		dbName,
		sslStr,
	)
}

// Write queries on multiple lines in the code but wrap them in a single line for better logging
func SqlFmt(qs string) string {
	return strings.TrimSpace(strings.Replace(qs, "\n", " ", -1))
}

func RunMigrations(migrationsDir, dbUrl string) {
	m, err := migrate.New(fmt.Sprintf("file://%s", migrationsDir), dbUrl)
	if err != nil {
		log.Fatalf("could not migrate on '%s': %+v", dbUrl, err)
	}
	defer m.Close()

	// XXX: https://github.com/golang-migrate/migrate/issues/179
	if err := m.Up(); err != nil && err.Error() != "no change" {
		log.Fatal(err)
	}
	log.Print("migrations checked")
}
