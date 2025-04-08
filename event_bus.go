package gv

import (
	"sync"
)

// EventBus struct
type EventBus struct {
	subscribers map[string][]chan string
	mu          sync.RWMutex
}

// NewEventBus creates a new event bus
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]chan string),
	}
}

// Subscribe to an event
func (eb *EventBus) Subscribe(event string) chan string {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	ch := make(chan string)
	eb.subscribers[event] = append(eb.subscribers[event], ch)
	return ch
}

// Publish an event
func (eb *EventBus) Publish(event string, data string) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	for _, ch := range eb.subscribers[event] {
		go func(c chan string) {
			c <- data
		}(ch)
	}
}

// // Example usage
// func main() {
// 	bus := NewEventBus()

// 	// Component A subscribes
// 	eventChan := bus.Subscribe("user:created")

// 	// Component B listens in a goroutine
// 	go func() {
// 		for data := range eventChan {
// 			fmt.Println("Component A received event:", data)
// 		}
// 	}()

// 	// Component C dispatches an event
// 	bus.Publish("user:created", "New user added")

// 	// Prevents the main function from exiting immediately
// 	select {}
// }
