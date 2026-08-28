package catalog

import "testing"

func TestCatalogHasCultureFixtures(t *testing.T) {
	c := New()
	if len(c.Schools()) < 4 || len(c.Pieces()) < 6 || len(c.Stories()) < 3 {
		t.Fatalf("fixtures are incomplete")
	}
	if len(c.ImageWall()) != len(c.Schools())+len(c.Stories()) {
		t.Fatal("image wall does not combine culture items")
	}
}
