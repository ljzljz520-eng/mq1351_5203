package store

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"qin-culture-site/internal/domain"
)

type Snapshot struct {
	Version     int                           `json:"version"`
	Schools     []domain.QinSchool            `json:"schools"`
	Pieces      []domain.QinPiece             `json:"pieces"`
	Notations   []domain.Notation             `json:"notations"`
	Experiences []domain.ExperienceSubmission `json:"experiences"`
}

func (s *Store) ExportSnapshot(ctx context.Context) (Snapshot, error) {
	snapshot := Snapshot{Version: 1}
	rows, err := s.db.QueryContext(ctx, `SELECT id, name, contact, interest, message, status, created_at FROM experiences ORDER BY id`)
	if err != nil {
		return Snapshot{}, fmt.Errorf("snapshot experiences: %w", err)
	}
	for rows.Next() {
		var item domain.ExperienceSubmission
		if err := rows.Scan(&item.ID, &item.Name, &item.Contact, &item.Interest, &item.Message, &item.Status, &item.CreatedAt); err != nil {
			rows.Close()
			return Snapshot{}, fmt.Errorf("snapshot experience row: %w", err)
		}
		snapshot.Experiences = append(snapshot.Experiences, item)
	}
	if err := rows.Close(); err != nil {
		return Snapshot{}, fmt.Errorf("close experience snapshot: %w", err)
	}
	return snapshot, nil
}

func EncodeSnapshot(snapshot Snapshot) ([]byte, error) {
	if snapshot.Version == 0 {
		snapshot.Version = 1
	}
	return json.MarshalIndent(snapshot, "", "  ")
}

func DecodeSnapshot(data []byte) (Snapshot, error) {
	var snapshot Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return Snapshot{}, fmt.Errorf("decode snapshot: %w", err)
	}
	if snapshot.Version != 1 {
		return Snapshot{}, fmt.Errorf("unsupported snapshot version %d", snapshot.Version)
	}
	return snapshot, nil
}

func ValidateSnapshot(snapshot Snapshot) error {
	if snapshot.Version != 1 {
		return fmt.Errorf("snapshot version must be 1")
	}
	for _, item := range snapshot.Experiences {
		if strings.TrimSpace(item.Name) == "" || strings.TrimSpace(item.Status) == "" {
			return fmt.Errorf("snapshot contains incomplete experience")
		}
	}
	return nil
}

func SnapshotSummary(snapshot Snapshot) string {
	return fmt.Sprintf("版本 %d：%d 个琴派、%d 首名曲、%d 条谱例、%d 条体验", snapshot.Version, len(snapshot.Schools), len(snapshot.Pieces), len(snapshot.Notations), len(snapshot.Experiences))
}

func MergeSnapshots(base, incoming Snapshot) Snapshot {
	result := base
	if incoming.Version > result.Version {
		result.Version = incoming.Version
	}
	result.Schools = mergeSchools(result.Schools, incoming.Schools)
	result.Pieces = mergePieces(result.Pieces, incoming.Pieces)
	result.Notations = mergeNotations(result.Notations, incoming.Notations)
	result.Experiences = mergeExperiences(result.Experiences, incoming.Experiences)
	return result
}

func mergeSchools(base, incoming []domain.QinSchool) []domain.QinSchool {
	byID := map[string]domain.QinSchool{}
	for _, item := range base {
		byID[item.ID] = item
	}
	for _, item := range incoming {
		if item.ID != "" {
			byID[item.ID] = item
		}
	}
	result := make([]domain.QinSchool, 0, len(byID))
	for _, item := range byID {
		result = append(result, item)
	}
	return result
}

func mergePieces(base, incoming []domain.QinPiece) []domain.QinPiece {
	byID := map[string]domain.QinPiece{}
	for _, item := range base {
		byID[item.ID] = item
	}
	for _, item := range incoming {
		if item.ID != "" {
			byID[item.ID] = item
		}
	}
	result := make([]domain.QinPiece, 0, len(byID))
	for _, item := range byID {
		result = append(result, item)
	}
	return result
}

func mergeNotations(base, incoming []domain.Notation) []domain.Notation {
	byID := map[string]domain.Notation{}
	for _, item := range base {
		byID[item.ID] = item
	}
	for _, item := range incoming {
		if item.ID != "" {
			byID[item.ID] = item
		}
	}
	result := make([]domain.Notation, 0, len(byID))
	for _, item := range byID {
		result = append(result, item)
	}
	return result
}

func mergeExperiences(base, incoming []domain.ExperienceSubmission) []domain.ExperienceSubmission {
	byID := map[int64]domain.ExperienceSubmission{}
	for _, item := range base {
		byID[item.ID] = item
	}
	for _, item := range incoming {
		if item.ID != 0 {
			byID[item.ID] = item
		}
	}
	result := make([]domain.ExperienceSubmission, 0, len(byID))
	for _, item := range byID {
		result = append(result, item)
	}
	return result
}
