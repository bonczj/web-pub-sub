package pubsub

import (
	"fmt"
	"log"
	"sync"
)

type Pubsub struct {
	Subscribers map[string]chan ([]byte)
	Lock        *sync.Mutex
}

var instance *Pubsub
var instanceLock = &sync.Mutex{}

// Instance retrieves a singleton instance of type
// Pubsub. This allows the publisher and subscribers
// to be able to share state so that the pubsub nature
// actually works.
func Instance() *Pubsub {
	instanceLock.Lock()
	defer instanceLock.Unlock()

	if instance == nil {
		instance = &Pubsub{
			Subscribers: map[string]chan ([]byte){},
			Lock:        &sync.Mutex{},
		}
	}

	return instance
}

// Publish determines all current subscribers and publishes the
// message to each of them.
func (p *Pubsub) Publish(message []byte) {
	for id, ch := range p.Subscribers {
		log.Printf("Sending message %s to subscriber ID %s", message, id)
		ch <- message
	}
}

// Subscribe will attempt to add a new subscriber id and their channel
// to the pubsub system. If the subscriber ID already exists, an error
// will be returned.
func (p *Pubsub) Subscribe(id string, channel chan ([]byte)) error {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	if _, exists := p.Subscribers[id]; exists {
		msg := fmt.Sprintf("Subscriber ID %s is around bound to a channel", id)
		log.Println(msg)
		return fmt.Errorf(msg)
	}

	p.Subscribers[id] = channel
	log.Printf("Added subscriber ID %s", id)
	return nil
}

// Unsubscribe checks if a given subscriber ID exists and if it
// does, closes the channel and removes it from the list of
// subscribers.
func (p *Pubsub) Unsubscribe(id string) {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	if ch, exists := p.Subscribers[id]; exists {
		close(ch)
		delete(p.Subscribers, id)
	}
}

// Clear closes any subscriber channels that might still be open and removes the subscriber
func (p *Pubsub) Clear() {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	for id, ch := range p.Subscribers {
		close(ch)
		delete(p.Subscribers, id)
	}
}
