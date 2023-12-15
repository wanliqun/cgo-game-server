package stress

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rcrowley/go-metrics"
	"github.com/sirupsen/logrus"
	"github.com/wanliqun/cgo-game-server/client"
)

type ScheduleOption struct {
	Workers         int
	RobotsPerWorker int

	Interval time.Duration
	Timeout  time.Duration

	ClientFactory func() *client.Client
}

type Schduler struct {
	option  ScheduleOption
	workers []*Worker
	tps     metrics.Meter
}

func NewScheduler(o ScheduleOption) *Schduler {
	var idx int
	var workers []*Worker
	tps := metrics.NewMeter()

	// Arrange robots to different workers
	for i := 0; i < o.Workers; i++ {
		worker := Worker{index: i, tps: tps}

		for j := 0; j < o.RobotsPerWorker; j++ {
			worker.Add(&Rotbot{
				name:   fmt.Sprintf("robot%d", idx),
				client: o.ClientFactory(),
			})

			idx++
		}

		workers = append(workers, &worker)
	}

	return &Schduler{option: o, workers: workers, tps: tps}
}

func (s *Schduler) ProduceTPS(task Task) {
	logrus.Info("Begin to produce TPS")

	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), s.option.Timeout)
	defer cancel()

	for _, w := range s.workers {
		defer w.Close()
		if err := w.Init(); err != nil {
			logrus.WithError(err).Errorf("Failed to Init worker %d", w.index)
			return
		}

		wg.Add(1)
		go w.Repeat(ctx, s.option.Interval, task, &wg)
	}

	wg.Add(1)
	go s.monitorTPS(ctx, &wg)

	wg.Wait()
}

func (s *Schduler) monitorTPS(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()

	start := time.Now()
	for {
		select {
		case <-ctx.Done():
			s.log(start)
			return
		case <-ticker.C:
			s.log(start)
		}
	}
}

func (s *Schduler) log(start time.Time) {
	logrus.Infof(
		"TPS: avg = %.2f, m1 = %.2f, m5 = %.2f, m15 = %.2f, elapsed = %v",
		s.tps.RateMean(), s.tps.Rate1(), s.tps.Rate5(), s.tps.Rate15(), time.Since(start),
	)
}
