package migrations

import (
	"context"
	"database/sql"
	"fmt"
)

func Run(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, createProductsTable); err != nil {
		return fmt.Errorf("product migrations: %w", err)
	}
	return nil
}

const createProductsTable = `
CREATE TABLE IF NOT EXISTS products (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT          NOT NULL,
    brand       TEXT          NOT NULL DEFAULT '',
    description TEXT          NOT NULL DEFAULT '',
    price       NUMERIC(12,2) NOT NULL DEFAULT 0,
    stock       INT           NOT NULL DEFAULT 0,
    category    TEXT          NOT NULL DEFAULT '',
    image_url   TEXT          NOT NULL DEFAULT '',
    is_active   BOOLEAN       NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMP     NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP     NOT NULL DEFAULT NOW()
);
ALTER TABLE products ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT TRUE;
CREATE INDEX IF NOT EXISTS idx_products_category  ON products(category)  WHERE is_active = TRUE;
CREATE INDEX IF NOT EXISTS idx_products_is_active ON products(is_active);
`
