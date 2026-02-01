package scheduler

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Scheduler struct {
	logger   *zap.Logger
	stopChan chan struct{}
	wg       sync.WaitGroup
	running  bool
	mu       sync.Mutex
}

func New(logger *zap.Logger) *Scheduler {
	return &Scheduler{
		logger:   logger,
		stopChan: make(chan struct{}),
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	s.wg.Add(1)
	go s.run(ctx)

	s.logger.Info("scheduler started")
}

func (s *Scheduler) run(ctx context.Context) {
	defer s.wg.Done()

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("scheduler stopped due to context cancellation")
			return
		case <-s.stopChan:
			s.logger.Info("scheduler stopped")
			return
		case <-ticker.C:
			s.logger.Debug("scheduler tick - CT log polling placeholder")
		}
	}
}

func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}

	close(s.stopChan)
	s.wg.Wait()
	s.running = false
}
