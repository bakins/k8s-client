package client

// Valid values for event types (new types could be added in future)
const (
	// Information only and will not cause any problems
	EventTypeNormal string = "Normal"
	// These events are to warn that something might go wrong
	EventTypeWarning string = "Warning"
)

type (

	// EventInterface is the methods needed for interacting with Events.
	EventInterface interface {
		CreateEvent(namespace string, item *Event) (*Event, error)
		GetEvent(namespace, name string) (result *Event, err error)
		ListEvents(namespace string, opts *ListOptions) (*EventList, error)
		WatchEvents(namespace string, opts *WatchOptions, events chan EventWatchEvent) error
		DeleteEvent(namespace, name string) error
		UpdateEvent(namespace string, item *Event) (*Event, error)
	}

	// EventWatchEvent is the methods needed for interacting with an Event watch
	EventWatchEvent interface {
		Type() WatchEventType
		Object() (*Event, error)
	}

	// EventSource contains information for an event.
	EventSource struct {
		Component string `json:"component,omitempty"`
		Host      string `json:"host,omitempty"`
	}

	//Event is a report of an event somewhere in the cluster.
	Event struct {
		TypeMeta       `json:",inline"`
		ObjectMeta     `json:"metadata,omitempty"`
		InvolvedObject ObjectReference `json:"involvedObject"`
		Reason         string          `json:"reason,omitempty"`
		Message        string          `json:"message,omitempty"`
		Source         EventSource     `json:"source,omitempty"`
		FirstTimestamp Time            `json:"firstTimestamp,omitempty"`
		LastTimestamp  Time            `json:"lastTimestamp,omitempty"`
		Count          int32           `json:"count,omitempty"`
		Type           string          `json:"type,omitempty"`
	}

	// EventList is a list of events..
	EventList struct {
		TypeMeta `json:",inline"`
		ListMeta `json:"metadata,omitempty"`
		// list of horizontal pod autoscaler objects.
		Items []Event `json:"items"`
	}
)

// NewEvent creates a new Event struct
func NewEvent(namespace, name string) *Event {
	return &Event{
		TypeMeta:   NewTypeMeta("Event", "v1"),
		ObjectMeta: NewObjectMeta(namespace, name),
	}
}
