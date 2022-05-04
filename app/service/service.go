package service

import (
	"github.com/pieterclaerhout/go-log"
	"redis_channels_playground/app/dto"
	"redis_channels_playground/app/repo"
)

type Service struct {
	repo  *repo.Redis
	repo2 *repo.Redis
}

func NewService(r *repo.Redis, r2 *repo.Redis) *Service {
	go func() {
		for err := r.Pull(); err == nil; err = r.Pull() {
		}
	}()

	go func() {
		for err := r2.Pull(); err == nil; err = r2.Pull() {
		}
	}()

	return &Service{repo: r, repo2: r}
}

func (r *Service) PostMessage() {
	err := r.repo.Push(&dto.Message{
		Name:    "Andre",
		Age:     1,
		Message: "hello world",
	})
	if err != nil {
		log.Errorf("error %v", err)
		return
	}
}
