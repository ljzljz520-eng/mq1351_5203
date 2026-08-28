package music

import (
	"errors"
	"sync"

	"qin-culture-site/internal/domain"
)

var ErrNoSelection = errors.New("no music selection")

type Player struct {
	mu       sync.Mutex
	pieces   []domain.QinPiece
	current  domain.QinPiece
	active   *playback
	sequence int
}

type playback struct {
	finish chan struct{}
	closed chan struct{}
}

func NewPlayer(pieces []domain.QinPiece) *Player {
	return &Player{pieces: append([]domain.QinPiece(nil), pieces...)}
}

func (p *Player) Select(id string) (domain.QinPiece, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	selected, err := domain.ValidatePieceSelection(p.pieces, id)
	if err != nil {
		return domain.QinPiece{}, err
	}
	p.current = selected
	p.sequence++
	return selected, nil
}

func (p *Player) Current() (domain.QinPiece, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.current.ID == "" {
		return domain.QinPiece{}, false
	}
	return p.current, true
}

func (p *Player) SelectionNumber() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.sequence
}

func (p *Player) Play(done func(string)) error {
	p.mu.Lock()
	if p.current.ID == "" {
		p.mu.Unlock()
		return ErrNoSelection
	}
	if p.active != nil {
		close(p.active.finish)
	}
	active := &playback{finish: make(chan struct{}), closed: make(chan struct{})}
	p.active = active
	title := p.current.Title
	p.mu.Unlock()
	go func() {
		defer close(active.closed)
		<-active.finish
		defer done(title)
	}()
	return nil
}

func (p *Player) Finish() {
	p.mu.Lock()
	active := p.active
	p.active = nil
	p.mu.Unlock()
	if active != nil {
		close(active.finish)
		<-active.closed
	}
}

func (p *Player) Stop() {
	p.Finish()
}

func (p *Player) IsPlaying() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.active != nil
}

func (p *Player) AvailableIDs() []string {
	p.mu.Lock()
	defer p.mu.Unlock()
	ids := make([]string, 0, len(p.pieces))
	for _, piece := range p.pieces {
		ids = append(ids, piece.ID)
	}
	return ids
}
