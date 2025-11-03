// Package service is the business logic connects the api with the internal mechanisms
package service

import (
	"github.com/topvennie/spotify_organizer/internal/database/repository"
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Uncomment this if ever needed (rollback in a service)
// func (s *Service) withRollback(ctx context.Context, fn func(context.Context) error) error {
// 	return s.repo.WithRollback(ctx, fn)
// }
