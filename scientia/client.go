package scientia

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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
			Timeout: time.Second * 100,
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

// Login - takes the student's shortcode and password to generate an authorisation token
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

// GetTokens returns the tokens as a LoginTokenStruct which can marshalled to json
func (c *APIClient) GetTokens() LoginTokens {
	return LoginTokens{
		AccessToken:  c.accessToken,
		RefreshToken: c.refreshToken,
	}
}

// AddTokens adds the tokens to the client
func (c *APIClient) AddTokens(tokens LoginTokens) {
	c.accessToken = tokens.AccessToken
	c.refreshToken = tokens.RefreshToken
}

// GetCourses fetchs the courses for the current academic year
func (c *APIClient) GetCourses() ([]Course, error) {
	year := getCurrentAcademicYear()

	req, err := http.NewRequest("GET", baseURL+"courses/"+year, nil)
	req.Header.Add("Authorization", "Bearer "+c.accessToken)

	if err != nil {
		panic(err)
	}
	resp, err := c.Do(req)
	if err != nil {
		// TODO? should this be a panic? what if we used an expired token?
		panic(err)
	}

	var courses []Course
	err = json.NewDecoder(resp.Body).Decode(&courses)
	if err != nil {
		return nil, errors.New("Error fetching your courses from scientia, have you logged in?")
	}
	return courses, err
}

// ListFiles returns the list of resources that are files for a course
func (c *APIClient) ListFiles(courseCode string) ([]Resource, error) {
	year := getCurrentAcademicYear()

	req, err := http.NewRequest("GET", fmt.Sprintf("%sresources?year=%s&course=%s", baseURL, year, courseCode), nil)
	req.Header.Add("Authorization", "Bearer "+c.accessToken)

	if err != nil {
		panic(err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	var resources []Resource
	err = json.NewDecoder(resp.Body).Decode(&resources)

	files := make([]Resource, 0)
	for _, resource := range resources {
		if resource.Type == "file" {
			// TODO: HACKY - we dont actually know if it should be a pdf but it is more often than not a pdf
			if !strings.HasSuffix(resource.Title, ".pdf") {
				resource.Title = resource.Title + ".pdf"
			}

			files = append(files, resource)
		}
	}

	return files, nil
}

// Download downloads the given resource from the API
func (c *APIClient) Download(resource Resource) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%sresources/%d/file", baseURL, resource.ID), nil)
	req.Header.Add("Authorization", "Bearer "+c.accessToken)

	if err != nil {
		panic(err)
	}
	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
