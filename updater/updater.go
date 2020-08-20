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
	RanAt      sql.NullTime
	Error      string
	Locked     bool
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

	err := unloadAll(t.db)
	if err != nil {
		logger.Log.Warn().Err(err).Msg("não foi possível remover a trava dos trabalhos")
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
	job, err := getNextJob(t.db)
	if errors.Is(err, sql.ErrNoRows) {
		return false
	}

	if err != nil {
		logger.Log.Warn().Err(err).Msg("erro ao conseguir próximo trabalho em TopicUpdater")
		return false // maybe database slow? returning false will make it sleep for a few seconds
	}

	if worker.fn == nil {
		logger.Log.Warn().Msg("função de worker nula em TopicUpdater")
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
		err = deleteJob(t.db, job.ID)
		if err != nil {
			logger.Log.Warn().Err(err).Msgf("erro ao apagar job %d", job.ID)
		}
	} else {
		if len(job.RetryWaits) > 0 {
			afterTime := time.Now().Add(job.RetryWaits[0])
			err := updateTopicJob(t.db, job.ID, afterTime, job.RetryWaits[1:])
			if err != nil {
				logger.Log.Warn().Err(err).Msgf("não foi possível atualizar o job %d", job.ID)
			}
		} else {
			err = updateJobFailed(t.db, job.ID, sql.NullTime{
				Time:  ranAt,
				Valid: true,
			}, jobErr.Error())
			if err != nil {
				logger.Log.Warn().Err(err).Msgf("erro ao atualizar job %d", job.ID)
			}
		}
	}

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

func (t *TopicUpdater) EnqueueTopicID(topicID int) (int64, error) {
	return enqueueTopicID(t.db, topicID, 10)
}

func (t *TopicUpdater) EnqueueTopicIDWithPriority(topicID int, priority int) (int64, error) {
	return enqueueTopicID(t.db, topicID, priority)
}
