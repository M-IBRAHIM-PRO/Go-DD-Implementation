package events

import (
    "log"
    "library_system/internal/domain"
)

// LogDispatcher — a simple implementation that logs events.
// In production this could publish to Kafka, RabbitMQ, etc.
// The domain doesn't care — it only knows the EventDispatcher interface.
type LogDispatcher struct{}

func NewLogDispatcher() *LogDispatcher {
    return &LogDispatcher{}
}

func (d *LogDispatcher) Dispatch(events []domain.DomainEvent) {
    for _, e := range events {
        log.Printf("[DOMAIN EVENT] %s at %s", e.EventName(), e.OccurredAt().Format("2006-01-02 15:04:05"))
    }
}