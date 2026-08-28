package service

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"qin-culture-site/internal/domain"
)

type CuratedSection struct {
	ID      string
	Heading string
	Intro   string
	Cards   []domain.Card
	Callout string
}

type CuratedPage struct {
	Title    string
	Subtitle string
	Sections []CuratedSection
	Stats    domain.CultureStats
}

func (s *Service) Curate(ctx context.Context, query string) (CuratedPage, error) {
	if s.catalog == nil {
		return CuratedPage{}, fmt.Errorf("catalog is not configured")
	}
	if err := domain.ValidateCulture(s.catalog.Schools(), s.catalog.Pieces(), s.catalog.Notations()); err != nil {
		return CuratedPage{}, err
	}
	model, err := s.Browse(ctx)
	if err != nil {
		return CuratedPage{}, err
	}
	result := s.catalog.Search(query)
	sections := []CuratedSection{
		{ID: "schools", Heading: "琴派地图", Intro: "从地域进入琴声的来处。", Cards: cardsFromSchools(model.Schools), Callout: "先辨地域，再听气韵"},
		{ID: "pieces", Heading: "名曲入口", Intro: "每首曲子都是一条通往山水的路径。", Cards: cardsFromPieces(model.Pieces), Callout: "选择一首，给自己三分钟"},
		{ID: "stories", Heading: "传承故事", Intro: "人和手艺让声音继续生长。", Cards: cardsFromStories(model.FeaturedStories), Callout: fmt.Sprintf("资料检索命中 %d 条", len(result.Hits))},
	}
	return CuratedPage{Title: "听见山水", Subtitle: "古琴艺术专题导览", Sections: sections, Stats: result.Stats}, nil
}

func cardsFromSchools(items []domain.QinSchool) []domain.Card {
	cards := make([]domain.Card, 0, len(items))
	for _, item := range items {
		cards = append(cards, domain.SchoolCard(item))
	}
	return cards
}

func cardsFromPieces(items []domain.QinPiece) []domain.Card {
	cards := make([]domain.Card, 0, len(items))
	for _, item := range items {
		cards = append(cards, domain.Card{Heading: item.Title, Kicker: item.Mood, Text: item.Summary, Link: "#pieces"})
	}
	return cards
}

func cardsFromStories(items []domain.HeritageStory) []domain.Card {
	cards := make([]domain.Card, 0, len(items))
	for _, item := range items {
		cards = append(cards, domain.StoryCard(item))
	}
	return cards
}

func (p CuratedPage) Section(id string) (CuratedSection, bool) {
	for _, section := range p.Sections {
		if section.ID == id {
			return section, true
		}
	}
	return CuratedSection{}, false
}

func (p CuratedPage) SearchableText() string {
	parts := []string{p.Title, p.Subtitle}
	for _, section := range p.Sections {
		parts = append(parts, section.Heading, section.Intro, section.Callout)
		for _, card := range section.Cards {
			parts = append(parts, card.Heading, card.Text)
		}
	}
	return strings.Join(parts, " ")
}

func (p CuratedPage) SortedSections() []CuratedSection {
	sections := append([]CuratedSection(nil), p.Sections...)
	sort.SliceStable(sections, func(i, j int) bool { return sections[i].ID < sections[j].ID })
	return sections
}

func (p CuratedPage) IsReady() bool {
	return len(p.Sections) >= 3 && p.Stats.HasMinimumContent()
}
