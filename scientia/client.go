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

const (
	scientiaBaseURL        = "https://scientia.doc.ic.ac.uk/api/"
	accessTokenCookieName  = "access_token_cookie"
	refreshTokenCookieName = "refresh_token_cookie"
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
		baseURL: scientiaBaseURL,
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

	resp, err := c.Post(c.baseURL+"auth/login", "application/json", bytes.NewBuffer(jsonValue))
	if err := checkResponse(resp, err); err != nil {
		return err
	}

	c.setAuthTokens(resp)

	log.WithFields(log.Fields{
		"accessToken":  c.accessToken,
		"refreshToken": c.refreshToken,
	}).Debug("Successfully logged in")

	return nil
}

func (c *APIClient) setAuthTokens(resp *http.Response) {
	for _, cookie := range resp.Cookies() {
		fmt.Println(cookie.Name, cookie.Value)
		switch cookie.Name {
		case accessTokenCookieName:
			c.accessToken = cookie.Value
		case refreshTokenCookieName:
			c.refreshToken = cookie.Value
		}
	}
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
	req.Header.Set("Cookie", "access_token_cookie="+c.accessToken)
	resp, err := c.Client.Do(req)

	// token has not expired
	if err != nil || (resp.StatusCode == http.StatusOK && resp.StatusCode != http.StatusUnauthorized) {
		return resp, err
	}

	refreshReq, _ := http.NewRequest("POST", c.baseURL+"auth/refresh", nil)

	req.Header.Set("Cookie", "refresh_token_cookie="+c.refreshToken)
	refreshResp, err := c.Client.Do(refreshReq)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't refresh access token, please login again")
	}

	c.setAuthTokens(refreshResp)
	req.Header.Set("Cookie", "access_token_cookie="+c.accessToken)
	return c.Client.Do(req)
}

// GetCourses fetchs the courses for the current academic year
func (c *APIClient) GetCourses() ([]Course, error) {
	year := getCurrentAcademicYear()

	req, err := http.NewRequest("GET", c.baseURL+"courses/"+year, nil)
	if err != nil {
		panic(err)
	}
	resp, err := c.Do(req)
	if err = checkResponse(resp, err); err != nil {
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

	req, err := http.NewRequest("GET", fmt.Sprintf("%sresources?year=%s&course=%s", c.baseURL, year, courseCode), nil)

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
				resource.Title = resource.Title + path.Ext(resource.Path)
			}

			files = append(files, resource)
		}
	}

	return files, nil
}

//GetFileLastModified returns the last time the resource was modified on the server
func (c *APIClient) GetFileLastModified(resourceID int) (time.Time, error) {
	req, err := http.NewRequest("HEAD", fmt.Sprintf("%sresources/%d/file", c.baseURL, resourceID), nil)
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
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%sresources/%d/file", c.baseURL, resourceID), nil)
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
