package music

import (
	"fmt"
	"sync"

	"qin-culture-site/internal/domain"
)

type Queue struct {
	mu      sync.Mutex
	items   []domain.QinPiece
	index   int
	loop    bool
	started bool
}

func NewQueue(items []domain.QinPiece) *Queue {
	return &Queue{items: append([]domain.QinPiece(nil), items...)}
}

func (q *Queue) Add(piece domain.QinPiece) error {
	if piece.ID == "" || piece.Title == "" {
		return fmt.Errorf("queue item requires id and title")
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, piece)
	return nil
}

func (q *Queue) Remove(id string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	for index, piece := range q.items {
		if piece.ID == id {
			q.items = append(q.items[:index], q.items[index+1:]...)
			if q.index >= len(q.items) {
				q.index = 0
			}
			return true
		}
	}
	return false
}

func (q *Queue) Current() (domain.QinPiece, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 || q.index < 0 || q.index >= len(q.items) {
		return domain.QinPiece{}, false
	}
	return q.items[q.index], true
}

func (q *Queue) Next() (domain.QinPiece, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return domain.QinPiece{}, false
	}
	if q.index+1 >= len(q.items) {
		if !q.loop {
			return q.items[q.index], false
		}
		q.index = 0
	} else {
		q.index++
	}
	q.started = true
	return q.items[q.index], true
}

func (q *Queue) Previous() (domain.QinPiece, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return domain.QinPiece{}, false
	}
	if q.index == 0 {
		if !q.loop {
			return q.items[q.index], false
		}
		q.index = len(q.items) - 1
	} else {
		q.index--
	}
	return q.items[q.index], true
}

func (q *Queue) SetLoop(loop bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.loop = loop
}

func (q *Queue) IsLooping() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.loop
}

func (q *Queue) Items() []domain.QinPiece {
	q.mu.Lock()
	defer q.mu.Unlock()
	return append([]domain.QinPiece(nil), q.items...)
}

func (q *Queue) Position() (int, int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.index, len(q.items)
}
