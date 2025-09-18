package memory

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/bstncartwright/beyond-database-sql-driver-pattern/03/event/driver"
)

var store sync.Map

type conn struct{}

func (c *conn) Push(ctx context.Context, topic driver.Topic, message driver.Message) error {
	val, _ := store.LoadOrStore(topic.Name, make(chan driver.Message))
	ch, ok := val.(chan driver.Message)
	if !ok {
		return fmt.Errorf("map value for key %q was not of correct type", topic.Name)
	}
	ch <- message
	return nil
}

func (c *conn) Subscribe(ctx context.Context, topics []driver.Topic) error {
	for _, topic := range topics {
		val, _ := store.LoadOrStore(topic.Name, make(chan driver.Message))
		ch, ok := val.(chan driver.Message)
		if !ok {
			return fmt.Errorf("map value for key %q was not of correct type", topic.Name)
		}

		go func(t driver.Topic, b chan driver.Message) {
			for msg := range b {
				if err := t.Consumer(msg); err != nil {
					log.Printf("Error consuming on topic %q: %s", t.Name, err)
				}
			}
		}(topic, ch)
	}

	<-ctx.Done()
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return nil
}
