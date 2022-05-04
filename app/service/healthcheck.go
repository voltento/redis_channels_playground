package service

import "sync"

type HealthCheck interface {
	HealthCheck() error
}

type HealthChecker struct {
	services []HealthCheck
	wg       sync.WaitGroup
	err      chan error
}

func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		services: []HealthCheck{},
		err:      make(chan error, 1),
	}
}

func (h *HealthChecker) Healthy() error {

	for _, s := range h.services {
		h.wg.Add(1)
		go func(s HealthCheck) {
			if err := s.HealthCheck(); err != nil {
				h.err <- err
			} else {
				h.err <- nil
			}
			defer h.wg.Done()
		}(s)
	}

	h.wg.Wait()
	for i := 0; i < len(h.services); i += 1 {
		err := <-h.err
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *HealthChecker) RegisterService(s HealthCheck) {
	h.services = append(h.services, s)
}
