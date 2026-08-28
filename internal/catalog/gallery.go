package catalog

import (
	"fmt"
	"sort"
	"strings"

	"qin-culture-site/internal/domain"
)

type GalleryItem struct {
	ID         string
	Title      string
	Category   string
	Caption    string
	ImageHint  string
	Accent     string
	SortWeight int
}

func (c *Catalog) Gallery() []GalleryItem {
	items := make([]GalleryItem, 0, len(c.schools)+len(c.stories)+len(c.pieces))
	for _, school := range c.schools {
		items = append(items, GalleryItem{ID: school.ID, Title: school.Name, Category: "琴派", Caption: school.Description, ImageHint: school.Region + "山水", Accent: "松烟", SortWeight: 10})
	}
	for _, story := range c.stories {
		weight := 5
		if story.Featured {
			weight = 20
		}
		items = append(items, GalleryItem{ID: story.ID, Title: story.Title, Category: "传承", Caption: story.Excerpt(140), ImageHint: story.Place, Accent: "古木", SortWeight: weight})
	}
	for _, piece := range c.pieces {
		items = append(items, GalleryItem{ID: piece.ID, Title: piece.Title, Category: "名曲", Caption: piece.Summary, ImageHint: piece.Mood, Accent: "琴弦", SortWeight: 8})
	}
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].SortWeight == items[j].SortWeight {
			return items[i].Title < items[j].Title
		}
		return items[i].SortWeight > items[j].SortWeight
	})
	return items
}

func (c *Catalog) GalleryByCategory(category string) []GalleryItem {
	term := strings.TrimSpace(category)
	items := []GalleryItem{}
	for _, item := range c.Gallery() {
		if term == "" || item.Category == term {
			items = append(items, item)
		}
	}
	return items
}

func (g GalleryItem) AltText() string {
	return fmt.Sprintf("%s：%s，意象为%s", g.Category, g.Title, g.ImageHint)
}

func (g GalleryItem) IsFeatured() bool {
	return g.SortWeight >= 20
}

func GalleryCategories(items []GalleryItem) []string {
	seen := map[string]bool{}
	result := []string{}
	for _, item := range items {
		if !seen[item.Category] {
			seen[item.Category] = true
			result = append(result, item.Category)
		}
	}
	sort.Strings(result)
	return result
}

func GalleryStats(items []GalleryItem) domain.CultureStats {
	stats := domain.CultureStats{}
	for _, item := range items {
		switch item.Category {
		case "琴派":
			stats.SchoolCount++
		case "名曲":
			stats.PieceCount++
		case "传承":
			stats.StoryCount++
		}
	}
	return stats
}
