package mancala

import (
	"context"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewService(gorm *gorm.DB) *Repository {
	return &Repository{db: gorm}
}

type MancalaDB struct {
	slug  string `gorm:"primarykey"`
	state []byte `gorm:"type:bytea"`
}

func (s *Repository) Save(ctx context.Context, m *Mancala) (*Mancala, error) {

	state, _ := json.Marshal(m)
	mdb := MancalaDB{
		slug:  m.LobbySlug,
		state: state,
	}
	err := s.db.Save(&mdb).Error
	if err != nil {
		return nil, fmt.Errorf("failed to save the game: %w", err)
	}
	return m, nil
}

func (s *Repository) Get(ctx context.Context, slug string) (*Mancala, error) {
	var mdb MancalaDB
	tx := s.db.Where("slug = ?", slug).First(&mdb)
	if err := tx.Error; err != nil {
		return nil, fmt.Errorf("failed to get a lobby: %w", err)

	}
	var m Mancala
	err := json.Unmarshal(mdb.state, &m)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal the game: %w", err)
	}
	return &m, nil
}
