package store

import (
	"context"
	"fmt"

	"qin-culture-site/internal/domain"
)

func schemaStatements() []string {
	return []string{
		`PRAGMA foreign_keys = ON`,
		`CREATE TABLE IF NOT EXISTS experiences (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL, contact TEXT NOT NULL, interest TEXT NOT NULL, message TEXT NOT NULL, status TEXT NOT NULL, created_at TEXT NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS qin_schools (id TEXT PRIMARY KEY, name TEXT NOT NULL, region TEXT NOT NULL, description TEXT NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS qin_pieces (id TEXT PRIMARY KEY, title TEXT NOT NULL, composer TEXT NOT NULL, school_id TEXT NOT NULL, audio_path TEXT NOT NULL, duration_seconds INTEGER NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS notations (id TEXT PRIMARY KEY, piece_id TEXT NOT NULL, label TEXT NOT NULL, excerpt TEXT NOT NULL, difficulty TEXT NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS audit_entries (id INTEGER PRIMARY KEY AUTOINCREMENT, action TEXT NOT NULL, entity TEXT NOT NULL, entity_id TEXT NOT NULL, detail TEXT NOT NULL, created_at TEXT NOT NULL)`,
	}
}

func (s *Store) SeedCulture(ctx context.Context, schools []domain.QinSchool, pieces []domain.QinPiece, notations []domain.Notation) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin culture seed: %w", err)
	}
	for _, school := range schools {
		if _, err := tx.ExecContext(ctx, `INSERT INTO qin_schools(id, name, region, description) VALUES(?, ?, ?, ?) ON CONFLICT(id) DO NOTHING`, school.ID, school.Name, school.Region, school.Description); err != nil {
			tx.Rollback()
			return fmt.Errorf("seed QinSchool: %w", err)
		}
	}
	for _, piece := range pieces {
		if _, err := tx.ExecContext(ctx, `INSERT INTO qin_pieces(id, title, composer, school_id, audio_path, duration_seconds) VALUES(?, ?, ?, ?, ?, ?) ON CONFLICT(id) DO NOTHING`, piece.ID, piece.Title, piece.Composer, piece.SchoolID, piece.AudioPath, piece.DurationSeconds); err != nil {
			tx.Rollback()
			return fmt.Errorf("seed QinPiece: %w", err)
		}
	}
	for _, notation := range notations {
		if _, err := tx.ExecContext(ctx, `INSERT INTO notations(id, piece_id, label, excerpt, difficulty) VALUES(?, ?, ?, ?, ?) ON CONFLICT(id) DO NOTHING`, notation.ID, notation.PieceID, notation.Label, notation.Excerpt, notation.Difficulty); err != nil {
			tx.Rollback()
			return fmt.Errorf("seed Notation: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit culture seed: %w", err)
	}
	return nil
}
