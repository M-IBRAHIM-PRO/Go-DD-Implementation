package domain

import "time"

// DomainEvent is the base interface every event must satisfy.
// It's in the root domain package — shared across all bounded contexts.
type DomainEvent interface {
    EventName() string
    OccurredAt() time.Time
}

// EventDispatcher is an interface — domain doesn't care HOW events are fired.
// Could be in-memory, could be Kafka, domain doesn't know or care.
type EventDispatcher interface {
    Dispatch(events []DomainEvent)
}