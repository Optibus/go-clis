package sets

import (
	"errors"
	"fmt"
	"sync"
)

// Set represents a collection of unique strings
type Set struct {
	maps map[string]struct{}
}

type DoubleSet struct {
	list   *Set
	checks map[string]bool
	mu     sync.Mutex
}

func NewDoubleSet() *DoubleSet {
	return &DoubleSet{
		list:   NewSet(),
		checks: make(map[string]bool),
	}
}

// NewSet creates a new Set instance
func NewSet() *Set {
	return &Set{
		maps: make(map[string]struct{}),
	}
}

func (s *DoubleSet) Add(str string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.list.add(str)
	if !s.hasBeenChecked(str) {
		s.checks[str] = false
	}
}

func (s *DoubleSet) hasBeenChecked(str string) bool {
	return s.checks[str]
}

func (s *DoubleSet) HasBeenChecked(str string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.checks[str]
}

func (s *DoubleSet) GetNoneChecked() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	newSet := NewSet()
	for _, str := range s.list.get() {
		if !s.hasBeenChecked(str) {
			newSet.add(str)
		}
	}
	return newSet.get()
}

func (s *DoubleSet) Check(str string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.list.has(str) {
		return errors.New("cant check key not exists")
	}
	s.checks[str] = true
	return nil
}

func (s *Set) has(str string) bool {
	if _, ok := s.maps[str]; ok {
		return true
	}
	return false
}

// Add adds a string to the Set if it doesn't already exist
func (s *Set) add(str string) error {
	if s.has(str) {
		return errors.New("element already exists")
	}
	s.maps[str] = struct{}{}
	return nil
}

// Get retrieves all elements from the Set
func (s *Set) get() []string {
	var result []string
	for str := range s.maps {
		result = append(result, str)
	}
	return result
}

// String returns a string representation of the Set
func (s *Set) string() string {
	return fmt.Sprintf("%v", s.get())
}

func Unique(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
