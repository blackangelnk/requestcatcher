package storage

import "github.com/blackangelnk/requestcatcher/internal/request"

type DummyStorage struct{}

func NewDummy() DummyStorage {
	return DummyStorage{}
}

func (s DummyStorage) Save(req request.CaughtRequest) (request.CaughtRequest, error) {
	return req, nil
}

func (s DummyStorage) Get() ([]request.CaughtRequest, error) {
	return make([]request.CaughtRequest, 0), nil
}
