package store

import (
	"context"
	"path/filepath"
	"testing"

	"qin-culture-site/internal/domain"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "qin.sqlite")
	ctx := context.Background()
	first, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	saved, err := first.SaveExperience(ctx, domain.ExperienceSubmission{Name: "山水", Contact: "13800000000", Interest: "名曲聆听", Message: "流水"})
	if err != nil {
		first.Close()
		t.Fatal(err)
	}
	if err := first.Close(); err != nil {
		t.Fatal(err)
	}
	second, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()
	loaded, err := second.FindExperience(ctx, saved.ID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Name != "山水" || loaded.Status != "pending" {
		t.Fatalf("persistence mismatch: %#v", loaded)
	}
}
