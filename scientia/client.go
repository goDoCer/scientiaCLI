package scientia

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// APIClient is the client for the scientia API
type APIClient struct {
	http.Client

	baseURL      string
	accessToken  string
	refreshToken string
}

// NewAPIClient creates a new instance of APIClient
func NewAPIClient() APIClient {
	return APIClient{
		Client: http.Client{
			Timeout: time.Second * 10,
		},
		baseURL: "https://api-materials.doc.ic.ac.uk/",
	}
}

type loginDetails struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginTokens contains the access and refresh tokens
type LoginTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

//Login - takes the student's shortcode and password to generate an authorisation token
func (c *APIClient) Login(username string, password string) error {
	details := loginDetails{Username: username, Password: password}

	jsonValue, err := json.Marshal(details)
	if err != nil {
		return err
	}

	resp, err := c.Post(baseURL+"auth/login", "application/json", bytes.NewBuffer(jsonValue))
	if err := checkResponse(resp, err); err != nil {
		return err
	}

	var response LoginTokens
	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		return errors.Wrap(err, "Error decoding login response")
	}

	c.accessToken = response.AccessToken
	c.refreshToken = response.RefreshToken

	log.WithFields(log.Fields{
		"accessToken":  c.accessToken,
		"refreshToken": c.refreshToken,
	}).Info("Successfully logged in")

	return nil
}

func (c *APIClient) GetTokens() LoginTokens {
	return LoginTokens{
		AccessToken:  c.accessToken,
		RefreshToken: c.refreshToken,
	}
}

func (c *APIClient) AddTokens(tokens LoginTokens) {
	c.accessToken = tokens.AccessToken
	c.refreshToken = tokens.RefreshToken
}

func (c *APIClient) GetCourses() {
	year := getCurrentAcademicYear()

	resp, err := c.Get(baseURL + "courses/" + year)
	if err != nil {
		panic(err)
	}
	log.Info(resp)
}
