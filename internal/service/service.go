package service

import (
	"context"
	"fmt"

	"qin-culture-site/internal/catalog"
	"qin-culture-site/internal/domain"
	"qin-culture-site/internal/store"
)

type Service struct {
	catalog *catalog.Catalog
	store   *store.Store
}

type HomeModel struct {
	Schools         []domain.QinSchool
	Pieces          []domain.QinPiece
	FeaturedStories []domain.HeritageStory
	ImageWall       []domain.Card
	Courtesy        []domain.Courtesy
}

func New(c *catalog.Catalog, s *store.Store) *Service {
	return &Service{catalog: c, store: s}
}

func (s *Service) Browse(ctx context.Context) (HomeModel, error) {
	if s.catalog == nil {
		return HomeModel{}, fmt.Errorf("catalog is not configured")
	}
	model := HomeModel{Schools: s.catalog.Schools(), Pieces: s.catalog.SortedPieces(), FeaturedStories: s.catalog.FeaturedStories(), ImageWall: s.catalog.ImageWall(), Courtesy: s.catalog.Courtesy()}
	if s.store != nil {
		for _, school := range model.Schools {
			if err := s.store.SaveQinSchool(ctx, school); err != nil {
				return HomeModel{}, err
			}
		}
		for _, piece := range model.Pieces {
			if err := s.store.SaveQinPiece(ctx, piece); err != nil {
				return HomeModel{}, err
			}
		}
	}
	return model, nil
}

func (s *Service) Catalog() *catalog.Catalog {
	return s.catalog
}

func (s *Service) Store() *store.Store {
	return s.store
}
