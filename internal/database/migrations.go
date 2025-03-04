package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigrations(pool *pgxpool.Pool) error {
	ctx := context.Background()

	migrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS links (
			id SERIAL PRIMARY KEY,
			short_code VARCHAR(20) UNIQUE NOT NULL,
			original_url TEXT NOT NULL,
			title VARCHAR(200),
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			is_active BOOLEAN DEFAULT true,
			expires_at TIMESTAMPTZ,
			max_clicks INTEGER,
			password_hash VARCHAR(255),
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_links_short_code ON links(short_code)`,
		`CREATE INDEX IF NOT EXISTS idx_links_user_id ON links(user_id)`,
		`CREATE TABLE IF NOT EXISTS clicks (
			id SERIAL PRIMARY KEY,
			link_id INTEGER REFERENCES links(id) ON DELETE CASCADE NOT NULL,
			ip_address VARCHAR(45),
			user_agent TEXT,
			referer TEXT,
			country VARCHAR(3),
			city VARCHAR(100),
			device VARCHAR(20),
			browser VARCHAR(50),
			os VARCHAR(50),
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_clicks_link_id ON clicks(link_id)`,
		`CREATE INDEX IF NOT EXISTS idx_clicks_created_at ON clicks(created_at)`,
		`CREATE TABLE IF NOT EXISTS tags (
			id SERIAL PRIMARY KEY,
			name VARCHAR(50) NOT NULL,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			UNIQUE(name, user_id)
		)`,
		`CREATE TABLE IF NOT EXISTS link_tags (
			link_id INTEGER REFERENCES links(id) ON DELETE CASCADE,
			tag_id INTEGER REFERENCES tags(id) ON DELETE CASCADE,
			PRIMARY KEY (link_id, tag_id)
		)`,
	}

	for i, m := range migrations {
		if _, err := pool.Exec(ctx, m); err != nil {
			return fmt.Errorf("migration %d: %w", i+1, err)
		}
	}

	return nil
}
