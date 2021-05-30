package trade

import (
	"context"
	"time"

	"github.com/pohsi/pktrade/pkg/log"
)

type Service interface {
	Get(ctx context.Context, id string) (Album, error)
	Query(ctx context.Context, offset, limit int) ([]Album, error)
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, input CreateAlbumRequest) (Album, error)
	Update(ctx context.Context, id string, input UpdateAlbumRequest) (Album, error)
	Delete(ctx context.Context, id string) (Album, error)
}

type service struct {
	repo   Repository
	logger log.Logger
}

func CreateService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

func (s service) Get(ctx context.Context, id string) (Album, error) {
	album, err := s.repo.Get(ctx, id)
	if err != nil {
		return Album{}, err
	}
	return Album{album}, nil
}

func (s service) Buy(ctx context.Context, req CreateAlbumRequest) (Album, error) {
	if err := req.Validate(); err != nil {
		return Album{}, err
	}
	id := entity.GenerateID()
	now := time.Now()
	err := s.repo.Create(ctx, entity.Album{
		ID:        id,
		Name:      req.Name,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return Album{}, err
	}
	return s.Get(ctx, id)
}

type Trade struct {
	entity.Album
}

type CreateAlbumRequest struct {
	Name string `json:"name"`
}
