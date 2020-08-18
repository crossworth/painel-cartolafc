package openid

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/markbates/goth"
)

type Session struct {
	AuthURL     string
	AccessToken string
	ExpiresAt   time.Time
	ID          int
}

func (s *Session) GetAuthURL() (string, error) {
	if s.AuthURL == "" {
		return "", errors.New(goth.NoAuthUrlErrorMessage)
	}
	return s.AuthURL, nil
}

func (s *Session) Marshal() string {
	b, _ := json.Marshal(s)
	return string(b)
}

func (s *Session) Authorize(provider goth.Provider, params goth.Params) (string, error) {
	p := provider.(*Provider)
	token, err := p.config.Exchange(goth.ContextForClient(p.Client()), params.Get("code"))
	if err != nil {
		return "", err
	}

	if !token.Valid() {
		return "", errors.New("token inválido recebido")
	}

	id, ok := token.Extra("user_id").(float64)
	if !ok {
		return "", errors.New("não foi possível conseguir o ID do usuário")
	}

	s.AccessToken = token.AccessToken
	s.ExpiresAt = token.Expiry
	s.ID = int(id)
	return s.AccessToken, err
}

func (p *Provider) UnmarshalSession(data string) (goth.Session, error) {
	sess := new(Session)
	err := json.NewDecoder(strings.NewReader(data)).Decode(&sess)
	return sess, err
}
