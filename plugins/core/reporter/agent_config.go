package reporter

type AgentConfigEventType int32

const (
	MODIFY AgentConfigEventType = iota
	DELETED
)

type AgentConfigChangeWatcher interface {
	Key() string
	Notify(eventType AgentConfigEventType, newValue string)
	Value() string
}
