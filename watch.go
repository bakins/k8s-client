package client

import (
	"encoding/json"
	"fmt"
)

const (
	// WatchEvent.Type constants
	WatchEventTypeAdded    WatchEventType = "ADDED"
	WatchEventTypeModified WatchEventType = "MODIFIED"
	WatchEventTypeDeleted  WatchEventType = "DELETED"
	WatchEventTypeError    WatchEventType = "ERROR"
)

type (
	WatchEventType string

	// WatchEvent objects are streamed from the api server in response to a watch request.
	// These are not API objects and may not be changed in a backward-incompatible way.
	WatchEvent struct {
		// the type of watch event; may be ADDED, MODIFIED, DELETED, or ERROR
		Type WatchEventType `json:"type,omitempty"`
		// For added or modified objects, this is the new object; for deleted objects,
		// it's the state of the object immediately prior to its deletion.
		// For errors, it's an api.Status.
		Object json.RawMessage `json:"object,omitempty"`
	}
)

// UnmarshalObject tries to unmarshal the Object field of the given event into the given object.
func (e *WatchEvent) UnmarshalObject(object interface{}) error {
	if len(e.Object) == 0 {
		return fmt.Errorf("Object is empty in event of type '%s'", e.Type)
	}
	return json.Unmarshal(e.Object, object)
}
