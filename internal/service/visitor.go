package service

import (
	"context"
	"fmt"
	"strings"

	"qin-culture-site/internal/domain"
	"qin-culture-site/internal/store"
)

type VisitorJourney struct {
	VisitorName string
	Interest    string
	Selected    []domain.QinPiece
	Completed   bool
}

func (s *Service) StartJourney(name, interest string) VisitorJourney {
	return VisitorJourney{VisitorName: strings.TrimSpace(name), Interest: strings.TrimSpace(interest), Selected: []domain.QinPiece{}}
}

func (j *VisitorJourney) AddPiece(piece domain.QinPiece) error {
	if piece.ID == "" {
		return fmt.Errorf("piece is required")
	}
	for _, selected := range j.Selected {
		if selected.ID == piece.ID {
			return fmt.Errorf("piece already selected")
		}
	}
	j.Selected = append(j.Selected, piece)
	return nil
}

func (j *VisitorJourney) RemovePiece(id string) bool {
	for index, piece := range j.Selected {
		if piece.ID == id {
			j.Selected = append(j.Selected[:index], j.Selected[index+1:]...)
			return true
		}
	}
	return false
}

func (j *VisitorJourney) Complete() error {
	if j.VisitorName == "" || j.Interest == "" {
		return fmt.Errorf("visitor name and interest are required")
	}
	if len(j.Selected) == 0 {
		return fmt.Errorf("select at least one piece")
	}
	j.Completed = true
	return nil
}

func (j VisitorJourney) Summary() string {
	names := make([]string, 0, len(j.Selected))
	for _, piece := range j.Selected {
		names = append(names, piece.Title)
	}
	return fmt.Sprintf("%s · %s · %s", j.VisitorName, j.Interest, strings.Join(names, "、"))
}

func (s *Service) SaveJourney(ctx context.Context, journey VisitorJourney) (domain.ExperienceSubmission, error) {
	if err := journey.Complete(); err != nil {
		return domain.ExperienceSubmission{}, err
	}
	if s.store == nil {
		return domain.ExperienceSubmission{}, fmt.Errorf("store is not configured")
	}
	return s.store.SaveExperience(ctx, domain.ExperienceSubmission{Name: journey.VisitorName, Contact: "journey", Interest: journey.Interest, Message: journey.Summary()})
}

func (s *Service) AuditJourney(ctx context.Context, item domain.ExperienceSubmission) error {
	if s.store == nil {
		return fmt.Errorf("store is not configured")
	}
	return s.store.RecordAudit(ctx, store.AuditEntry{Action: "journey-saved", Entity: "ExperienceSubmission", EntityID: fmt.Sprint(item.ID), Detail: item.Summary(), CreatedAt: item.CreatedAt})
}
