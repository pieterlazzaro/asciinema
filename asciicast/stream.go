package asciicast

import (
	"errors"
	"sync"
	"time"
)

type Stream struct {
	Frames        []Frame
	elapsedTime   time.Duration
	lastWriteTime time.Time
	maxWait       time.Duration
	lock          *sync.Mutex
}

type StreamWriter struct {
	*Stream
	t EventType
}

func (s *Stream) NewEventStreamWriter(t EventType) *StreamWriter {
	return &StreamWriter{
		Stream: s,
		t:      t,
	}
}

func (s *StreamWriter) Write(p []byte) (int, error) {
	if s.Stream == nil {
		return 0, errors.New("invalid stream")
	}

	return s.WriteFrame(s.t, p), nil
}

func NewStream(maxWait float64) *Stream {
	if maxWait <= 0 {
		maxWait = 1.0
	}
	return &Stream{
		lastWriteTime: time.Now(),
		maxWait:       time.Duration(maxWait*1000000) * time.Microsecond,
		lock:          &sync.Mutex{},
	}
}

func (s *Stream) WriteFrame(t EventType, p []byte) int {
	frame := Frame{}
	frame.EventType = t
	frame.Time = s.incrementElapsedTime().Seconds()
	frame.EventData = make([]byte, len(p))
	n := copy(frame.EventData, p)

	s.Frames = append(s.Frames, frame)

	return n
}

func (s *Stream) Write(p []byte) (int, error) {

	return s.WriteFrame(OutputEvent, p), nil
}

func (s *Stream) Close() {
	s.incrementElapsedTime()

	if len(s.Frames) > 0 && string(s.Frames[len(s.Frames)-1].EventData) == "exit\r\n" {
		s.Frames = s.Frames[:len(s.Frames)-1]
	}
}

func (s *Stream) Duration() time.Duration {
	return s.elapsedTime
}

func (s *Stream) incrementElapsedTime() time.Duration {
	s.lock.Lock()
	now := time.Now()
	d := now.Sub(s.lastWriteTime)

	if s.maxWait > 0 && d > s.maxWait {
		d = s.maxWait
	}
	if d <= 0 {
		d = time.Millisecond * 500
	}

	s.elapsedTime += d
	s.lastWriteTime = now
	s.lock.Unlock()
	return s.elapsedTime
}
