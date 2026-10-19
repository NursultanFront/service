// Package migrate applies the database migrations and the seed data.
package migrate

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"math"

	dbfiles "github.com/ardanlabs/service/db"
	sqldb "github.com/ardanlabs/service/internal/sys/db"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
)

// NOTE(AI): darwin was replaced by golang-migrate in migration step (d); the SQL
// lives in db/migrations and db/fixtures at the repository root (package db).

const seedFile = "fixtures/seed_test_data.sql"

// Migrate attempts to bring the database up to date with the migrations in
// db/migrations.
func Migrate(ctx context.Context, db *sqlx.DB) (err error) {
	if err := sqldb.StatusCheck(ctx, db); err != nil {
		return fmt.Errorf("status check database: %w", err)
	}

	src, err := iofs.New(dbfiles.Migrations, "migrations")
	if err != nil {
		return fmt.Errorf("open migrations: %w", err)
	}

	// A dedicated connection: closing the migrator releases only this
	// connection, never the caller's pool.
	conn, err := db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("acquire connection: %w", err)
	}

	drv, err := postgres.WithConnection(ctx, conn, &postgres.Config{})
	if err != nil {
		return errors.Join(fmt.Errorf("construct migrate driver: %w", err), conn.Close())
	}

	m, err := migrate.NewWithInstance("iofs", src, "postgres", drv)
	if err != nil {
		return errors.Join(fmt.Errorf("construct migrate: %w", err), drv.Close())
	}

	defer func() {
		srcErr, drvErr := m.Close()
		if err == nil {
			err = errors.Join(srcErr, drvErr)
		}
	}()

	if err := baselineDarwin(ctx, db, m); err != nil {
		return fmt.Errorf("darwin baseline: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate up: %w", err)
	}

	return nil
}

// baselineDarwin runs once for a database that was migrated by the previous
// darwin-based tool: when golang-migrate has no version recorded yet and the
// darwin_migrations table exists, the migrations darwin already applied are
// marked as applied (darwin 1.01 ... 1.05 == files 000001 ... 000005).
func baselineDarwin(ctx context.Context, db *sqlx.DB, m *migrate.Migrate) error {
	_, _, err := m.Version()
	switch {
	case err == nil:
		return nil

	case !errors.Is(err, migrate.ErrNilVersion):
		return fmt.Errorf("read version: %w", err)
	}

	var hasDarwin bool
	if err := db.QueryRowContext(ctx, "SELECT to_regclass('public.darwin_migrations') IS NOT NULL").Scan(&hasDarwin); err != nil {
		return fmt.Errorf("check darwin table: %w", err)
	}

	if !hasDarwin {
		return nil
	}

	var darwinVersion float64
	if err := db.QueryRowContext(ctx, "SELECT COALESCE(MAX(version), 0) FROM darwin_migrations").Scan(&darwinVersion); err != nil {
		return fmt.Errorf("read darwin version: %w", err)
	}

	version := int(math.Round((darwinVersion - 1) * 100))
	if version <= 0 {
		return nil
	}

	if err := m.Force(version); err != nil {
		return fmt.Errorf("force version %d: %w", version, err)
	}

	return nil
}

// Seed runs the idempotent fixtures in db/fixtures against the database.
func Seed(ctx context.Context, db *sqlx.DB) (err error) {
	if err := sqldb.StatusCheck(ctx, db); err != nil {
		return fmt.Errorf("status check database: %w", err)
	}

	seedDoc, err := fs.ReadFile(dbfiles.Fixtures, seedFile)
	if err != nil {
		return fmt.Errorf("read seed: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if errTx := tx.Rollback(); errTx != nil {
			if errors.Is(errTx, sql.ErrTxDone) {
				return
			}

			err = fmt.Errorf("rollback: %w", errTx)
			return
		}
	}()

	if _, err := tx.Exec(string(seedDoc)); err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}
