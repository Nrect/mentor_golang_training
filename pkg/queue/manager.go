package queue

import (
	"context"
	"sync"
	"time"
)

// QueueManager интерфейс менеджера очередей
type QueueManager interface {
	Put(queueName string, message string)
	Get(ctx context.Context, queueName string, timeout time.Duration) (string, bool)
}

// queueManager управляет очередями сообщений
type queueManager struct {
	queues      map[string][]string
	subscribers map[string][]chan string
	mu          sync.RWMutex
}

// NewQueueManager создает новый менеджер очередей
func NewQueueManager() QueueManager {
	return &queueManager{
		queues:      make(map[string][]string),
		subscribers: make(map[string][]chan string),
	}
}

// Put добавляет сообщение в указанную очередь
func (qm *queueManager) Put(queueName string, message string) {
	qm.mu.Lock()
	defer qm.mu.Unlock()

	if subscribers, exists := qm.subscribers[queueName]; exists && len(subscribers) > 0 {
		subscriber := subscribers[0]
		qm.subscribers[queueName] = subscribers[1:]
		subscriber <- message
		close(subscriber)
		return
	}

	qm.queues[queueName] = append(qm.queues[queueName], message)
}

// Get извлекает первое сообщение из очереди FIFO или ожидает таймаут
func (qm *queueManager) Get(ctx context.Context, queueName string, timeout time.Duration) (string, bool) {
	qm.mu.Lock()

	if messages, exists := qm.queues[queueName]; exists && len(messages) > 0 {
		message := messages[0]
		qm.queues[queueName] = messages[1:]
		qm.mu.Unlock()
		return message, true
	}

	subscriber := make(chan string)
	qm.subscribers[queueName] = append(qm.subscribers[queueName], subscriber)
	qm.mu.Unlock()

	select {
	case message := <-subscriber:
		return message, true

	case <-ctx.Done():
		qm.removeSubscriber(queueName, subscriber)
		return "", false

	case <-time.After(timeout):
		qm.removeSubscriber(queueName, subscriber)
		return "", false
	}
}

// removeSubscriber удаляет подписчика из списка
func (qm *queueManager) removeSubscriber(queueName string, subscriber chan string) {
	qm.mu.Lock()
	defer qm.mu.Unlock()

	for i, sub := range qm.subscribers[queueName] {
		if sub == subscriber {
			qm.subscribers[queueName] = append(qm.subscribers[queueName][:i], qm.subscribers[queueName][i+1:]...)
			break
		}
	}
}
