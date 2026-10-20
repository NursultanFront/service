// Package db embeds the SQL migrations and local fixtures of this repository
// so the admin tool and the test harness can apply them without files on disk.
//
// NOTE(AI): created in migration step (d) (darwin -> golang-migrate).
package db

import "embed"

// Migrations holds db/migrations/*.sql in golang-migrate format
// (NNNNNN_<name>.up.sql / NNNNNN_<name>.down.sql).
//
//go:embed migrations/*.sql
var Migrations embed.FS

// Fixtures holds db/fixtures/*.sql: idempotent local seed data
// (INSERT ... ON CONFLICT DO NOTHING).
//
//go:embed fixtures/*.sql
var Fixtures embed.FS
