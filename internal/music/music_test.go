package music

import (
	"testing"
	"time"

	"qin-culture-site/internal/catalog"
)

func TestMusicTitleFollowsSelection(t *testing.T) {
	player := NewPlayer(catalog.New().Pieces())
	if _, err := player.Select("liushui"); err != nil {
		t.Fatal(err)
	}
	titles := make(chan string, 1)
	if err := player.Play(func(title string) { titles <- title }); err != nil {
		t.Fatal(err)
	}
	if _, err := player.Select("meihua"); err != nil {
		t.Fatal(err)
	}
	player.Finish()
	select {
	case title := <-titles:
		if title != "梅花三弄" {
			t.Fatalf("title should follow current selection, got %q", title)
		}
	case <-time.After(time.Second):
		t.Fatal("playback callback did not complete")
	}
}

func TestPlaylistSearch(t *testing.T) {
	playlist := NewPlaylist(catalog.New().Pieces())
	if got := playlist.Search("流水"); len(got) != 1 || got[0].ID != "liushui" {
		t.Fatalf("unexpected search result: %#v", got)
	}
	if _, ok := playlist.At(-1); ok {
		t.Fatal("negative index should be rejected")
	}
}
