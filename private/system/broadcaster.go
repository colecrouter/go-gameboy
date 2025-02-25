package system

import (
	"sync"
	"time"
)

// Broadcaster handles subscription and broadcasting of both m-cycles and t-cycles.
type Broadcaster struct {
	mu      sync.Mutex
	mSubs   []chan struct{}
	tSubs   []chan struct{}
	closeCh chan struct{}
}

// NewBroadcaster initializes and returns a new Broadcaster.
func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		mSubs:   make([]chan struct{}, 0),
		tSubs:   make([]chan struct{}, 0),
		closeCh: make(chan struct{}),
	}
}

// SubscribeM adds a new m-cycle subscriber and returns its channel.
func (b *Broadcaster) SubscribeM() <-chan struct{} {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan struct{}, 1) // Buffered to prevent blocking
	b.mSubs = append(b.mSubs, ch)
	return ch
}

// SubscribeT adds a new t-cycle subscriber and returns its channel.
func (b *Broadcaster) SubscribeT() <-chan struct{} {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan struct{}, 1) // Buffered to prevent blocking
	b.tSubs = append(b.tSubs, ch)
	return ch
}

// BroadcastM sends an m-cycle tick to all m-cycle subscribers.
func (b *Broadcaster) BroadcastM() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, ch := range b.mSubs {
		select {
		case ch <- struct{}{}:
			// Successfully sent m-cycle tick
		default:
			// Subscriber not ready; skip to prevent blocking
		}
	}
}

// BroadcastT sends a t-cycle tick to all t-cycle subscribers.
func (b *Broadcaster) BroadcastT() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, ch := range b.tSubs {
		select {
		case ch <- struct{}{}:
			// Successfully sent t-cycle tick
		default:
			// Subscriber not ready; skip to prevent blocking
		}
	}
}

// Close closes all subscriber channels and stops the broadcaster.
func (b *Broadcaster) Close() {
	close(b.closeCh)
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, ch := range b.mSubs {
		close(ch)
	}
	for _, ch := range b.tSubs {
		close(ch)
	}
	b.mSubs = nil
	b.tSubs = nil
}

// ClockGenerator generates clock ticks at specified intervals.
// It emits 4 t-cycles per 1 m-cycle.
func ClockGenerator(b *Broadcaster, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	tCycleCount := 0
	for {
		select {
		case <-ticker.C:
			tCycleCount++
			if tCycleCount%4 == 0 {
				// On every 4th t-cycle, emit m-cycle tick first, then t-cycle
				b.BroadcastM()
				b.BroadcastT()
			} else {
				// Regular t-cycle tick
				b.BroadcastT()
			}
		case <-b.closeCh:
			return
		}
	}
}
