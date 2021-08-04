package storage

import (
	"testing"
	"time"

	"github.com/blackangelnk/requestcatcher/internal/config"
	"github.com/blackangelnk/requestcatcher/internal/request"
	"github.com/stretchr/testify/assert"
)

func TestMemInit(t *testing.T) {
	cfg := &config.MemStorageConfig{}
	s := NewMem(cfg)

	assert.NotNil(t, s)
	requests, err := s.Get()
	assert.NoError(t, err)
	assert.Len(t, requests, 0)
}

func TestMemSave(t *testing.T) {
	cfg := &config.MemStorageConfig{}
	s := NewMem(cfg)
	cr := request.CaughtRequest{
		Id: 123,
	}
	_, err := s.Save(cr)
	assert.NoError(t, err)
}

func TestMemGet(t *testing.T) {
	cfg := &config.MemStorageConfig{}
	s := NewMem(cfg)

	cr := request.CaughtRequest{
		Id:            1,
		Time:          request.Time(time.Now()),
		Method:        "GET",
		ContentLength: 0,
		RemoteAddr:    "127.0.0.1:123",
		Url:           "/",
		Body:          "",
	}

	cr, err := s.Save(cr)
	assert.NoError(t, err)
	requests, err := s.Get()
	assert.NoError(t, err)
	assert.Len(t, requests, 1)
	assert.Equal(t, cr.Id, requests[0].Id)
}
