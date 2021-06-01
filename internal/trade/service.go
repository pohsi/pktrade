package trade

import (
	"context"

	"github.com/pohsi/pktrade/pkg/log"
)

type Service interface {
}

type service struct {
	repo   Repository
	logger log.Logger
}

func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

// Count returns the number of albums.
func (s service) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}
