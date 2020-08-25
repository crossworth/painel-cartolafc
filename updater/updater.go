package updater

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/crossworth/cartola-web-admin/logger"
)

// this is based on: https://github.com/btubbs/pgq

type TopicUpdateJob struct {
	ID         int64
	CreatedAt  time.Time
	TopicID    int
	Priority   int
	RunAfter   time.Time
	RetryWaits Durations
	RanAt      time.Time
	Error      string
}

type WorkerFunc func(TopicUpdateJob) error

type Worker struct {
	running             bool
	fn                  WorkerFunc
	deleteJobOnComplete bool
}

type TopicUpdater struct {
	db                 *sql.DB
	workersLock        sync.Mutex
	workers            []Worker
	stopChan           chan bool
	jobPollingInterval time.Duration
}

func NewTopicUpdater(db *sql.DB) *TopicUpdater {
	t := &TopicUpdater{
		db:                 db,
		stopChan:           make(chan bool),
		jobPollingInterval: 10 * time.Millisecond,
	}

	return t
}

func (t *TopicUpdater) run() {
	for {
		select {
		case <-t.stopChan:
			close(t.stopChan)
			return
		default:
			t.workersLock.Lock()
			var worker *Worker

			for i := range t.workers {
				if t.workers[i].running {
					continue
				}

				t.workers[i].running = true
				worker = &t.workers[i]
				break
			}
			t.workersLock.Unlock()

			if worker != nil {
				go func(worker *Worker) {
					t.performNextJob(worker)
					t.workersLock.Lock()
					defer t.workersLock.Unlock()
					worker.running = false
				}(worker)
			}

			time.Sleep(t.jobPollingInterval)
		}
	}
}

func (t *TopicUpdater) performNextJob(worker *Worker) bool {
	tx, err := t.db.Begin()
	if err != nil {
		logger.Log.Warn().Err(err).Msg("erro ao iniciar uma transaction")
		return false
	}

	job, err := getNextJob(tx)
	if errors.Is(err, sql.ErrNoRows) {
		_ = tx.Rollback()
		return false
	}

	if err != nil {
		logger.Log.Warn().Err(err).Msg("erro ao conseguir próximo trabalho em TopicUpdater")
		_ = tx.Rollback()
		return false
	}

	if worker.fn == nil {
		logger.Log.Warn().Msg("função de worker nula em TopicUpdater")
		_ = tx.Rollback()
		return false
	}

	ranAt := time.Now()
	var jobErr error
	func() {
		defer func() {
			if r := recover(); r != nil {
				jobErr = fmt.Errorf("panic ao executar o job %d com tópico %d: Job: %v, Panic: %v", job.ID, job.TopicID, jobErr, r)
			}
		}()
		jobErr = worker.fn(job)
	}()

	if jobErr == nil && worker.deleteJobOnComplete {
		err = deleteJob(tx, job.ID)
		if err != nil {
			logger.Log.Warn().Err(err).Msgf("erro ao apagar job %d", job.ID)
		}
	} else {
		if len(job.RetryWaits) > 0 {
			afterTime := time.Now().Add(job.RetryWaits[0])
			err := updateTopicJob(tx, job.ID, afterTime, job.RetryWaits[1:])
			if err != nil {
				logger.Log.Warn().Err(err).Msgf("não foi possível atualizar o job %d", job.ID)
			}
		} else {
			err = updateJobFailed(tx, job.ID, ranAt, jobErr.Error())
			if err != nil {
				logger.Log.Warn().Err(err).Msgf("erro ao atualizar job %d", job.ID)
			}
		}
	}

	_ = tx.Commit()
	return jobErr != nil
}

func (t *TopicUpdater) StartProcessing() {
	t.stopChan = make(chan bool)
	go t.run()
}

func (t *TopicUpdater) StopProcessing() {
	t.stopChan <- true
}

func (t *TopicUpdater) RegisterWorker(workerFunc WorkerFunc, deleteJobOnComplete bool) {
	t.workersLock.Lock()
	defer t.workersLock.Unlock()
	t.workers = append(t.workers, Worker{
		running:             false,
		fn:                  workerFunc,
		deleteJobOnComplete: deleteJobOnComplete,
	})
}

func (t *TopicUpdater) EnqueueTopicID(topicID int) error {
	return enqueueTopicID(t.db, topicID, 10)
}

func (t *TopicUpdater) EnqueueTopicIDWithPriority(topicID int, priority int) error {
	return enqueueTopicID(t.db, topicID, priority)
}
