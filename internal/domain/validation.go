package domain

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNameRequired     = errors.New("name is required")
	ErrContactRequired  = errors.New("contact is required")
	ErrInterestRequired = errors.New("interest is required")
	ErrPieceRequired    = errors.New("piece is required")
)

func (e ExperienceSubmission) Validate() error {
	if strings.TrimSpace(e.Name) == "" {
		return ErrNameRequired
	}
	if strings.TrimSpace(e.Contact) == "" {
		return ErrContactRequired
	}
	if strings.TrimSpace(e.Interest) == "" {
		return ErrInterestRequired
	}
	if len([]rune(e.Name)) > 80 {
		return fmt.Errorf("name exceeds 80 characters")
	}
	if len([]rune(e.Contact)) > 120 {
		return fmt.Errorf("contact exceeds 120 characters")
	}
	if len([]rune(e.Message)) > 1000 {
		return fmt.Errorf("message exceeds 1000 characters")
	}
	return nil
}

func NormalizeSubmission(input ExperienceSubmission) (ExperienceSubmission, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Contact = strings.TrimSpace(input.Contact)
	input.Interest = strings.TrimSpace(input.Interest)
	input.Message = strings.TrimSpace(input.Message)
	if err := input.Validate(); err != nil {
		return ExperienceSubmission{}, err
	}
	if input.Status == "" {
		input.Status = "pending"
	}
	if input.Status != "pending" && input.Status != "reviewed" && input.Status != "archived" {
		return ExperienceSubmission{}, fmt.Errorf("unknown status %q", input.Status)
	}
	return input, nil
}

func ValidatePieceSelection(pieces []QinPiece, id string) (QinPiece, error) {
	for _, piece := range pieces {
		if piece.ID == id {
			return piece, nil
		}
	}
	return QinPiece{}, fmt.Errorf("%w: %s", ErrPieceRequired, id)
}

func ValidateStudyQuery(query string) string {
	return strings.TrimSpace(strings.ToLower(query))
}

func MatchStudyText(query string, values ...string) bool {
	term := ValidateStudyQuery(query)
	if term == "" {
		return true
	}
	for _, value := range values {
		if strings.Contains(strings.ToLower(value), term) {
			return true
		}
	}
	return false
}
