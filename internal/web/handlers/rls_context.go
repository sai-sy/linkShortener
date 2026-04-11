package handlers

import (
	"context"

	"github.com/google/uuid"
	"github.com/sai-sy/linkShortener/internal/db"
)

func (h *Handler) withProfileContext(ctx context.Context, profileID uuid.UUID) (*db.Queries, func() error, func() error, error) {
	tx, err := h.conn.Begin(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	if _, err := tx.Exec(ctx, "SET LOCAL app.profile_id = $1", profileID.String()); err != nil {
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
