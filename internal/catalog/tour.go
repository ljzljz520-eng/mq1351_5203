package catalog

import (
	"fmt"
	"sort"

	"qin-culture-site/internal/domain"
)

type TourStop struct {
	Order       int
	Title       string
	Description string
	Href        string
	Duration    int
}

type CuratedTour struct {
	ID          string
	Name        string
	Audience    string
	Description string
	Stops       []TourStop
}

func (c *Catalog) Tours() []CuratedTour {
	return []CuratedTour{
		{ID: "first-listen", Name: "第一次听琴", Audience: "初次来访", Description: "用三首气质不同的曲目认识古琴的留白与呼吸。", Stops: c.tourStops([]string{"gaoshan", "meihua", "liushui"})},
		{ID: "craft-and-lineage", Name: "手艺与传承", Audience: "研究者", Description: "从琴派地域、斫琴手艺和传承故事理解古琴生态。", Stops: c.tourStops([]string{"xiaoxiang", "yangchun", "woye"})},
		{ID: "notation-practice", Name: "减字谱练习", Audience: "学习者", Description: "按从泛音到吟猱的难度顺序安排练习线索。", Stops: c.tourStops([]string{"gaoshan", "liushui", "meihua"})},
	}
}

func (c *Catalog) tourStops(ids []string) []TourStop {
	stops := make([]TourStop, 0, len(ids))
	for index, id := range ids {
		piece, ok := c.Piece(id)
		if !ok {
			continue
		}
		stops = append(stops, TourStop{Order: index + 1, Title: piece.Title, Description: piece.Summary, Href: "#pieces-" + piece.ID, Duration: piece.DurationSeconds})
	}
	return stops
}

func (t CuratedTour) TotalMinutes() int {
	total := 0
	for _, stop := range t.Stops {
		total += stop.Duration
	}
	return (total + 59) / 60
}

func (t CuratedTour) Labels() []string {
	labels := make([]string, 0, len(t.Stops))
	for _, stop := range t.Stops {
		labels = append(labels, fmt.Sprintf("%d. %s", stop.Order, stop.Title))
	}
	return labels
}

func SortTours(tours []CuratedTour) []CuratedTour {
	result := append([]CuratedTour(nil), tours...)
	sort.SliceStable(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result
}

func TourStats(tours []CuratedTour) domain.CultureStats {
	stats := domain.CultureStats{}
	for _, tour := range tours {
		stats.PieceCount += len(tour.Stops)
	}
	return stats
}
