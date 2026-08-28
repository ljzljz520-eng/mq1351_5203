package catalog

import (
	"sort"
	"strings"

	"qin-culture-site/internal/domain"
)

type SearchResult struct {
	Hits  []domain.SearchHit
	Stats domain.CultureStats
}

func (c *Catalog) Stats() domain.CultureStats {
	return domain.BuildStats(c.schools, c.pieces, c.notations, c.courtesy, c.stories)
}

func (c *Catalog) Search(query string) SearchResult {
	hits := make([]domain.SearchHit, 0)
	term := strings.TrimSpace(query)
	for _, school := range c.schools {
		weight := 2
		if strings.Contains(strings.ToLower(school.Name+school.Region), strings.ToLower(term)) {
			weight = 5
		}
		hits = append(hits, domain.NewSearchHit("school", school.ID, school.Name, school.Description, weight))
	}
	for _, piece := range c.pieces {
		weight := 3
		if strings.Contains(strings.ToLower(piece.Title+piece.Composer), strings.ToLower(term)) {
			weight = 6
		}
		hits = append(hits, domain.NewSearchHit("piece", piece.ID, piece.Title, piece.Summary, weight))
	}
	for _, story := range c.stories {
		hits = append(hits, domain.NewSearchHit("story", story.ID, story.Title, story.Body, 1))
	}
	return SearchResult{Hits: domain.FilterHits(hits, term), Stats: c.Stats()}
}

func (c *Catalog) SchoolsByRegion(region string) []domain.QinSchool {
	term := strings.ToLower(strings.TrimSpace(region))
	items := make([]domain.QinSchool, 0)
	for _, school := range c.schools {
		if term == "" || strings.Contains(strings.ToLower(school.Region), term) {
			items = append(items, school)
		}
	}
	sort.SliceStable(items, func(i, j int) bool { return items[i].Region < items[j].Region })
	return items
}

func (c *Catalog) PiecesBySchool(schoolID string) []domain.QinPiece {
	items := make([]domain.QinPiece, 0)
	for _, piece := range c.pieces {
		if piece.SchoolID == schoolID {
			items = append(items, piece)
		}
	}
	sort.SliceStable(items, func(i, j int) bool { return items[i].DurationSeconds < items[j].DurationSeconds })
	return items
}

func (c *Catalog) PiecesByMood(mood string) []domain.QinPiece {
	term := strings.ToLower(strings.TrimSpace(mood))
	items := make([]domain.QinPiece, 0)
	for _, piece := range c.pieces {
		if term == "" || strings.Contains(strings.ToLower(piece.Mood), term) {
			items = append(items, piece)
		}
	}
	return items
}

func (c *Catalog) Timeline() []domain.TimelineEntry {
	return domain.Timeline(c.stories)
}

func (c *Catalog) NotationsForPiece(pieceID string) []domain.Notation {
	items := []domain.Notation{}
	for _, notation := range c.notations {
		if notation.PieceID == pieceID {
			items = append(items, notation)
		}
	}
	return items
}
