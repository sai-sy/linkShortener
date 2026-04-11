package workspace

import (
	"context"

	"github.com/sai-sy/linkShortener/internal/db"
	"github.com/sai-sy/linkShortener/internal/service/rls"
)

type Service struct {
	rls *rls.Service
}

func New(rlsSvc *rls.Service) *Service {
	return &Service{rls: rlsSvc}
}

func (s *Service) List(ctx context.Context, profileID int64, page int) ([]db.Workspace, error) {
	limit := int32(25)
	offset := int32((page - 1) * 25)
	var workspaces []db.Workspace
	return workspaces, s.rls.WithProfile(ctx, profileID, func(q *db.Queries) error {
		rows, err := q.GetWorkspacesPage(ctx, db.GetWorkspacesPageParams{
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			return err
		}
		workspaces = rows
		return nil
	})
}
