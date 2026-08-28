package service

import (
	"context"
	"fmt"
	"sort"

	"qin-culture-site/internal/domain"
)

type Inspection struct {
	Ready       bool
	Summary     string
	Warnings    []string
	StorageHint string
}

func (s *Service) Inspect(ctx context.Context) Inspection {
	inspection := Inspection{Warnings: []string{}}
	if s.catalog == nil {
		inspection.Warnings = append(inspection.Warnings, "目录未连接")
		return inspection
	}
	stats := s.catalog.Stats()
	inspection.Summary = domain.CompactSummary(stats)
	if !stats.HasMinimumContent() {
		inspection.Warnings = append(inspection.Warnings, "文化资料尚未达到完整展示线")
	}
	if s.store == nil {
		inspection.Warnings = append(inspection.Warnings, "体验提交不会持久化")
		inspection.StorageHint = "memory"
		return inspection
	}
	metrics, err := s.store.Metrics(ctx)
	if err != nil {
		inspection.Warnings = append(inspection.Warnings, err.Error())
		inspection.StorageHint = "unavailable"
		return inspection
	}
	if !metrics.Healthy() {
		inspection.Warnings = append(inspection.Warnings, "存储指标异常")
	}
	inspection.StorageHint = fmt.Sprintf("%d 条持久化记录", metrics.Total())
	inspection.Ready = len(inspection.Warnings) == 0
	return inspection
}

func (s *Service) PieceRecommendations(pieceID string, limit int) []domain.QinPiece {
	if s.catalog == nil || limit <= 0 {
		return []domain.QinPiece{}
	}
	selected, found := s.catalog.Piece(pieceID)
	if !found {
		return []domain.QinPiece{}
	}
	items := s.catalog.PiecesByMood(selected.Mood)
	result := make([]domain.QinPiece, 0, limit)
	for _, item := range items {
		if item.ID == selected.ID {
			continue
		}
		result = append(result, item)
		if len(result) >= limit {
			break
		}
	}
	if len(result) < limit {
		for _, item := range s.catalog.SortedPieces() {
			if item.ID == selected.ID {
				continue
			}
			seen := false
			for _, current := range result {
				if current.ID == item.ID {
					seen = true
					break
				}
			}
			if !seen {
				result = append(result, item)
			}
			if len(result) >= limit {
				break
			}
		}
	}
	sort.SliceStable(result, func(i, j int) bool { return result[i].Title < result[j].Title })
	return result
}

func (s *Service) PieceLabel(id string) string {
	if s.catalog == nil {
		return ""
	}
	piece, ok := s.catalog.Piece(id)
	if !ok {
		return ""
	}
	return piece.Label()
}
