package routemap

import (
	"context"
	"strings"

	"github.com/sai-sy/linkShortener/internal/db"
	"github.com/sai-sy/linkShortener/internal/service/rls"
)

type Service struct {
	rls *rls.Service
}

func New(rlsSvc *rls.Service) *Service {
	return &Service{rls: rlsSvc}
}

func (s *Service) Create(ctx context.Context, profileID int64, destination string) error {
	destination = normalizeDestination(destination)
	return s.rls.WithProfile(ctx, profileID, func(q *db.Queries) error {
		workspace, err := q.GetWorkspaceByProfile(ctx, profileID)
		if err != nil {
			return err
		}
		return q.InsertRoutemap(ctx, db.InsertRoutemapParams{
			Destination: destination,
			WorkspaceID: workspace.ID,
		})
	})
}

func (s *Service) UpdateDestination(ctx context.Context, profileID int64, id int64, destination string) error {
	destination = normalizeDestination(destination)
	return s.rls.WithProfile(ctx, profileID, func(q *db.Queries) error {
		return q.UpdateRoutemapDestination(ctx, db.UpdateRoutemapDestinationParams{
			ID:          id,
			Destination: destination,
		})
	})
}

func (s *Service) List(ctx context.Context, profileID int64, page int) ([]db.GetRoutemapsPageRow, error) {
	limit := int32(25)
	offset := int32((page - 1) * 25)
	var routemaps []db.GetRoutemapsPageRow
	return routemaps, s.rls.WithProfile(ctx, profileID, func(q *db.Queries) error {
		rows, err := q.GetRoutemapsPage(ctx, db.GetRoutemapsPageParams{
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			return err
		}
		routemaps = rows
		return nil
	})
}

func normalizeDestination(value string) string {
	if value == "" {
		return value
	}
	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		return value
	}
	return "https://" + value
}
