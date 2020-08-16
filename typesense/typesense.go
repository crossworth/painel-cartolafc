package typesense

import (
	"errors"
	"strconv"

	"github.com/GianOrtiz/typesense-go"

	"github.com/crossworth/cartola-web-admin/model"
)

type TypeSense struct {
	client *typesense.Client
}

func NewSearch(host string, port string, apiKey string) (*TypeSense, error) {
	s := &TypeSense{}
	s.client = typesense.NewClient(
		&typesense.Node{
			Host:     host,
			Port:     port,
			Protocol: "http",
			APIKey:   apiKey,
		},
		2,
	)

	err := s.client.Ping()
	if err != nil {
		return s, err
	}

	return s, nil
}

func (t *TypeSense) CreateCollections() error {
	topicsSchema := typesense.CollectionSchema{
		Name: "topics",
		Fields: []typesense.CollectionField{
			{
				Name: "id",
				Type: "string",
			},
			{
				Name: "title",
				Type: "string",
			},
			{
				Name: "created_at",
				Type: "int32",
			},
			{
				Name: "updated_at",
				Type: "int32",
			},
			{
				Name: "created_by",
				Type: "int32",
			},
		},
		DefaultSortingField: "created_at",
	}

	_, err := t.client.RetrieveCollection("topics")
	if err != nil && !errors.Is(err, typesense.ErrCollectionNotFound) {
		return err
	}

	if errors.Is(err, typesense.ErrCollectionNotFound) {
		_, err := t.client.CreateCollection(topicsSchema)
		if err != nil {
			return err
		}
	}

	commentsSchema := typesense.CollectionSchema{
		Name: "comments",
		Fields: []typesense.CollectionField{
			{
				Name: "id",
				Type: "string",
			},
			{
				Name: "text",
				Type: "string",
			},
			{
				Name: "date",
				Type: "int32",
			},
			{
				Name: "topic_id",
				Type: "int32",
			},
			{
				Name: "from_id",
				Type: "int32",
			},
		},
		DefaultSortingField: "date",
	}

	_, err = t.client.RetrieveCollection("comments")
	if err != nil && !errors.Is(err, typesense.ErrCollectionNotFound) {
		return err
	}

	if errors.Is(err, typesense.ErrCollectionNotFound) {
		_, err := t.client.CreateCollection(commentsSchema)
		if err != nil {
			return err
		}
	}

	profilesSchema := typesense.CollectionSchema{
		Name: "profiles",
		Fields: []typesense.CollectionField{
			{
				Name: "id",
				Type: "string",
			},
			{
				Name: "first_name",
				Type: "string",
			},
			{
				Name: "last_name",
				Type: "string",
			},
			{
				Name: "screen_name",
				Type: "string",
			}, {
				Name: "photo",
				Type: "string",
			},
			{
				Name: "id_int",
				Type: "int32",
			},
		},
		DefaultSortingField: "id_int",
	}

	_, err = t.client.RetrieveCollection("profiles")
	if err != nil && !errors.Is(err, typesense.ErrCollectionNotFound) {
		return err
	}

	if errors.Is(err, typesense.ErrCollectionNotFound) {
		_, err := t.client.CreateCollection(profilesSchema)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *TypeSense) DropCollection(name string) error {
	_, err := t.client.DeleteCollection(name)
	return err
}

type Topic struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
	CreatedBy int    `json:"created_by"`
}

type Comment struct {
	ID      string `json:"id"`
	Text    string `json:"text"`
	Date    int    `json:"date"`
	TopicID int    `json:"topic_id"`
	FromID  int    `json:"from_id"`
}

type Profile struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	ScreenName string `json:"screen_name"`
	Photo      string `json:"photo"`
	IDInt      int    `json:"id_int"`
}

func (t *TypeSense) InsertTopic(topic Topic) error {
	documentResponse := t.client.IndexDocument("topics", topic)
	return documentResponse.Error
}

func (t *TypeSense) InsertComment(comment Comment) error {
	documentResponse := t.client.IndexDocument("comments", comment)
	return documentResponse.Error
}

func (t *TypeSense) InsertProfile(profile Profile) error {
	documentResponse := t.client.IndexDocument("profiles", profile)
	return documentResponse.Error
}

func ToTypeSenseTopic(topic model.Topic) Topic {
	return Topic{
		ID:        strconv.Itoa(topic.ID),
		Title:     topic.Title,
		CreatedAt: topic.CreatedAt,
		UpdatedAt: topic.UpdatedAt,
		CreatedBy: topic.CreatedBy,
	}
}

func ToTypeSenseComment(comment model.Comment) Comment {
	return Comment{
		ID:      strconv.Itoa(comment.ID),
		Text:    comment.Text,
		Date:    comment.Date,
		TopicID: comment.TopicID,
		FromID:  comment.FromID,
	}
}

func ToTypeSenseProfile(profile model.Profile) Profile {
	return Profile{
		ID:         strconv.Itoa(profile.ID),
		FirstName:  profile.FirstName,
		LastName:   profile.LastName,
		ScreenName: profile.ScreenName,
		Photo:      profile.Photo,
		IDInt:      profile.ID,
	}
}
