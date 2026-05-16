package migrations

import (
	"context"
	"database/sql"
	"fmt"
)

func Run(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, schema); err != nil {
		return fmt.Errorf("notification migrations: %w", err)
	}
	return nil
}

const schema = `
CREATE TABLE IF NOT EXISTS notifications (
    id         BIGSERIAL PRIMARY KEY,
    recipient  TEXT NOT NULL,
    type       TEXT NOT NULL DEFAULT 'email',
    message    TEXT NOT NULL,
    status     TEXT NOT NULL DEFAULT 'sent',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
`
