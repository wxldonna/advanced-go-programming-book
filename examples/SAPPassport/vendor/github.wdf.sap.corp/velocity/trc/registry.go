package trc

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
)

/* default registry */

// defReg is the trace topic registry (maintaining global and local trace levels for all known topics)
var defReg = newRegistry()

// InitTraceTopic creates and registrs a trace topic.
// The function should be called in package initialization phase only!
// InitTraceTopic will panic in case of
// (a) name conflicts with the name of an already registered topics
// (b) serialization errors for the given context
func InitTraceTopic(name string, description string, info ...Info) Tracer {
	return defReg.initTraceTopic(name, description, info...)
}

// ReconfigGlobal re-sets the global trace level.
func ReconfigGlobal(globalLevel Level) error {
	return defReg.reconfigGlobal(globalLevel)
}

// ReconfigTopic re-sets the trace level of the trace topic with the given name.
// If no such trace topic is registered, an error is returned
func ReconfigTopic(topicName string, localLevel Level) error {
	return defReg.reconfigTopic(topicName, localLevel)
}

func ReconfigFromString(cfg string) error {
	return defRegFlag.Set(cfg)
}

// GetGlobalLevel returns the current global trace level.
//
// PLEASE NOTE: This function must not be used for runtime checks
// during trace operations.
// It is meant to be used for trace level configuration related
// operations (e.g. deliver current value to config UI).
func GetGlobalLevel() Level {
	return defReg.getGlobalLevel()
}

// TopicSnapshot is a map of registered topics at a given point in time
type TopicsSnapshot map[string]interface {
	Description() string
	Level() Level
}

func (tp TopicsSnapshot) String() string {
	var b bytes.Buffer
	b.WriteString("TopicsSnapshot[")
	first := true
	for k, v := range tp {
		if first {
			first = false
		} else {
			b.WriteRune(',')
		}
		fmt.Fprintf(&b, "%q:%s", k, v)
	}
	b.WriteString("]")
	return b.String()
}

// GetTopics returns a snapshot of the currently registered topics
func GetTopics() TopicsSnapshot {
	return defReg.getTopics()
}

/* registry implementation */

// Create new registry with default global default level info.
func newRegistry() *registry {
	return &registry{globalLevel: LevelInfo, topicsMap: make(topicsMap)}
}

type registry struct {
	sync.Mutex        // Mutex is used as lock for consistent registry operations (less frequent); on topic level atomics are used for synchronized access
	globalLevel Level // globalLevel tracks the global trace level. For trace performance reasons the value is redundanty contained in each topic
	topicsMap         // topicsMap manages the trace levels for all topics known to this registry
}
type topicsMap map[string]*traceTopic // topicsMap maps: topic name (string) -> topic (*traceTopic)

// InitTraceTopic should be called in package initialization phase only!
// InitTraceTopic will panic in case of
// (a) name conflicts with the name of an already registered topics
// (b) serialization errors for the given context
func (rt *registry) initTraceTopic(name string, description string, info ...Info) Tracer {
	topic := newTraceTopic(name, description)
	if err := rt.registerTopic(topic); err != nil {
		panic(err)
	}
	return newTopicTracer(topic, info...)
}

// registerTopic adds the given topic.
// In case a topic with the same name is already registered, an error is returned
func (rt *registry) registerTopic(topic *traceTopic) error {
	rt.Lock()
	defer rt.Unlock()
	_, alreadyRegistered := rt.topicsMap[topic.Name]
	if alreadyRegistered {
		return errors.New("Trace topic registration failed: topic with the same name already registered: " + string(topic.Name))
	}
	topic.setGlobal(rt.globalLevel)
	rt.topicsMap[topic.Name] = topic
	return nil
}

// reconfigTopic re-sets the trace level of the trace topic with the given name.
// If no such trace topic is registered, an error is returned
func (rt *registry) reconfigTopic(topicName string, level Level) error {
	rt.Lock()
	defer rt.Unlock()
	t, ok := rt.topicsMap[topicName]
	if ok {
		t.setLocal(level)
		return nil
	}
	return fmt.Errorf("Unable to configure trace topic. Topic not found: %s", topicName)
}

// reconfigGlobal re-sets the global trace level.
// The value is stored and propagated into each registered trace topic
// for fast access
func (rt *registry) reconfigGlobal(level Level) error {
	rt.Lock()
	defer rt.Unlock()
	rt.globalLevel = level
	for _, t := range rt.topicsMap {
		t.setGlobal(level)
	}
	return nil
}

// getGlobalLevel returns the current global trace level.
//
// PLEASE NOTE: This function must not be used for runtime checks
// during trace operations.
// It is meant to be used for trace level configuration related
// operations (e.g. deliver current value to config UI).
func (rt *registry) getGlobalLevel() Level {
	rt.Lock()
	defer rt.Unlock()
	return rt.globalLevel
}

// snapshotTopic implements anonymous topic interface in TopicSnapshot map
type snapshotTopic struct {
	description string
	level       Level
}

func (st *snapshotTopic) Description() string {
	return st.description
}

func (st *snapshotTopic) Level() Level {
	return st.level
}

func (st *snapshotTopic) String() string {
	return fmt.Sprintf("{Description:%q,Level:%q}", st.Description(), st.Level())
}

// getTopics returns a snapshot of the currently registered topics
func (rt *registry) getTopics() TopicsSnapshot {
	rt.Lock()
	defer rt.Unlock()
	snapshot := make(TopicsSnapshot)
	for topicName, topic := range rt.topicsMap {
		snapshot[topicName] = &snapshotTopic{
			description: topic.Description,
			level:       topic.local(),
		}
	}
	return snapshot
}
