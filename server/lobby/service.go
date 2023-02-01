package lobby

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

var (
	ErrNotFound = fmt.Errorf("lobby not found")
)

type Service struct {
	db *gorm.DB
}

func NewService(gorm *gorm.DB) *Service {
	return &Service{db: gorm}
}

// Create creates a new lobby
func (s *Service) Create(ctx context.Context) (*Lobby, error) {
	slug, _ := uuid.NewUUID()
	l := &Lobby{Slug: slug.String()}
	err := s.db.Create(l).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create lobby: %w", err)
	}
	return l, nil
}

// Delete deletes existing lobby
func (s *Service) Delete(ctx context.Context, slug string) error {
	return nil
}

// Update updates existing lobby
func (s *Service) Get(ctx context.Context, slug string) (*Lobby, error) {
	var l Lobby
	tx := s.db.Where("slug = ?", slug).First(&l)
	if err := tx.Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrNotFound
		default:
			return nil, fmt.Errorf("failed to get a lobby: %w", err)
		}
	}
	return &l, nil
}
