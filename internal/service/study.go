package service

import (
	"context"
	"fmt"

	"qin-culture-site/internal/catalog"
	"qin-culture-site/internal/domain"
)

type StudyRequest struct {
	PieceID    string
	Query      string
	Difficulty string
}

type StudyModel struct {
	Bundle            domain.StudyBundle
	DifficultyOptions []string
	SelectedQuery     string
}

func (s *Service) Study(ctx context.Context, request StudyRequest) (StudyModel, error) {
	if s.catalog == nil {
		return StudyModel{}, fmt.Errorf("catalog is not configured")
	}
	filter := catalog.StudyFilter{PieceID: request.PieceID, Query: request.Query, Difficulty: request.Difficulty}
	bundle := s.catalog.FilterStudy(filter)
	if bundle.Piece.ID == "" {
		return StudyModel{}, fmt.Errorf("unknown piece %q", request.PieceID)
	}
	for _, notation := range bundle.Notations {
		if s.store != nil {
			if err := s.store.SaveNotation(ctx, notation); err != nil {
				return StudyModel{}, err
			}
		}
	}
	return StudyModel{Bundle: bundle, DifficultyOptions: s.catalog.DifficultyOptions(), SelectedQuery: domain.ValidateStudyQuery(request.Query)}, nil
}

func (s *Service) StudySummary(model StudyModel) string {
	if !model.Bundle.Complete() {
		return "资料尚未完整"
	}
	return fmt.Sprintf("%s：%d 条谱例，%d 条礼仪，%d 个传承故事", model.Bundle.Piece.Title, len(model.Bundle.Notations), len(model.Bundle.Courtesy), len(model.Bundle.Stories))
}
