package store

import (
	"context"
	"fmt"
)

type StorageMetrics struct {
	Experiences int
	Schools     int
	Pieces      int
	Notations   int
	Audits      int
}

func (s *Store) Metrics(ctx context.Context) (StorageMetrics, error) {
	counts, err := s.CountCultureRows(ctx)
	if err != nil {
		return StorageMetrics{}, err
	}
	var experiences, audits int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM experiences`).Scan(&experiences); err != nil {
		return StorageMetrics{}, fmt.Errorf("count experiences: %w", err)
	}
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM audit_entries`).Scan(&audits); err != nil {
		return StorageMetrics{}, fmt.Errorf("count audits: %w", err)
	}
	return StorageMetrics{Experiences: experiences, Schools: counts["qin_schools"], Pieces: counts["qin_pieces"], Notations: counts["notations"], Audits: audits}, nil
}

func (m StorageMetrics) Total() int {
	return m.Experiences + m.Schools + m.Pieces + m.Notations + m.Audits
}

func (m StorageMetrics) Healthy() bool {
	return m.Schools >= 0 && m.Pieces >= 0 && m.Notations >= 0 && m.Experiences >= 0 && m.Audits >= 0
}

func (m StorageMetrics) Labels() []string {
	return []string{fmt.Sprintf("体验 %d", m.Experiences), fmt.Sprintf("琴派 %d", m.Schools), fmt.Sprintf("曲目 %d", m.Pieces), fmt.Sprintf("谱例 %d", m.Notations), fmt.Sprintf("审计 %d", m.Audits)}
}
