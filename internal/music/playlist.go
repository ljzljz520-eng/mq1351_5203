package music

import (
	"sort"
	"strings"

	"qin-culture-site/internal/domain"
)

type Playlist struct {
	items []domain.QinPiece
}

func NewPlaylist(items []domain.QinPiece) Playlist {
	copyItems := append([]domain.QinPiece(nil), items...)
	sort.SliceStable(copyItems, func(i, j int) bool { return copyItems[i].Title < copyItems[j].Title })
	return Playlist{items: copyItems}
}

func (p Playlist) Items() []domain.QinPiece {
	return append([]domain.QinPiece(nil), p.items...)
}

func (p Playlist) Search(query string) []domain.QinPiece {
	term := strings.ToLower(strings.TrimSpace(query))
	if term == "" {
		return p.Items()
	}
	result := make([]domain.QinPiece, 0)
	for _, item := range p.items {
		if strings.Contains(strings.ToLower(item.Title), term) || strings.Contains(strings.ToLower(item.Composer), term) || strings.Contains(strings.ToLower(item.Mood), term) {
			result = append(result, item)
		}
	}
	return result
}

func (p Playlist) TotalDuration() int {
	total := 0
	for _, item := range p.items {
		total += item.DurationSeconds
	}
	return total
}

func (p Playlist) At(index int) (domain.QinPiece, bool) {
	if index < 0 || index >= len(p.items) {
		return domain.QinPiece{}, false
	}
	return p.items[index], true
}
