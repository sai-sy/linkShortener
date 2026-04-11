package rls

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/sai-sy/linkShortener/internal/db"
)

type Service struct {
	conn *pgx.Conn
	db   *db.Queries
}

func New(db *db.Queries, conn *pgx.Conn) *Service {
	return &Service{conn: conn, db: db}
}

func (s *Service) WithProfile(ctx context.Context, profileID int64, fn func(q *db.Queries) error) error {
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("rls begin: %w", err)
	}

	if _, err := tx.Exec(ctx, "SELECT set_config('app.profile_id', $1, true)", strconv.FormatInt(profileID, 10)); err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("rls set profile: %w", err)
	}

	queries := s.db.WithTx(tx)
	if err := fn(queries); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("rls commit: %w", err)
	}

	return nil
}
