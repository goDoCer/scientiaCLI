package scientia

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"
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
func (c *APIClient) GetCourses() []Course {
	year := getCurrentAcademicYear()

	req, err := http.NewRequest("GET", baseURL+"courses/"+year, nil)
	req.Header.Add("Authorization", "Bearer "+c.accessToken)

	if err != nil {
		panic(err)
	}
	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}

	var courses []Course
	err = json.NewDecoder(resp.Body).Decode(&courses)
	return courses
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
			files = append(files, resource)
		}
	}

	return files, nil
}

// Download downloads the given resource from the API
func (c *APIClient) Download(resource Resource) error {
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

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fileExtension := path.Ext(resource.Path)
	return os.WriteFile(resource.Title+fileExtension, data, 0o777)
}

// DownloadCourse downloads all the files for a course
func (c *APIClient) DownloadCourse(course Course) error {
	files, err := c.ListFiles(course.Code)
	if err != nil {
		return err
	}
	dirName := course.Code + "-" + course.Title
	os.Mkdir(dirName, 0o777)
	os.Chdir(dirName)

	// TODO: This feels disgusting so maybe find a better way to separate concerns
	bar := progressbar.Default(int64(len(files)), "Downloading files")

	for _, file := range files {
		if err := c.Download(file); err != nil {
			return err
		}
		bar.Add(1)
	}
	os.Chdir("..")
	return nil
}
