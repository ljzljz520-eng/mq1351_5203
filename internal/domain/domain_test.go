package domain

import "testing"

func TestSubmissionNormalization(t *testing.T) {
	item, err := NormalizeSubmission(ExperienceSubmission{Name: "  琴友  ", Contact: " mail@example.com ", Interest: " 名曲聆听 ", Message: "  想听流水  "})
	if err != nil {
		t.Fatal(err)
	}
	if item.Name != "琴友" || item.Status != "pending" || item.Message != "想听流水" {
		t.Fatalf("unexpected normalized item: %#v", item)
	}
}

func TestStudyTextMatching(t *testing.T) {
	if !MatchStudyText("流水", "流水·泛音开篇", "泛音") {
		t.Fatal("expected match")
	}
	if MatchStudyText("不存在", "流水", "高山") {
		t.Fatal("unexpected match")
	}
}
