package storage

import (
	models_ "Practice2/internal/models"
	"sync"
)

type MemoryStore struct {
	mu     sync.Mutex
	tasks  map[int]models_.Task
	nextID int
}

func NewStore() *MemoryStore {
	return &MemoryStore{
		tasks:  make(map[int]models_.Task),
		nextID: 1,
	}
}

func (s *MemoryStore) Create(title string, description string) models_.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := models_.Task{
		ID:          s.nextID,
		Title:       title,
		Description: description,
		Completed:   false,
	}

	s.tasks[s.nextID] = task
	s.nextID++

	return task
}

func (s *MemoryStore) GetAll() []models_.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	list := make([]models_.Task, 0, len(s.tasks))

	for _, t := range s.tasks {
		list = append(list, t)
	}

	return list
}

func (s *MemoryStore) Get(id int) (models_.Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, ok := s.tasks[id]
	return t, ok
}

func (s *MemoryStore) Update(id int, completed bool) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, ok := s.tasks[id]
	if !ok {
		return false
	}

	t.Completed = completed
	s.tasks[id] = t
	return true
}
