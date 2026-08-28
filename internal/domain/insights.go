package domain

import (
	"fmt"
	"sort"
	"strings"
)

type CultureStats struct {
	SchoolCount       int
	PieceCount        int
	NotationCount     int
	CourtesyCount     int
	StoryCount        int
	FeaturedStories   int
	TotalAudioSeconds int
}

type SearchHit struct {
	Kind    string
	ID      string
	Title   string
	Excerpt string
	Weight  int
}

type TimelineEntry struct {
	Label string
	Era   string
	Place string
	Body  string
}

func BuildStats(schools []QinSchool, pieces []QinPiece, notations []Notation, courtesy []Courtesy, stories []HeritageStory) CultureStats {
	stats := CultureStats{SchoolCount: len(schools), PieceCount: len(pieces), NotationCount: len(notations), CourtesyCount: len(courtesy), StoryCount: len(stories)}
	for _, piece := range pieces {
		stats.TotalAudioSeconds += piece.DurationSeconds
	}
	for _, story := range stories {
		if story.Featured {
			stats.FeaturedStories++
		}
	}
	return stats
}

func (s CultureStats) AudioMinutes() string {
	return fmt.Sprintf("%d 分钟", (s.TotalAudioSeconds+59)/60)
}

func (s CultureStats) HasMinimumContent() bool {
	return s.SchoolCount >= 4 && s.PieceCount >= 6 && s.NotationCount >= 4 && s.CourtesyCount >= 3 && s.StoryCount >= 3
}

func (s CultureStats) Labels() []string {
	return []string{fmt.Sprintf("%d 个琴派", s.SchoolCount), fmt.Sprintf("%d 首名曲", s.PieceCount), fmt.Sprintf("%d 条谱例", s.NotationCount), s.AudioMinutes()}
}

func NewSearchHit(kind, id, title, excerpt string, weight int) SearchHit {
	return SearchHit{Kind: kind, ID: id, Title: strings.TrimSpace(title), Excerpt: strings.TrimSpace(excerpt), Weight: weight}
}

func RankHits(hits []SearchHit) []SearchHit {
	result := append([]SearchHit(nil), hits...)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Weight == result[j].Weight {
			return result[i].Title < result[j].Title
		}
		return result[i].Weight > result[j].Weight
	})
	return result
}

func FilterHits(hits []SearchHit, query string) []SearchHit {
	term := strings.ToLower(strings.TrimSpace(query))
	if term == "" {
		return RankHits(hits)
	}
	filtered := make([]SearchHit, 0)
	for _, hit := range hits {
		if strings.Contains(strings.ToLower(hit.Title), term) || strings.Contains(strings.ToLower(hit.Excerpt), term) {
			filtered = append(filtered, hit)
		}
	}
	return RankHits(filtered)
}

func Timeline(stories []HeritageStory) []TimelineEntry {
	items := make([]TimelineEntry, 0, len(stories))
	for _, story := range stories {
		items = append(items, TimelineEntry{Label: story.Title, Era: story.Era, Place: story.Place, Body: story.Excerpt(160)})
	}
	sort.SliceStable(items, func(i, j int) bool { return items[i].Era < items[j].Era })
	return items
}

func ValidateCulture(schools []QinSchool, pieces []QinPiece, notations []Notation) error {
	if len(schools) == 0 {
		return fmt.Errorf("no QinSchool records")
	}
	if len(pieces) == 0 {
		return fmt.Errorf("no QinPiece records")
	}
	ids := map[string]bool{}
	for _, piece := range pieces {
		if ids[piece.ID] {
			return fmt.Errorf("duplicate QinPiece %s", piece.ID)
		}
		ids[piece.ID] = true
	}
	for _, notation := range notations {
		if !ids[notation.PieceID] {
			return fmt.Errorf("notation %s references unknown piece", notation.ID)
		}
	}
	return nil
}
