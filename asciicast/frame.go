package asciicast

import (
	"encoding/json"
	"errors"
	"fmt"
)

type EventType string

const (
	UnknownEvent EventType = ""
	OutputEvent  EventType = "o"
	InputEvent   EventType = "i"
	ResizeEvent  EventType = "r"
	MarkerEvent  EventType = "m"
	ExitEvent    EventType = "x"
)

func decodeEvent(s string) (EventType, error) {
	switch s {
	case "o":
		return OutputEvent, nil
	case "i":
		return InputEvent, nil
	case "r":
		return ResizeEvent, nil
	case "m":
		return MarkerEvent, nil
	case "x":
		return ExitEvent, nil
	default:
		return UnknownEvent, errors.New("unknown event type")
	}
}

type Frame struct {
	Time      float64 // Delay
	EventType EventType
	EventData []byte //Data
}

func (f *Frame) MarshalJSON() ([]byte, error) {
	s, _ := json.Marshal(string(f.EventData))
	json := fmt.Sprintf(`[%.6f, "o", %s]`, f.Time, s)
	return []byte(json), nil
}

func (f *Frame) UnmarshalJSON(data []byte) error {
	var x interface{}

	err := json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	xx := x.([]interface{})
	if len(xx) == 3 {
		eventType, err := decodeEvent(xx[1].(string))
		if err != nil {
			return err
		}

		f.Time = xx[0].(float64)
		f.EventType = eventType
		s := []byte(xx[2].(string))
		b := make([]byte, len(s))
		copy(b, s)
		f.EventData = b
	} else if len(xx) == 2 {
		f.Time = xx[0].(float64)
		s := []byte(xx[1].(string))
		b := make([]byte, len(s))
		copy(b, s)
		f.EventData = b
		f.EventType = OutputEvent
	}
	return nil
}
