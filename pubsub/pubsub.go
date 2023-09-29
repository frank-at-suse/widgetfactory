package pubsub

import "sync"

type StreamMessageKind string

const (
	StreamMessageKindCreate StreamMessageKind = "create"
	StreamMessageKindDelete StreamMessageKind = "delete"
	StreamMessageKindLoad   StreamMessageKind = "load"
	StreamMessageKindError  StreamMessageKind = "error"
)

type StreamMessage struct {
	Kind   StreamMessageKind
	Object any
}

type Agent struct {
	mu     sync.Mutex
	subs   map[string][]chan StreamMessage
	quit   chan struct{}
	closed bool
}

func New() *Agent {
	return &Agent{
		subs: make(map[string][]chan StreamMessage),
		quit: make(chan struct{}),
	}
}

func (a *Agent) Publish(topic string, message StreamMessage) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.closed {
		return
	}

	for _, ch := range a.subs[topic] {
		ch <- message
	}
}

func (a *Agent) Subscribe(topic string) <-chan StreamMessage {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.closed {
		return nil
	}

	ch := make(chan StreamMessage)
	a.subs[topic] = append(a.subs[topic], ch)

	return ch
}

func (a *Agent) Close() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.closed {
		return
	}

	a.closed = true
	close(a.quit)

	for _, ch := range a.subs {
		for _, sub := range ch {
			close(sub)
		}
	}
}
