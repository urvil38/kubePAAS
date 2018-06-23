package util

import (
	"fmt"
	"sync"
	"time"
)

type spinner struct {
	mu       *sync.Mutex
	frames   []rune
	pos      int
	len      int
	active   bool
	stopChan chan struct{}
	text     string
}

var spin = `|/-\`

func (s *spinner) Next() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	r := s.frames[s.pos%s.len]
	s.pos++
	return string(r)
}

func NewSpinner(text string) *spinner {
	return &spinner{
		frames:   []rune(spin),
		pos:      0,
		len:      len(spin),
		active:   false,
		stopChan: make(chan struct{}, 1),
		text:     text,
		mu:       &sync.Mutex{},
	}
}

func (s *spinner) Start() {
	if s.active {
		return
	}

	s.active = true

	go func() {
		for {
			select {
			case <-s.stopChan:
				return
			default:
				fmt.Printf("\r%s %s", s.text, s.Next())
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

func (s *spinner) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.active {
		s.active = false
		s.stopChan <- struct{}{}
	}
}
