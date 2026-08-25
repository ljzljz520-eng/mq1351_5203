package service

import (
	"fmt"
	"sort"
	"strings"

	"qin-culture-site/internal/domain"
)

type NavigationItem struct {
	Label string
	Href  string
	Items []NavigationItem
}

func Navigation() []NavigationItem {
	return []NavigationItem{
		{Label: "琴学入门", Href: "#intro", Items: []NavigationItem{{Label: "琴派", Href: "#schools"}, {Label: "礼仪", Href: "#courtesy"}}},
		{Label: "聆听名曲", Href: "#pieces", Items: []NavigationItem{{Label: "曲目表", Href: "#pieces"}, {Label: "播放体验", Href: "#player"}}},
		{Label: "传承脉络", Href: "#stories", Items: []NavigationItem{{Label: "故事", Href: "#stories"}, {Label: "减字谱", Href: "#notation"}}},
	}
}

func PieceOptions(pieces []domain.QinPiece) []string {
	options := make([]string, 0, len(pieces))
	for _, piece := range pieces {
		options = append(options, fmt.Sprintf("%s|%s", piece.ID, piece.Title))
	}
	sort.Strings(options)
	return options
}

func JoinLabels(values []string) string {
	clean := make([]string, 0, len(values))
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			clean = append(clean, strings.TrimSpace(value))
		}
	}
	return strings.Join(clean, " · ")
}
