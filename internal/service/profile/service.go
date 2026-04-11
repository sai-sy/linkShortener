package profile

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

func (s *Service) GetByID(ctx context.Context, profileID int64) (db.GetProfileByIDRow, error) {
	var profile db.GetProfileByIDRow
	return profile, s.rls.WithProfile(ctx, profileID, func(q *db.Queries) error {
		row, err := q.GetProfileByID(ctx, profileID)
		if err != nil {
			return err
		}
		profile = row
		return nil
	})
}
