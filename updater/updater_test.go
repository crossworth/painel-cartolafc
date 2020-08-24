package updater

import (
	"database/sql"
	"log"
	"sync"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func TestNewTopicUpdater(t *testing.T) {
	db, err := sql.Open("postgres", "postgres://cartola:root@192.168.0.65/cartola?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		t.Fatal(err)
	}

	updater := NewTopicUpdater(db)
	updater.RegisterWorker(func(job TopicUpdateJob) error {
		log.Printf("1: Executando job: %d", job.ID)
		time.Sleep(10 * time.Second)
		return nil
	}, true)
	updater.RegisterWorker(func(job TopicUpdateJob) error {
		log.Printf("2: Executando job: %d", job.ID)
		time.Sleep(10 * time.Second)
		return nil
	}, true)

	updater.StartProcessing()
	for i := 0; i < 1000; i++ {
		err = updater.EnqueueTopicID(i)
		if err != nil {
			log.Println(err)
		}
	}

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
