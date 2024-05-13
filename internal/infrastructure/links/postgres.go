package links

import (
	"context"
	"database/sql"
	"time"

	"github.com/iamnoturkkitty/shortener/internal/domain/links"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(cs string) (*PostgresRepo, error) {
	db, err := sql.Open("pgx", cs)
	if err != nil {
		return nil, err
	}

	if _, err := db.ExecContext(context.Background(), `CREATE TABLE IF NOT EXISTS links (
			"id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
			"short_url" VARCHAR(250) NOT NULL UNIQUE,
			"created" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			"original_url" TEXT)`); err != nil {
		return nil, err
	}

	return &PostgresRepo{
		db: db,
	}, nil
}

func (r *PostgresRepo) SaveLink(l links.Link) error {
	rs, err := r.db.ExecContext(context.Background(), "INSERT INTO links  (short_url, original_url) VALUES ($1, $2) ON CONFLICT DO NOTHING", l.Hash(), l.URL())
	if err != nil {
		return err
	}

	ra, err := rs.RowsAffected()
	if err != nil {
		return err
	}

	if ra == 0 {
		return links.ErrLinkDuplicate
	}

	return nil
}

func (r *PostgresRepo) SaveLinkBatch(ls []links.Link) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	smtp, err := tx.PrepareContext(context.Background(), "INSERT INTO links  (short_url, original_url) VALUES ($1, $2);")
	if err != nil {
		return err
	}

	defer smtp.Close()

	for _, l := range ls {
		_, err := smtp.ExecContext(context.Background(), l.Hash(), l.URL())
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PostgresRepo) GetLink(hash string) (*links.Link, error) {
	smtp, err := r.db.PrepareContext(context.Background(), "SELECT id, short_url, original_url FROM links WHERE short_url=$1")
	if err != nil {
		return nil, err
	}

	defer smtp.Close()

	row := smtp.QueryRowContext(context.Background(), hash)

	var l links.StoredLink

	if err := row.Scan(&l.ID, &l.Hash, &l.URL); err != nil {
		return nil, err
	}

	return links.NewLink(l.ID, l.URL, l.Hash)
}

func (r *PostgresRepo) Test() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := r.db.PingContext(ctx); err != nil {
		return err
	}

	return nil
}
