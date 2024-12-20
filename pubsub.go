package main

import (
	"fmt"
	"sync"
)

// PubSub represents the publish-subscribe system
type PubSub struct {
	mu         sync.RWMutex
	subscribers map[string][]chan Value // Map of channels to their subscribers
}

// NewPubSub creates a new PubSub instance
func NewPubSub() *PubSub {
	return &PubSub{
		subscribers: make(map[string][]chan Value),
	}
}

// Subscribe adds a client to a channel's list of subscribers
func (ps *PubSub) Subscribe(channel string) <-chan Value {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan Value, 10) // Buffered channel for messages
	ps.subscribers[channel] = append(ps.subscribers[channel], ch)
	return ch
}

// Unsubscribe removes a client from a channel's list of subscribers
func (ps *PubSub) Unsubscribe(channel string, sub <-chan Value) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if subs, ok := ps.subscribers[channel]; ok {
		for i, ch := range subs {
			if ch == sub {
				ps.subscribers[channel] = append(subs[:i], subs[i+1:]...)
				close(ch)
				break
			}
		}
		// If no subscribers remain, delete the channel entry
		if len(ps.subscribers[channel]) == 0 {
			delete(ps.subscribers, channel)
		}
	}
}

// Publish sends a message to all subscribers of a channel
func (ps *PubSub) Publish(channel string, message Value) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if subs, ok := ps.subscribers[channel]; ok {
		for _, ch := range subs {
			// Non-blocking send to prevent slow clients from halting the publisher
			select {
			case ch <- message:
			default:
				fmt.Println("Dropping message for slow subscriber")
			}
		}
	}
}
