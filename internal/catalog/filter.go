package catalog

import (
	"sort"
	"strings"

	"qin-culture-site/internal/domain"
)

type StudyFilter struct {
	PieceID    string
	Query      string
	Difficulty string
}

func (c *Catalog) FilterStudy(filter StudyFilter) domain.StudyBundle {
	piece, found := c.Piece(filter.PieceID)
	if !found {
		return domain.StudyBundle{}
	}
	term := strings.ToLower(strings.TrimSpace(filter.Query))
	notations := make([]domain.Notation, 0)
	for _, notation := range c.notations {
		if notation.PieceID != piece.ID {
			continue
		}
		if filter.Difficulty != "" && notation.Difficulty != filter.Difficulty {
			continue
		}
		if term != "" && !strings.Contains(strings.ToLower(notation.Label+notation.Excerpt+notation.Technique), term) {
			continue
		}
		notations = append(notations, notation)
	}
	sort.SliceStable(notations, func(i, j int) bool { return notations[i].ID < notations[j].ID })
	return domain.StudyBundle{Piece: piece, Notations: notations, Courtesy: c.Courtesy(), Stories: c.Stories()}
}

func (c *Catalog) DifficultyOptions() []string {
	seen := map[string]bool{}
	options := []string{}
	for _, notation := range c.notations {
		if !seen[notation.Difficulty] {
			seen[notation.Difficulty] = true
			options = append(options, notation.Difficulty)
		}
	}
	sort.Strings(options)
	return options
}

func (c *Catalog) FeaturedStories() []domain.HeritageStory {
	items := []domain.HeritageStory{}
	for _, story := range c.stories {
		if story.Featured {
			items = append(items, story)
		}
	}
	return items
}
