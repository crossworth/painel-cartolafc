package database

import (
	"github.com/crossworth/cartola-web-admin/model"
)

func (d *DatabaseService) FindProfileByID(id int) ([]model.Profile, error) {
	return []model.Profile{}, nil
}

func (d *DatabaseService) FindTopicByUser(id int) ([]model.Topic, error) {
	return []model.Topic{}, nil
}

func (d *DatabaseService) FindTopicCount(id int) (uint64, error) {
	return 0, nil
}

func (d *DatabaseService) FindCommentsByUser(id int) ([]model.Comment, error) {
	return []model.Comment{}, nil
}

func (d *DatabaseService) FindCommentCount(id int) (uint64, error) {
	return 0, nil
}

func (d *DatabaseService) FindProfileHistory(id int) ([]model.ProfileNames, error) {
	return []model.ProfileNames{}, nil
}
