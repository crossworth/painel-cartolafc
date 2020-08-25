package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"github.com/goccy/go-yaml"
	vkapi "github.com/himidori/golang-vk-api"
	_ "github.com/lib/pq"
	"github.com/travelaudience/go-sx"

	"github.com/crossworth/cartola-web-admin/updater"
)

type Account struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
}

type Settings struct {
	DBUser     string    `yaml:"db_user"`
	DBPassword string    `yaml:"db_password"`
	DBHost     string    `yaml:"db_host"`
	DBDatabase string    `yaml:"db_database"`
	GroupID    int       `yaml:"group_id"`
	Accounts   []Account `yaml:"accounts"`
}

const settingsFile = "worker_settings.yml"

func main() {
	var settings Settings

	output, err := ioutil.ReadFile(settingsFile)
	if err != nil {
		log.Fatalf("não foi possível abrir o arquivo de configurações %v", err)
	}

	err = yaml.Unmarshal(output, &settings)
	if err != nil {
		log.Fatalf("não foi possível decodicar o arquivo de configurações %v", err)
	}

	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		settings.DBUser,
		settings.DBPassword,
		settings.DBHost,
		settings.DBDatabase,
	))
	if err != nil {
		log.Fatalf("não foi possível criar a conexão com o banco de dados")
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("não foi possível conectar no banco de dados")
	}

	if settings.GroupID == 0 {
		log.Fatalf("você deve informar o ID do grupo")
	}

	if settings.Accounts == nil || len(settings.Accounts) == 0 {
		log.Fatalf("você deve informar pelo menos 1 conta")
	}

	for i, account := range settings.Accounts {
		if account.Email == "" {
			log.Fatalf("você deve informar o email para a conta %d", i)
		}

		if account.Password == "" {
			log.Fatalf("você deve informar a senha para a conta %d", i)
		}
	}

	var accountsPoll []*vkapi.VKClient

	for _, account := range settings.Accounts {
		androidClient, err := vkapi.NewVKClient(vkapi.DeviceAndroid, account.Email, account.Password)
		if err != nil {
			log.Fatalf("erro ao criar o cliente Android para conta %s, %v", account.Email, err)
		}

		iPhoneClient, err := vkapi.NewVKClient(vkapi.DeviceIPhone, account.Email, account.Password)
		if err != nil {
			log.Fatalf("erro ao criar o cliente iPhone para conta %s, %v", account.Email, err)
		}

		WPhoneClient, err := vkapi.NewVKClient(vkapi.DeviceWPhone, account.Email, account.Password)
		if err != nil {
			log.Fatalf("erro ao criar o cliente WindowsPhone para conta %s, %v", account.Email, err)
		}

		accountsPoll = append(accountsPoll, androidClient, iPhoneClient, WPhoneClient)
	}

	topicUpdater := updater.NewTopicUpdater(db)

	var wg sync.WaitGroup
	wg.Add(len(accountsPoll))

	for i, account := range accountsPoll {
		topicUpdater.RegisterWorker(work(db, settings.GroupID, i, account), true)
	}
	log.Printf("adicionado %d workers\n", len(accountsPoll))

	topicUpdater.StartProcessing()

	wg.Wait()
}

func insertProfiles(db *sql.DB, profiles []Profile) error {
	mapProfiles := make(map[int]Profile)

	for _, profile := range profiles {
		mapProfiles[profile.ID] = profile
	}

	var profilesUnique []Profile

	for _, profile := range mapProfiles {
		profilesUnique = append(profilesUnique, profile)
	}

	return sx.Do(db, func(tx *sx.Tx) {
		profileQuery := `INSERT INTO profiles VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO UPDATE SET first_name = $2, last_name = $3, screen_name = $4, photo = $5`

		for _, profile := range profilesUnique {
			tx.MustExec(profileQuery, profile.ID, profile.FirstName, profile.LastName, profile.ScreenName, profile.Photo)
		}
	})
}

func insertTopic(db *sql.DB, topic Topic) error {
	return sx.Do(db, func(tx *sx.Tx) {
		topicQuery := `INSERT INTO topics VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT (id) DO UPDATE SET title = $2, is_closed = $3, is_fixed = $4, created_at = $5, updated_at = $6, created_by = $7, updated_by = $8`
		tx.MustExec(topicQuery, topic.ID, topic.Title, topic.IsClosed, topic.IsFixed, topic.CreatedAt, topic.UpdatedAt, topic.CreatedBy.ID, topic.UpdatedBy.ID, false)
	})
}

