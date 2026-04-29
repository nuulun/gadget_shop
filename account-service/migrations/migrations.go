package migrations

import (
	"context"
	"database/sql"
	"fmt"
)

func Run(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, schema); err != nil {
		return fmt.Errorf("account migrations: %w", err)
	}
	return nil
}

const schema = `
CREATE TABLE IF NOT EXISTS users (
    id          BIGSERIAL PRIMARY KEY,
    login       TEXT NOT NULL UNIQUE,
    email       TEXT NOT NULL UNIQUE,
    phone       TEXT NOT NULL DEFAULT '',
    first_name  TEXT NOT NULL DEFAULT '',
    last_name   TEXT NOT NULL DEFAULT '',
    middle_name TEXT NOT NULL DEFAULT '',
    age         INT  NOT NULL DEFAULT 0,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);
`
