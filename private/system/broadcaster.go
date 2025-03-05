package system

import (
	"sync"
)

// ClockType represents the type of clock cycle & edge.
type ClockType int

const (
	MRisingEdge ClockType = iota
	MFallingEdge
	TRisingEdge
	TFallingEdge
)

// Broadcaster handles subscription and broadcasting of both m-cycles and t-cycles.
type Broadcaster struct {
	mu    sync.Mutex
	subs  [4][]chan struct{}
	count int64
}

// Subscribe adds a new subscriber for the specified clock type and returns its channel.
func (b *Broadcaster) Subscribe(c ClockType) <-chan struct{} {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan struct{}, 1) // Buffered to prevent blocking
	b.subs[c] = append(b.subs[c], ch)
	return ch
}

func (b *Broadcaster) Broadcast(c ClockType) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, ch := range b.subs[c] {
		select {
		case ch <- struct{}{}:
			// Successfully sent cycle tick
		default:
			// Subscriber not ready; skip to prevent blocking
		}
	}
}

func (b *Broadcaster) TClock() {
	if b.count%4 == 0 {
		b.Broadcast(MRisingEdge)
		b.Broadcast(TRisingEdge)

		b.Broadcast(MFallingEdge)
		b.Broadcast(TFallingEdge)
	} else {
		b.Broadcast(TRisingEdge)
		b.Broadcast(TFallingEdge)
	}

	b.count++
}
