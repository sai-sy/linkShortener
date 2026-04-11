package handlers

import (
	"context"
	"log"
	"strconv"

	"github.com/sai-sy/linkShortener/internal/db"
)

func (h *Handler) withProfileContext(ctx context.Context, profileID int64) (*db.Queries, func() error, func() error, error) {
	tx, err := h.conn.Begin(ctx)
	if err != nil {
		log.Printf("rls context: begin failed: %v", err)
		return nil, nil, nil, err
	}

	if _, err := tx.Exec(ctx, "SELECT set_config('app.profile_id', $1, true)", strconv.FormatInt(profileID, 10)); err != nil {
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
