package stress

import (
	"context"
	"sync"
	"time"

	"github.com/rcrowley/go-metrics"
	"github.com/sirupsen/logrus"
	"go.uber.org/multierr"
)

// Worker uses separate goroutine to do work with specified robots.
type Worker struct {
	index  int // worker index
	robots []*Rotbot
	tps    metrics.Meter
}

func (w *Worker) Init() error {
	for _, robot := range w.robots {
		err := robot.GetReady()
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *Worker) Add(r *Rotbot) {
	w.robots = append(w.robots, r)
}

// Repeat repeatly do work with the specified task.
func (w *Worker) Repeat(
	ctx context.Context, interval time.Duration, task Task, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for iteration := 0; ; iteration++ {
		select {
		case <-ctx.Done():
			logrus.WithFields(logrus.Fields{
				"worker":    w.index,
				"iteration": iteration,
			}).Info("Stress testing worker completed.")
			return
		case <-ticker.C:
			if err := w.do(task, iteration); err != nil {
				// Stop goroutine once error occurred.
				logrus.WithField("worker", w.index).
					WithError(err).
					Errorf("Failed to do job once at iteration = %v", iteration)
				continue
			}
		}
	}
}

func (w *Worker) do(task Task, iteration int) (err error) {
	for _, robot := range w.robots {
		e := task.Do(robot, iteration)
		if e != nil {
			err = multierr.Combine(err, e)
			continue
		}
		w.tps.Mark(1)
	}

	return err
}

func (w *Worker) Close() {
	for _, robot := range w.robots {
		robot.Close()
	}
}
