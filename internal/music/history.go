package music

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"qin-culture-site/internal/domain"
)

type PlayEvent struct {
	Sequence int
	PieceID  string
	Title    string
	Action   string
	At       string
}

type History struct {
	mu    sync.Mutex
	items []PlayEvent
	clock func() time.Time
	limit int
}

func NewHistory(limit int, clock func() time.Time) *History {
	if limit < 1 {
		limit = 20
	}
	if clock == nil {
		clock = func() time.Time { return time.Unix(0, 0).UTC() }
	}
	return &History{clock: clock, limit: limit}
}

func (h *History) Record(piece domain.QinPiece, action string) PlayEvent {
	h.mu.Lock()
	defer h.mu.Unlock()
	event := PlayEvent{Sequence: len(h.items) + 1, PieceID: piece.ID, Title: piece.Title, Action: action, At: h.clock().UTC().Format(time.RFC3339)}
	h.items = append(h.items, event)
	if len(h.items) > h.limit {
		h.items = append([]PlayEvent(nil), h.items[len(h.items)-h.limit:]...)
	}
	return event
}

func (h *History) Items() []PlayEvent {
	h.mu.Lock()
	defer h.mu.Unlock()
	return append([]PlayEvent(nil), h.items...)
}

func (h *History) Latest() (PlayEvent, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.items) == 0 {
		return PlayEvent{}, false
	}
	return h.items[len(h.items)-1], true
}

func (h *History) CountFor(pieceID string) int {
	h.mu.Lock()
	defer h.mu.Unlock()
	count := 0
	for _, event := range h.items {
		if event.PieceID == pieceID && event.Action == "finish" {
			count++
		}
	}
	return count
}

func FormatEvent(event PlayEvent) string {
	return fmt.Sprintf("%s · %s · %s", event.Title, event.Action, event.At)
}

func SortEvents(items []PlayEvent) []PlayEvent {
	result := append([]PlayEvent(nil), items...)
	sort.SliceStable(result, func(i, j int) bool { return result[i].Sequence < result[j].Sequence })
	return result
}
