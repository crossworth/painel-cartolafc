package openid

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/markbates/goth"
	"golang.org/x/oauth2"
)

var (
	authURL      = "https://oauth.vk.com/authorize"
	tokenURL     = "https://oauth.vk.com/access_token"
	endpointUser = "https://api.vk.com/method/users.get"
	apiVersion   = "5.71"
)

// This is based on github.com/markbates/goth VK implementation
func New(clientKey, secret, callbackURL string, scopes ...string) *Provider {
	p := &Provider{
		ClientKey:    clientKey,
		Secret:       secret,
		CallbackURL:  callbackURL,
		providerName: "vk",
	}
	p.config = newConfig(p, scopes)
	return p
}

type Provider struct {
	ClientKey    string
	Secret       string
	CallbackURL  string
	HTTPClient   *http.Client
	config       *oauth2.Config
	providerName string
	version      string
}

func (p *Provider) Name() string {
	return p.providerName
}

func (p *Provider) SetName(name string) {
	p.providerName = name
}

func (p *Provider) Client() *http.Client {
	return goth.HTTPClientWithFallBack(p.HTTPClient)
}

func (p *Provider) BeginAuth(state string) (goth.Session, error) {
	url := p.config.AuthCodeURL(state)
	session := &Session{
		AuthURL: url,
	}

	return session, nil
}

func (p *Provider) FetchUser(session goth.Session) (goth.User, error) {
	sess := session.(*Session)
	user := goth.User{
		AccessToken: sess.AccessToken,
		Provider:    p.Name(),
		ExpiresAt:   sess.ExpiresAt,
		UserID:      strconv.Itoa(sess.ID),
	}

	if user.AccessToken == "" {
		return user, fmt.Errorf("%s não é possível retornar os dados do usuário sem um access token", p.providerName)
	}

	fields := "photo_200,nickname"
	requestURL := fmt.Sprintf("%s?fields=%s&access_token=%s&v=%s", endpointUser, fields, sess.AccessToken, apiVersion)
	response, err := p.Client().Get(requestURL)
	if err != nil {
		return user, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return user, fmt.Errorf("%s respondeu com status %d tentar solicita os dados do usuário", p.providerName, response.StatusCode)
	}

	bits, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return user, err
	}

	err = json.NewDecoder(bytes.NewReader(bits)).Decode(&user.RawData)
	if err != nil {
		return user, err
	}

	err = userFromReader(bytes.NewReader(bits), &user)
	return user, err
}

func userFromReader(reader io.Reader, user *goth.User) error {
	response := struct {
		Response []struct {
			ID        int64  `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			NickName  string `json:"nickname"`
			Photo200  string `json:"photo_200"`
		} `json:"response"`
	}{}

	err := json.NewDecoder(reader).Decode(&response)
	if err != nil {
		return err
	}

	if len(response.Response) == 0 {
		return fmt.Errorf("vk não é possível conseguir os dados do usuário")
	}

	u := response.Response[0]

	user.UserID = strconv.FormatInt(u.ID, 10)
	user.FirstName = u.FirstName
	user.LastName = u.LastName
	user.NickName = u.NickName
	user.AvatarURL = u.Photo200

	return err
}

func (p *Provider) Debug(debug bool) {}

func (p *Provider) RefreshToken(refreshToken string) (*oauth2.Token, error) {
	return nil, errors.New("refresh token não é fornecido pelo VK")
}

func (p *Provider) RefreshTokenAvailable() bool {
	return false
}

func newConfig(provider *Provider, scopes []string) *oauth2.Config {
	c := &oauth2.Config{
		ClientID:     provider.ClientKey,
		ClientSecret: provider.Secret,
		RedirectURL:  provider.CallbackURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
		Scopes: []string{},
	}

	defaultScopes := map[string]struct{}{}

	for _, scope := range scopes {
		if _, exists := defaultScopes[scope]; !exists {
			c.Scopes = append(c.Scopes, scope)
		}
	}

	return c
}
