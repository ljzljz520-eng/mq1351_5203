package store

import (
	"context"
	"fmt"
)

type AuditEntry struct {
	ID        int64
	Action    string
	Entity    string
	EntityID  string
	Detail    string
	CreatedAt string
}

func (s *Store) RecordAudit(ctx context.Context, entry AuditEntry) error {
	if entry.Action == "" || entry.Entity == "" {
		return fmt.Errorf("audit action and entity are required")
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO audit_entries(action, entity, entity_id, detail, created_at) VALUES(?, ?, ?, ?, ?)`, entry.Action, entry.Entity, entry.EntityID, entry.Detail, entry.CreatedAt)
	if err != nil {
		return fmt.Errorf("record audit: %w", err)
	}
	return nil
}

func (s *Store) ListAudit(ctx context.Context, entity string) ([]AuditEntry, error) {
	query := `SELECT id, action, entity, entity_id, detail, created_at FROM audit_entries ORDER BY id`
	args := []any{}
	if entity != "" {
		query = `SELECT id, action, entity, entity_id, detail, created_at FROM audit_entries WHERE entity = ? ORDER BY id`
		args = append(args, entity)
	}
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query audit: %w", err)
	}
	defer rows.Close()
	items := []AuditEntry{}
	for rows.Next() {
		var entry AuditEntry
		if err := rows.Scan(&entry.ID, &entry.Action, &entry.Entity, &entry.EntityID, &entry.Detail, &entry.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan audit: %w", err)
		}
		items = append(items, entry)
	}
	return items, rows.Err()
}

func (s *Store) AuditCount(ctx context.Context) (int, error) {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM audit_entries`).Scan(&count); err != nil {
		return 0, fmt.Errorf("count audit: %w", err)
	}
	return count, nil
}