func insertComments(db *sql.DB, topicID int, comments []Comment) error {
	return sx.Do(db, func(tx *sx.Tx) {
		commentQuery := `INSERT INTO comments VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT (id) DO UPDATE SET date = $3, text = $4, likes = $5, reply_to_uid = $6, reply_to_cid = $7`
		attachmentQuery := `INSERT INTO attachments  VALUES($1, $2) ON CONFLICT (comment_id, content) DO UPDATE SET content = $1, comment_id = $2`

		for _, comment := range comments {
			tx.MustExec(commentQuery, comment.ID, comment.FromID, comment.Date, comment.Text, comment.Likes, comment.ReplyToUID, comment.ReplyToCID, topicID, comment.FromID)

			for _, attachment := range comment.Attachments {
				tx.MustExec(attachmentQuery, attachment, comment.ID)
			}
		}
	})
}

func insertPoll(db *sql.DB, topicID int, poll Poll) error {
	return sx.Do(db, func(tx *sx.Tx) {
		pollQuery := `INSERT INTO polls VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (id) DO UPDATE SET question = $2, votes = $3, multiple = $4, end_date = $5, closed = $6`
		tx.MustExec(pollQuery, poll.ID, poll.Question, poll.Votes, poll.Multiple, poll.EndDate, poll.Closed, topicID)

		pollAnswerQuery := `INSERT INTO poll_answers VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO UPDATE SET text = $2, votes = $3, rate = $4`
		for _, answer := range poll.Answers {
			tx.MustExec(pollAnswerQuery, answer.ID, answer.Text, answer.Votes, answer.Votes, poll.ID)
		}
	})
}

func work(db *sql.DB, groupID int, workerID int, client *vkapi.VKClient) func(job updater.TopicUpdateJob) error {
	return func(job updater.TopicUpdateJob) error {
		log.Printf("Worker %d: processando job-%d (t%d)\n", workerID, job.ID, job.TopicID)

		topic, err := downloadTopic(client, groupID, job.TopicID)
		if err != nil {
			log.Printf("Worker %d: erro ao baixar tópico, job-%d (t%d), %v\n", workerID, job.ID, job.TopicID, err)
			return err
		}

		profilesTopic := []Profile{topic.CreatedBy, topic.UpdatedBy}

		err = insertTopic(db, topic)
		if err != nil {
			log.Printf("Worker %d: erro ao inserir tópico, job-%d (t%d), %v\n", workerID, job.ID, job.TopicID, err)
			return err
		}

		// NOTE(Pedro): We dont care about the topic anymore
		topic = Topic{}

		startComment := 0
		insertedComments := 0

		for {
			comments, total, profiles, poll, err := downloadComments(client, groupID, job.TopicID, startComment)
			if err != nil {
				log.Printf("Worker %d: erro ao baixar comentários, job-%d (t%d), %v\n", workerID, job.ID, job.TopicID, err)
				return err
			}

			if poll.ID != 0 {
				err = insertPoll(db, job.TopicID, poll)
				if err != nil {
					log.Printf("Worker %d: erro ao inserir poll, job-%d (t%d), %v\n", workerID, job.ID, job.TopicID, err)
					return err
				}
			}

			err = insertComments(db, job.TopicID, comments)
			if err != nil {
				log.Printf("Worker %d: erro ao inserir comentários, job-%d (t%d), %v\n", workerID, job.ID, job.TopicID, err)
				return err
			}

			err = insertProfiles(db, append(profiles, profilesTopic...))
			if err != nil {
				log.Printf("Worker %d: erro ao inserir profiles, job-%d (t%d), %v\n", workerID, job.ID, job.TopicID, err)
				return err
			}

			insertedComments += len(comments)
			startComment = comments[len(comments)-1].ID

			if insertedComments >= total {
				break
			}
		}

		log.Printf("Worker %d: processado job-%d (t%d) com sucesso\n", workerID, job.ID, job.TopicID)
		return nil
	}
}
