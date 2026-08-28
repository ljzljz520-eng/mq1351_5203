package store

import (
	"context"
	"fmt"

	"qin-culture-site/internal/domain"
)

func (s *Store) SaveQinSchool(ctx context.Context, school domain.QinSchool) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO qin_schools(id, name, region, description) VALUES(?, ?, ?, ?) ON CONFLICT(id) DO UPDATE SET name=excluded.name, region=excluded.region, description=excluded.description`, school.ID, school.Name, school.Region, school.Description)
	if err != nil {
		return fmt.Errorf("save QinSchool: %w", err)
	}
	return nil
}

func (s *Store) SaveQinPiece(ctx context.Context, piece domain.QinPiece) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO qin_pieces(id, title, composer, school_id, audio_path, duration_seconds) VALUES(?, ?, ?, ?, ?, ?) ON CONFLICT(id) DO UPDATE SET title=excluded.title, composer=excluded.composer, school_id=excluded.school_id, audio_path=excluded.audio_path, duration_seconds=excluded.duration_seconds`, piece.ID, piece.Title, piece.Composer, piece.SchoolID, piece.AudioPath, piece.DurationSeconds)
	if err != nil {
		return fmt.Errorf("save QinPiece: %w", err)
	}
	return nil
}

func (s *Store) SaveNotation(ctx context.Context, notation domain.Notation) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO notations(id, piece_id, label, excerpt, difficulty) VALUES(?, ?, ?, ?, ?) ON CONFLICT(id) DO UPDATE SET piece_id=excluded.piece_id, label=excluded.label, excerpt=excluded.excerpt, difficulty=excluded.difficulty`, notation.ID, notation.PieceID, notation.Label, notation.Excerpt, notation.Difficulty)
	if err != nil {
		return fmt.Errorf("save Notation: %w", err)
	}
	return nil
}

func (s *Store) CountCultureRows(ctx context.Context) (map[string]int, error) {
	counts := map[string]int{}
	for _, table := range []string{"qin_schools", "qin_pieces", "notations"} {
		var count int
		if err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+table).Scan(&count); err != nil {
			return nil, fmt.Errorf("count %s: %w", table, err)
		}
		counts[table] = count
	}
	return counts, nil
}
