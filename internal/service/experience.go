package service

import (
	"context"
	"fmt"

	"qin-culture-site/internal/domain"
)

func (s *Service) SubmitExperience(ctx context.Context, input domain.ExperienceSubmission) (domain.ExperienceSubmission, error) {
	if s.store == nil {
		return domain.ExperienceSubmission{}, fmt.Errorf("store is not configured")
	}
	item, err := s.store.SaveExperience(ctx, input)
	if err != nil {
		return domain.ExperienceSubmission{}, fmt.Errorf("submit experience: %w", err)
	}
	return item, nil
}

func (s *Service) Experiences(ctx context.Context) ([]domain.ExperienceSubmission, error) {
	if s.store == nil {
		return nil, fmt.Errorf("store is not configured")
	}
	return s.store.ListExperiences(ctx)
}

func (s *Service) ExperienceStatus(ctx context.Context, id int64) (string, error) {
	if s.store == nil {
		return "", fmt.Errorf("store is not configured")
	}
	item, err := s.store.FindExperience(ctx, id)
	if err != nil {
		return "", err
	}
	if item.IsPending() {
		return "等待琴友回访", nil
	}
	return item.Status, nil
}
