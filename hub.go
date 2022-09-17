package event

import (
	"fmt"
	"reflect"
	"sync"
)

type HubT struct {
	name   string
	mutex  sync.RWMutex
	topics map[string]any
	logger Logger
}

type Hub = *HubT

func NewHub(name string, logger Logger) Hub {
	return &HubT{
		name:   name,
		mutex:  sync.RWMutex{},
		topics: map[string]any{},
		logger: logger,
	}
}

func (me Hub) Name() string { return me.name }

func (me Hub) CreateTopic(name string, eventExample any, logger Logger) any {
	me.mutex.Lock()
	defer me.mutex.Unlock()

	_, has := me.topics[name]
	if has {
		panic(fmt.Errorf("topic %s already exists", name))
	}

	r := NewTopic(name, eventExample, logger)
	me.topics[name] = r

	return r
}

func (me Hub) HasTopic(name string) bool {
	me.mutex.RLock()
	defer me.mutex.RUnlock()

	_, has := me.topics[name]
	return has
}

func (me Hub) GetTopic(name string, eventExample any) any {
	me.mutex.RLock()
	defer me.mutex.RUnlock()

	r, has := me.topics[name]
	if has {
		expected := reflect.TypeOf(eventExample)
		actual := r.(Topic[any]).EventType()
		if expected != actual {
			panic(fmt.Errorf("expected topic event type is %v, however, the actual one is %v", expected, actual))
		}
	}

	return r
}

func GetTopic[K any](hub Hub, name string, eventExample K) Topic[K] {
	r := hub.GetTopic(name, eventExample)
	if r == nil {
		return nil
	}
	return r.(Topic[K])
}

func CreateTopic[K any](hub Hub, name string, eventExample any, logger Logger) Topic[K] {
	return hub.CreateTopic(name, eventExample, logger).(Topic[K])
}