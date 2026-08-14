package migrations

import "embed"

// FS contains all SQL migration files embedded at compile time.
// This makes the Go binary fully self-contained, eliminating the
// need to ship migration files alongside the executable.
//
//go:embed *.sql
var FS embed.FS
