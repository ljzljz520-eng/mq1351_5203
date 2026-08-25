package service

import (
	"context"
	"testing"

	"qin-culture-site/internal/catalog"
	"qin-culture-site/internal/store"
)

func TestBrowseCultureWorkflow(t *testing.T) {
	db, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	svc := New(catalog.New(), db)
	model, err := svc.Browse(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(model.Schools) != 4 || len(model.Pieces) != 6 || len(model.ImageWall) < 7 {
		t.Fatalf("unexpected home model: %#v", model)
	}
}

func TestStudyMaterialsWorkflow(t *testing.T) {
	svc := New(catalog.New(), nil)
	model, err := svc.Study(context.Background(), StudyRequest{PieceID: "meihua", Difficulty: "高级"})
	if err != nil {
		t.Fatal(err)
	}
	if model.Bundle.Piece.Title != "梅花三弄" || len(model.Bundle.Notations) != 1 || !model.Bundle.Complete() {
		t.Fatalf("unexpected study model: %#v", model)
	}
}
