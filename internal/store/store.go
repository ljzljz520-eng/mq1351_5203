package store

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
	"qin-culture-site/internal/domain"
)

type Store struct {
	db   *sql.DB
	path string
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, fmt.Errorf("database path is required")
	}
	if path != ":memory:" {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return nil, fmt.Errorf("create database directory: %w", err)
		}
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	store := &Store{db: db, path: path}
	if err := store.initialize(context.Background()); err != nil {
		db.Close()
		return nil, err
	}
	return store, nil
}

func (s *Store) initialize(ctx context.Context) error {
	for _, statement := range schemaStatements() {
		if _, err := s.db.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("initialize store: %w", err)
		}
	}
	return nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Store) Path() string {
	return s.path
}

func (s *Store) Ping(ctx context.Context) error {
	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping sqlite: %w", err)
	}
	return nil
}

func (s *Store) SaveExperience(ctx context.Context, input domain.ExperienceSubmission) (domain.ExperienceSubmission, error) {
	normalized, err := domain.NormalizeSubmission(input)
	if err != nil {
		return domain.ExperienceSubmission{}, err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.ExperienceSubmission{}, fmt.Errorf("begin experience transaction: %w", err)
	}
	result, err := tx.ExecContext(ctx, `INSERT INTO experiences(name, contact, interest, message, status, created_at) VALUES(?, ?, ?, ?, ?, ?)`, normalized.Name, normalized.Contact, normalized.Interest, normalized.Message, normalized.Status, normalized.CreatedAt)
	if err != nil {
		tx.Rollback()
		return domain.ExperienceSubmission{}, fmt.Errorf("insert experience: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return domain.ExperienceSubmission{}, fmt.Errorf("commit experience: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return domain.ExperienceSubmission{}, fmt.Errorf("read experience id: %w", err)
	}
	normalized.ID = id
	return normalized, nil
}

func (s *Store) ListExperiences(ctx context.Context) ([]domain.ExperienceSubmission, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, name, contact, interest, message, status, created_at FROM experiences ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("query experiences: %w", err)
	}
	defer rows.Close()
	items := make([]domain.ExperienceSubmission, 0)
	for rows.Next() {
		var item domain.ExperienceSubmission
		if err := rows.Scan(&item.ID, &item.Name, &item.Contact, &item.Interest, &item.Message, &item.Status, &item.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan experience: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read experience rows: %w", err)
	}
	return items, nil
}

func (s *Store) FindExperience(ctx context.Context, id int64) (domain.ExperienceSubmission, error) {
	var item domain.ExperienceSubmission
	err := s.db.QueryRowContext(ctx, `SELECT id, name, contact, interest, message, status, created_at FROM experiences WHERE id = ?`, id).Scan(&item.ID, &item.Name, &item.Contact, &item.Interest, &item.Message, &item.Status, &item.CreatedAt)
	if err != nil {
		return domain.ExperienceSubmission{}, fmt.Errorf("find experience %d: %w", id, err)
	}
	return item, nil
}
