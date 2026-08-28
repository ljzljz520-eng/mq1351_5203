package service

import (
	"context"
	"fmt"
	"strings"

	"qin-culture-site/internal/domain"
	"qin-culture-site/internal/store"
)

type CultureReport struct {
	Stats      domain.CultureStats
	Labels     []string
	Timeline   []domain.TimelineEntry
	Recent     []domain.ExperienceSubmission
	SearchTerm string
	SearchHits int
}

func (s *Service) Report(ctx context.Context, query string) (CultureReport, error) {
	if s.catalog == nil {
		return CultureReport{}, fmt.Errorf("catalog is not configured")
	}
	result := s.catalog.Search(query)
	recent := []domain.ExperienceSubmission{}
	if s.store != nil {
		items, err := s.store.ListExperiences(ctx)
		if err != nil {
			return CultureReport{}, err
		}
		recent = items
	}
	return CultureReport{Stats: result.Stats, Labels: result.Stats.Labels(), Timeline: s.catalog.Timeline(), Recent: recent, SearchTerm: strings.TrimSpace(query), SearchHits: len(result.Hits)}, nil
}

func (r CultureReport) Headline() string {
	return fmt.Sprintf("收录 %d 个琴派、%d 首名曲，等待每位访客找到自己的入琴之门", r.Stats.SchoolCount, r.Stats.PieceCount)
}

func (r CultureReport) Complete() bool {
	return r.Stats.HasMinimumContent() && len(r.Timeline) > 0
}

func (r CultureReport) AuditFor(submission domain.ExperienceSubmission) store.AuditEntry {
	return store.AuditEntry{Action: "experience-submitted", Entity: "ExperienceSubmission", EntityID: fmt.Sprint(submission.ID), Detail: submission.Summary(), CreatedAt: submission.CreatedAt}
}

func (r CultureReport) Lines() []string {
	lines := []string{r.Headline(), strings.Join(r.Labels, " / ")}
	for _, entry := range r.Timeline {
		lines = append(lines, fmt.Sprintf("%s | %s | %s", entry.Era, entry.Label, entry.Place))
	}
	return lines
}
