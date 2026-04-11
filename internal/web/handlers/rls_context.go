package handlers

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/sai-sy/linkShortener/internal/db"
)

func (h *Handler) withProfileContext(ctx context.Context, profileID uuid.UUID) (*db.Queries, func() error, func() error, error) {
	tx, err := h.conn.Begin(ctx)
	if err != nil {
		log.Printf("rls context: begin failed: %v", err)
		return nil, nil, nil, err
	}

	if _, err := tx.Exec(ctx, "SELECT set_config('app.profile_id', $1, true)", profileID.String()); err != nil {
		log.Printf("rls context: set app.profile_id failed: %v", err)
		_ = tx.Rollback(ctx)
		return nil, nil, nil, err
	}

	queries := h.db.WithTx(tx)
	commit := func() error {
		return tx.Commit(ctx)
	}
	rollback := func() error {
		return tx.Rollback(ctx)
	}

	return queries, commit, rollback, nil
}
