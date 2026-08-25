package catalog

import "testing"

func TestStudyFilterAndDifficulty(t *testing.T) {
	c := New()
	bundle := c.FilterStudy(StudyFilter{PieceID: "liushui", Query: "泛音"})
	if bundle.Piece.Title != "流水" || len(bundle.Notations) != 1 {
		t.Fatalf("unexpected filtered bundle: %#v", bundle)
	}
	if len(c.DifficultyOptions()) != 3 {
		t.Fatalf("unexpected difficulty options: %#v", c.DifficultyOptions())
	}
}
