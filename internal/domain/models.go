package domain

import (
	"fmt"
	"strings"
)

type QinSchool struct {
	ID          string
	Name        string
	Region      string
	Description string
	Founder     string
	Period      string
}

type QinPiece struct {
	ID              string
	Title           string
	Composer        string
	SchoolID        string
	AudioPath       string
	DurationSeconds int
	Mood            string
	Summary         string
}

type Notation struct {
	ID         string
	PieceID    string
	Label      string
	Excerpt    string
	Difficulty string
	Technique  string
}

type Courtesy struct {
	ID      string
	Title   string
	Content string
	Order   int
	Stage   string
}

type HeritageStory struct {
	ID       string
	Title    string
	Era      string
	Body     string
	Featured bool
	Place    string
}

type ExperienceSubmission struct {
	ID        int64
	Name      string
	Contact   string
	Interest  string
	Message   string
	Status    string
	CreatedAt string
}

type StudyBundle struct {
	Piece     QinPiece
	Notations []Notation
	Courtesy  []Courtesy
	Stories   []HeritageStory
}

func (s QinSchool) Slug() string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(s.ID), " ", "-"))
}

func (p QinPiece) Label() string {
	if p.Composer == "" {
		return p.Title
	}
	return fmt.Sprintf("%s · %s", p.Title, p.Composer)
}

func (p QinPiece) DurationLabel() string {
	minutes := p.DurationSeconds / 60
	seconds := p.DurationSeconds % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

func (n Notation) IsAdvanced() bool {
	return strings.EqualFold(n.Difficulty, "高级") || strings.EqualFold(n.Difficulty, "进阶")
}

func (c Courtesy) DisplayOrder() string {
	return fmt.Sprintf("%02d · %s", c.Order, c.Stage)
}

func (h HeritageStory) Excerpt(limit int) string {
	text := strings.TrimSpace(h.Body)
	if limit <= 0 || len([]rune(text)) <= limit {
		return text
	}
	runes := []rune(text)
	return string(runes[:limit]) + "…"
}

func (e ExperienceSubmission) IsPending() bool {
	return e.Status == "pending"
}

func (e ExperienceSubmission) Summary() string {
	parts := []string{strings.TrimSpace(e.Name), strings.TrimSpace(e.Interest)}
	return strings.Join(parts, " / ")
}

func (b StudyBundle) Complete() bool {
	return b.Piece.ID != "" && len(b.Notations) > 0 && len(b.Courtesy) > 0 && len(b.Stories) > 0
}
