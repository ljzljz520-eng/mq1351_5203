package catalog

import (
	"sort"

	"qin-culture-site/internal/domain"
)

type Catalog struct {
	schools   []domain.QinSchool
	pieces    []domain.QinPiece
	notations []domain.Notation
	courtesy  []domain.Courtesy
	stories   []domain.HeritageStory
}

func New() *Catalog {
	return &Catalog{
		schools:   fixtureSchools(),
		pieces:    fixturePieces(),
		notations: fixtureNotations(),
		courtesy:  fixtureCourtesy(),
		stories:   fixtureStories(),
	}
}

func (c *Catalog) Schools() []domain.QinSchool {
	return append([]domain.QinSchool(nil), c.schools...)
}

func (c *Catalog) Pieces() []domain.QinPiece {
	return append([]domain.QinPiece(nil), c.pieces...)
}

func (c *Catalog) Notations() []domain.Notation {
	return append([]domain.Notation(nil), c.notations...)
}

func (c *Catalog) Courtesy() []domain.Courtesy {
	return append([]domain.Courtesy(nil), c.courtesy...)
}

func (c *Catalog) Stories() []domain.HeritageStory {
	return append([]domain.HeritageStory(nil), c.stories...)
}

func (c *Catalog) Piece(id string) (domain.QinPiece, bool) {
	for _, piece := range c.pieces {
		if piece.ID == id {
			return piece, true
		}
	}
	return domain.QinPiece{}, false
}

func (c *Catalog) School(id string) (domain.QinSchool, bool) {
	for _, school := range c.schools {
		if school.ID == id {
			return school, true
		}
	}
	return domain.QinSchool{}, false
}

func (c *Catalog) SortedPieces() []domain.QinPiece {
	items := c.Pieces()
	sort.SliceStable(items, func(i, j int) bool { return items[i].Title < items[j].Title })
	return items
}

func (c *Catalog) ImageWall() []domain.Card {
	items := make([]domain.Card, 0, len(c.schools)+len(c.stories))
	for _, school := range c.schools {
		items = append(items, domain.SchoolCard(school))
	}
	for _, story := range c.stories {
		items = append(items, domain.StoryCard(story))
	}
	return items
}
