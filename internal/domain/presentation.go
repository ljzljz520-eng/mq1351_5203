package domain

import (
	"fmt"
	"html/template"
	"strings"
)

type Card struct {
	Heading string
	Kicker  string
	Text    string
	Link    string
}

type TableRow struct {
	Cells []string
}

func SchoolCard(s QinSchool) Card {
	return Card{Heading: s.Name, Kicker: s.Region, Text: s.Description, Link: "#schools"}
}

func PieceRow(p QinPiece) TableRow {
	return TableRow{Cells: []string{p.Title, p.Composer, p.Mood, p.DurationLabel(), p.AudioPath}}
}

func NotationRow(n Notation) TableRow {
	return TableRow{Cells: []string{n.Label, n.Difficulty, n.Technique, n.Excerpt}}
}

func StoryCard(h HeritageStory) Card {
	kicker := h.Era
	if h.Place != "" {
		kicker += " · " + h.Place
	}
	return Card{Heading: h.Title, Kicker: kicker, Text: h.Excerpt(110), Link: "#stories"}
}

func BuildAltText(piece QinPiece) string {
	return fmt.Sprintf("古琴曲目%s，演奏者%s，气质%s", piece.Title, piece.Composer, piece.Mood)
}

func SafeText(text string) template.HTML {
	return template.HTML(template.HTMLEscapeString(strings.TrimSpace(text)))
}
