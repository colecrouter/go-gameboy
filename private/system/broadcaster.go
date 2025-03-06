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

type Broadcaster struct {
	mu      sync.Mutex
	subs    [4][]chan struct{}
	ackSubs [4][]chan struct{}
	count   int64
}

// Subscribe adds a new subscriber for the specified clock type.
func (b *Broadcaster) Subscribe(c ClockType) (tick <-chan struct{}, ack chan<- struct{}) {
	b.mu.Lock()
	defer b.mu.Unlock()
	ticker := make(chan struct{}, 1)
	ackCh := make(chan struct{}, 1)
	b.subs[c] = append(b.subs[c], ticker)
	b.ackSubs[c] = append(b.ackSubs[c], ackCh)
	return ticker, ackCh
}

func (b *Broadcaster) broadcastWithAck(c ClockType) {
	b.mu.Lock()
	subs := b.subs[c]
	ackSubs := b.ackSubs[c]
	b.mu.Unlock()

	var wg sync.WaitGroup
	wg.Add(len(ackSubs))
	// Send tick to subscribers.
	for _, ch := range subs {
		select {
		case ch <- struct{}{}:
		default:
			// Skip if not ready.
		}
	}

	// Wait for all subscriber acknowledgements.
	for _, ackCh := range ackSubs {
		go func(ch chan struct{}) {
			<-ch // Wait for ack from subscriber.
			wg.Done()
		}(ackCh)
	}
	wg.Wait()
}

func (b *Broadcaster) TClock() {
	if b.count%4 == 0 {
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			b.broadcastWithAck(MRisingEdge)
			wg.Done()
		}()
		go func() {
			b.broadcastWithAck(TRisingEdge)
			wg.Done()
		}()
		wg.Wait()

		wg.Add(2)
		go func() {
			b.broadcastWithAck(MFallingEdge)
			wg.Done()
		}()
		go func() {
			b.broadcastWithAck(TFallingEdge)
			wg.Done()
		}()
	} else {
		b.broadcastWithAck(TRisingEdge)
		b.broadcastWithAck(TFallingEdge)
	}
	b.count++
}
