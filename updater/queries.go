package updater

import (
	"database/sql"
	"time"

	"github.com/travelaudience/go-sx"
)

func getNextJob(db *sql.DB) (TopicUpdateJob, error) {
	var result TopicUpdateJob

	err := sx.Do(db, func(tx *sx.Tx) {
		tx.MustQueryRow(`UPDATE topic_update_job
SET locked = true
WHERE id = (SELECT id
            FROM topic_update_job
            WHERE run_after
                < $1
              AND ran_at IS NULL
              AND locked IS false
            ORDER BY priority, run_after
            LIMIT 1)
RETURNING id, created_at, topic_id, priority, run_after, retry_waits, ran_at, error, locked;`, time.Now()).
			MustScan(
				&result.ID,
				&result.CreatedAt,
				&result.TopicID,
				&result.Priority,
				&result.RunAfter,
				&result.RetryWaits,
				&result.RanAt,
				&result.Error,
				&result.Locked,
			)
	})

	return result, err
}

func unloadAll(db *sql.DB) error {
	err := sx.Do(db, func(tx *sx.Tx) {
		tx.MustExec(`UPDATE topic_update_job
SET locked = false
WHERE locked = true`)
	})
	return err
}

func enqueueTopicID(db *sql.DB, topicID int, priority int) (int64, error) {
	return enqueueTopicJob(db, topicID, priority, time.Now(), []time.Duration{
		time.Second * 60,
		time.Second * 60 * 10,
		time.Second * 60 * 30,
		time.Second * 60 * 60,
	})
}

func enqueueTopicJob(db *sql.DB, topicID int, priority int, runAfter time.Time, durations Durations) (int64, error) {
	job := TopicUpdateJob{
		TopicID:    topicID,
		Priority:   priority,
		RunAfter:   runAfter,
		RetryWaits: durations,
	}

	var jobID int64

	err := sx.Do(db, func(tx *sx.Tx) {
		tx.MustQueryRow(`INSERT INTO topic_update_job (topic_id,
                              priority,
                              run_after,
                              retry_waits)
VALUES ($1,
        $2,
        $3,
        $4)
ON CONFLICT (topic_id,run_after) DO UPDATE SET priority = $2
RETURNING id`, job.TopicID, job.Priority, job.RunAfter, job.RetryWaits).MustScan(&jobID)
	})

	return jobID, err
}

func updateJobFailed(db *sql.DB, jobID int64, ranAt sql.NullTime, jobError string) error {
	err := sx.Do(db, func(tx *sx.Tx) {
		tx.MustExec(`UPDATE topic_update_job
SET ran_at = $1,
    error  = $2,
    locked = false
WHERE id = $3`, ranAt, jobError, jobID)
	})
	return err
}

func updateTopicJob(db *sql.DB, jobID int64, afterTime time.Time, durations Durations) error {
	err := sx.Do(db, func(tx *sx.Tx) {
		tx.MustExec(`UPDATE topic_update_job
SET run_after   = $1,
    retry_waits = $2,
    locked      = false
WHERE id = $3`, afterTime, durations, jobID)
	})
	return err
}

func deleteJob(db *sql.DB, jobID int64) error {
	err := sx.Do(db, func(tx *sx.Tx) {
		tx.MustExec(`DELETE
FROM topic_update_job
WHERE id = $1`, jobID)
	})
	return err
}
