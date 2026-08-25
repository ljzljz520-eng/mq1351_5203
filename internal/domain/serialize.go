package domain

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ExportRow struct {
	Category string
	ID       string
	Title    string
	Detail   string
}

func PieceExportRows(pieces []QinPiece) []ExportRow {
	rows := make([]ExportRow, 0, len(pieces))
	for _, piece := range pieces {
		rows = append(rows, ExportRow{Category: "名曲", ID: piece.ID, Title: piece.Title, Detail: piece.Label()})
	}
	return rows
}

func SchoolExportRows(schools []QinSchool) []ExportRow {
	rows := make([]ExportRow, 0, len(schools))
	for _, school := range schools {
		rows = append(rows, ExportRow{Category: "琴派", ID: school.ID, Title: school.Name, Detail: school.Region + " · " + school.Period})
	}
	return rows
}

func StoryExportRows(stories []HeritageStory) []ExportRow {
	rows := make([]ExportRow, 0, len(stories))
	for _, story := range stories {
		rows = append(rows, ExportRow{Category: "传承故事", ID: story.ID, Title: story.Title, Detail: story.Era + " · " + story.Place})
	}
	return rows
}

func ExportJSON(rows []ExportRow) ([]byte, error) {
	if rows == nil {
		rows = []ExportRow{}
	}
	return json.MarshalIndent(rows, "", "  ")
}

func ExportCSV(rows []ExportRow) string {
	lines := []string{"category,id,title,detail"}
	for _, row := range rows {
		lines = append(lines, strings.Join([]string{csvCell(row.Category), csvCell(row.ID), csvCell(row.Title), csvCell(row.Detail)}, ","))
	}
	return strings.Join(lines, "\n") + "\n"
}

func csvCell(value string) string {
	value = strings.ReplaceAll(value, "\"", "\"\"")
	if strings.ContainsAny(value, ",\n\r\"") {
		return fmt.Sprintf("\"%s\"", value)
	}
	return value
}

func HumanCount(value int) string {
	if value == 1 {
		return "1 条"
	}
	return fmt.Sprintf("%d 条", value)
}

func CompactSummary(stats CultureStats) string {
	return strings.Join([]string{fmt.Sprintf("琴派 %d", stats.SchoolCount), fmt.Sprintf("名曲 %d", stats.PieceCount), fmt.Sprintf("谱例 %d", stats.NotationCount), fmt.Sprintf("故事 %d", stats.StoryCount)}, " / ")
}
