package store

import "sync"

type Status string

const (
	Processing Status = "processing"
	Completed  Status = "completed"
)

type Response struct {
	StatusCode int
	Body       []byte
	Status     Status
}

type Store struct {
	mu    sync.Mutex
	store map[string]*Response
}

func NewStore() *Store {
	return &Store{
		store: make(map[string]*Response),
	}
}

func (s *Store) Get(key string) (*Response, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	resp, ok := s.store[key]
	return resp, ok
}

func (s *Store) Start(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.store[key]; exists {
		return false
	}

	s.store[key] = &Response{
		Status: Processing,
	}
	return true
}

func (s *Store) Save(key string, code int, body []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.store[key] = &Response{
		Status:     Completed,
		StatusCode: code,
		Body:       body,
	}
}
