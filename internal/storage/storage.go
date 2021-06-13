package storage

import "github.com/blackangelnk/requestcatcher/internal/request"

type Storage interface {
	Save(request.CaughtRequest) (request.CaughtRequest, error)
	Get() ([]request.CaughtRequest, error)
}
