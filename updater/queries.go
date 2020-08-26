package updater

import (
	"database/sql"
	"time"

	"github.com/travelaudience/go-sx"
)

func getNextJob(tx *sql.Tx) (TopicUpdateJob, error) {
	var result TopicUpdateJob

	row := tx.QueryRow(`SELECT id,
       created_at,
       topic_id,
       priority,
       run_after,
       retry_waits,
       ran_at,
       error
FROM topic_update_job
WHERE run_after < $1
  AND ran_at = '0001-01-01 00:00:00'
ORDER BY priority, run_after
LIMIT 1 FOR UPDATE SKIP LOCKED;`, time.Now().UTC())

	err := row.Scan(
		&result.ID,
		&result.CreatedAt,
		&result.TopicID,
		&result.Priority,
		&result.RunAfter,
		&result.RetryWaits,
		&result.RanAt,
		&result.Error,
	)

	return result, err
}

func enqueueTopicID(db *sql.DB, topicID int, priority int) error {
	return enqueueTopicJob(db, topicID, priority, time.Now(), []time.Duration{
		time.Second * 60,
		time.Second * 60 * 10,
		time.Second * 60 * 30,
		time.Second * 60 * 60,
	})
}

func enqueueTopicSliceID(db *sql.DB, ids []int, priority int) error {
	return enqueueTopicSliceJob(db, ids, priority, time.Now(), []time.Duration{
		time.Second * 60,
		time.Second * 60 * 10,
		time.Second * 60 * 30,
		time.Second * 60 * 60,
	})
}

func enqueueTopicJob(db *sql.DB, topicID int, priority int, runAfter time.Time, durations Durations) error {
	job := TopicUpdateJob{
		TopicID:    topicID,
		Priority:   priority,
		RunAfter:   runAfter,
		RetryWaits: durations,
	}

	err := sx.Do(db, func(tx *sx.Tx) {
		tx.MustExec(`INSERT INTO topic_update_job (topic_id,
                              priority,
                              run_after,
                              retry_waits,
                              ran_at)
VALUES ($1,
        $2,
        $3,
        $4,
        $5)
ON CONFLICT DO NOTHING`, job.TopicID, job.Priority, job.RunAfter, job.RetryWaits, time.Time{})
	})

	return err
}

func enqueueTopicSliceJob(db *sql.DB, ids []int, priority int, runAfter time.Time, durations Durations) error {
	var jobs []TopicUpdateJob

	for _, id := range ids {
		job := TopicUpdateJob{
			TopicID:    id,
			Priority:   priority,
			RunAfter:   runAfter,
			RetryWaits: durations,
		}

		jobs = append(jobs, job)
	}

	err := sx.Do(db, func(tx *sx.Tx) {
		for _, job := range jobs {
			tx.MustExec(`INSERT INTO topic_update_job (topic_id,
                              priority,
                              run_after,
                              retry_waits,
                              ran_at)
VALUES ($1,
        $2,
        $3,
        $4,
        $5)
ON CONFLICT DO NOTHING`, job.TopicID, job.Priority, job.RunAfter, job.RetryWaits, time.Time{})
		}
	})

	return err
}

func updateJobFailed(tx *sql.Tx, jobID int64, ranAt time.Time, jobError string) error {
	_, err := tx.Exec(`UPDATE topic_update_job
SET ran_at = $1,
    error  = $2
WHERE id = $3`, ranAt, jobError, jobID)

	return err
}

func updateTopicJob(tx *sql.Tx, jobID int64, afterTime time.Time, durations Durations) error {
	_, err := tx.Exec(`UPDATE topic_update_job
SET run_after   = $1,
    retry_waits = $2
WHERE id = $3`, afterTime, durations, jobID)
	return err
}

func deleteJob(tx *sql.Tx, jobID int64) error {
	_, err := tx.Exec(`DELETE
FROM topic_update_job
WHERE id = $1`, jobID)
	return err
}
