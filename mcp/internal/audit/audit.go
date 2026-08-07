package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Event struct {
	Timestamp  time.Time      `json:"timestamp"`
	Operation  string         `json:"operation"`
	Repository string         `json:"repository,omitempty"`
	Target     string         `json:"target,omitempty"`
	Outcome    string         `json:"outcome"`
	Metadata   map[string]any `json:"metadata,omitempty"`
}

type Logger struct {
	path string
	mu   sync.Mutex
}

func New(path string) *Logger { return &Logger{path: path} }

func (l *Logger) Record(event Event) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	event.Timestamp = time.Now().UTC()
	if err := os.MkdirAll(filepath.Dir(l.path), 0o700); err != nil {
		return fmt.Errorf("create audit directory: %w", err)
	}
	f, err := os.OpenFile(l.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("open audit log: %w", err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(event); err != nil {
		return fmt.Errorf("write audit event: %w", err)
	}
	return f.Sync()
}
