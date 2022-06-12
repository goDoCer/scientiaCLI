package scientia

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
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
	}).Debug("Successfully logged in")

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

// Do tries to make a request using the auth token
// if the request fails because the auth token has expired
//  it uses the refresh token to get a new auth token and make the request again
func (c *APIClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	resp, err := c.Client.Do(req)
	// token has not expired
	if err != nil || (resp.StatusCode == http.StatusOK && resp.StatusCode != http.StatusUnauthorized) {
		return resp, err
	}

	refreshReq, _ := http.NewRequest("POST", baseURL+"auth/refresh", nil)

	refreshReq.Header.Add("Authorization", "Bearer "+c.refreshToken)
	refreshResp, err := c.Client.Do(refreshReq)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't refresh access token, please login again")
	}

	var refreshResponse LoginTokens
	err = json.NewDecoder(refreshResp.Body).Decode(&refreshResponse)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't decode refresh token response")
	}

	c.accessToken = refreshResponse.AccessToken
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	return c.Client.Do(req)
}

// GetCourses fetchs the courses for the current academic year
func (c *APIClient) GetCourses() ([]Course, error) {
	year := getCurrentAcademicYear()

	req, err := http.NewRequest("GET", baseURL+"courses/"+year, nil)

	if err != nil {
		panic(err)
	}
	resp, err := c.Do(req)

	if err != nil {
		return nil, err
	}

	var courses []Course
	err = json.NewDecoder(resp.Body).Decode(&courses)
	if err != nil {
		return nil, errors.New("Error fetching your courses from scientia, have you logged in?")
	}
	return courses, nil
}

// ListFiles returns the list of resources that are files for a course
func (c *APIClient) ListFiles(courseCode string) ([]Resource, error) {
	year := getCurrentAcademicYear()

	req, err := http.NewRequest("GET", fmt.Sprintf("%sresources?year=%s&course=%s", baseURL, year, courseCode), nil)

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
			ext := path.Ext(resource.Title)
			if ext == "" {
				resource.Title = resource.Title + "." + path.Ext(resource.Path)
			}

			files = append(files, resource)
		}
	}

	return files, nil
}

//GetFileLastModified returns the last time the resource was modified on the server
func (c *APIClient) GetFileLastModified(resourceID int) (time.Time, error) {
	req, err := http.NewRequest("HEAD", fmt.Sprintf("%sresources/%d/file", baseURL, resourceID), nil)
	if err != nil {
		return time.Now(), err
	}

	lastModified := req.Header.Get("last-modified")
	if lastModified == "" {
		return time.Now(), errors.New("Couldn't get last modified header")
	}

	return time.Parse(time.RFC1123, lastModified)
}

// Download downloads the given resource from the API
func (c *APIClient) Download(ctx context.Context, resourceID int) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%sresources/%d/file", baseURL, resourceID), nil)
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
