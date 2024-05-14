package links

import (
	"context"
	"time"

	"github.com/iamnoturkkitty/shortener/internal/domain/links"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepo(cs string) (*PostgresRepo, error) {
	pool, err := pgxpool.New(context.Background(), cs)
	if err != nil {
		return nil, err
	}

	if _, err := pool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS links (
			"id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
			"short_url" VARCHAR(250) NOT NULL UNIQUE,
			"created" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			"original_url" TEXT)`); err != nil {
		return nil, err
	}

	return &PostgresRepo{
		db: pool,
	}, nil
}

func (r *PostgresRepo) SaveLink(ctx context.Context, l links.Link) error {
	rs, err := r.db.Exec(ctx, "INSERT INTO links  (short_url, original_url) VALUES ($1, $2) ON CONFLICT DO NOTHING", l.Hash(), l.URL())
	if err != nil {
		return err
	}

	ra := rs.RowsAffected()

	if ra == 0 {
		return links.ErrLinkDuplicate
	}

	return nil
}

func (r *PostgresRepo) SaveLinkBatch(ctx context.Context, ls []links.Link) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	for _, l := range ls {
		_, err := r.db.Exec(context.Background(), "INSERT INTO links  (short_url, original_url) VALUES ($1, $2);", l.Hash(), l.URL())
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *PostgresRepo) GetLink(ctx context.Context, hash string) (*links.Link, error) {
	row := r.db.QueryRow(ctx, "SELECT id, short_url, original_url FROM links WHERE short_url=$1", hash)

	var l links.StoredLink

	if err := row.Scan(&l.ID, &l.Hash, &l.URL); err != nil {
		return nil, err
	}

	return links.NewLink(l.ID, l.URL, l.Hash)
}

func (r *PostgresRepo) Test() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := r.db.Ping(ctx); err != nil {
		return err
	}

	return nil
}
