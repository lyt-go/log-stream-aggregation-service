package subscription

import (
	"logaggregation/internal/eventchannel"
	"logaggregation/internal/substate"
)

type Manager struct {
	hub      *eventchannel.Hub
	registry *substate.Registry
	cached   map[string]*eventchannel.Subscription
}

func New() *Manager {
	return &Manager{
		hub:      eventchannel.New(),
		registry: substate.New(),
		cached:   make(map[string]*eventchannel.Subscription),
	}
}

func (m *Manager) Subscribe(topic string) *eventchannel.Subscription {
	if subscription := m.cached[topic]; subscription != nil {
		return subscription
	}
	if !m.registry.CanSubscribe(topic) {
		closed := make(chan string)
		close(closed)
		return &eventchannel.Subscription{Events: closed, Cancel: func() {}}
	}
	subscription := m.hub.Subscribe(topic)
	originalCancel := subscription.Cancel
	subscription.Cancel = func() {
		m.registry.MarkClosed(topic, subscription.ID)
		originalCancel()
	}
	m.cached[topic] = subscription
	return subscription
}

func (m *Manager) Publish(topic, event string) { m.hub.Publish(topic, event) }
