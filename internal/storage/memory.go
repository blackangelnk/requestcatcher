package storage

import (
	"sync"

	"github.com/blackangelnk/requestcatcher/internal/config"
	"github.com/blackangelnk/requestcatcher/internal/request"
)

type memStorage struct {
	s []request.CaughtRequest
	m sync.RWMutex
}

func NewMem(cfg *config.MemStorageConfig) *memStorage {
	return &memStorage{
		s: make([]request.CaughtRequest, 0, cfg.Cap),
	}
}

func (s *memStorage) Save(req request.CaughtRequest) (request.CaughtRequest, error) {
	s.m.Lock()
	s.s = append(s.s, request.CaughtRequest{})
	copy(s.s[1:], s.s)
	s.s[0] = req
	s.m.Unlock()
	return req, nil
}

func (s *memStorage) Get() ([]request.CaughtRequest, error) {
	s.m.RLock()
	defer s.m.RUnlock()
	return s.s, nil
}
