package subscriptions

import (
	"sync"

	"github.com/mike-jacks/neo/utils"
)

type EventType string

const (
	ObjectNodeCreated EventType = "objectNodeCreated"
	ObjectNodeUpdated EventType = "objectNodeUpdated"
	ObjectNodeDeleted EventType = "objectNodeDeleted"

	ObjectRelationshipCreated EventType = "objectRelationshipCreated"
	ObjectRelationshipUpdated EventType = "objectRelationshipUpdated"
	ObjectRelationshipDeleted EventType = "objectRelationshipDeleted"

	DomainSchemaNodeCreated EventType = "domainSchemaNodeCreated"
	DomainSchemaNodeUpdated EventType = "domainSchemaNodeUpdated"
	DomainSchemaNodeDeleted EventType = "domainSchemaNodeDeleted"

	TypeSchemaNodeCreated EventType = "typeSchemaNodeCreated"
	TypeSchemaNodeUpdated EventType = "typeSchemaNodeUpdated"
	TypeSchemaNodeDeleted EventType = "typeSchemaNodeDeleted"

	RelationshipSchemaNodeCreated EventType = "relationshipSchemaNodeCreated"
	RelationshipSchemaNodeUpdated EventType = "relationshipSchemaNodeUpdated"
	RelationshipSchemaNodeDeleted EventType = "relationshipSchemaNodeDeleted"
)

type Subscriber struct {
	ID     string
	Events chan interface{}
}

type SubscriptionManager struct {
	subscribers map[EventType]map[string]*Subscriber
	mu          sync.RWMutex
}

func NewSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{
		subscribers: make(map[EventType]map[string]*Subscriber),
	}
}

func (m *SubscriptionManager) Subscribe(eventType EventType) *Subscriber {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.subscribers[eventType] == nil {
		m.subscribers[eventType] = make(map[string]*Subscriber)
	}

	subscriber := &Subscriber{
		ID:     utils.GenerateId(),
		Events: make(chan interface{}, 1),
	}

	m.subscribers[eventType][subscriber.ID] = subscriber

	return subscriber
}

func (m *SubscriptionManager) Unsubscribe(eventType EventType, subscriberID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if subscribers, ok := m.subscribers[eventType]; ok {
		if subscriber, exists := subscribers[subscriberID]; exists {
			close(subscriber.Events)
			delete(subscribers, subscriberID)
		}
	}
}

func (m *SubscriptionManager) Publish(eventType EventType, data interface{}) {
	m.mu.RLock()
	subscribers := m.subscribers[eventType]
	m.mu.RUnlock()

	for _, subscriber := range subscribers {
		select {
		case subscriber.Events <- data:
		default:
			// Channel is full, skip this subscriber
		}
	}
}
