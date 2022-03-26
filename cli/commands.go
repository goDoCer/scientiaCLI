package cli

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sync"
	"syscall"

	"scientia-cli/scientia"

	"github.com/schollz/progressbar/v3"
	cli "github.com/urfave/cli/v2"
	"golang.org/x/term"
)

var (
	client    scientia.APIClient
	tokenPath string
)

func init() {
	client = scientia.NewAPIClient()
	filepath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	tokenPath = path.Dir(filepath) + "/token.txt"
	loginDetails, _ := loadDetails()
	client.AddTokens(*loginDetails)
}

var commands = []*cli.Command{
	{
		Name:  "login",
		Usage: "login to the scientia API",
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return errors.New("missing shortcode")
			}

			username := c.Args().Get(0)
			if username == "" {
				return errors.New("Please enter your shortcode")
			}
			fmt.Print("Enter password: ")
			bytePassword, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				return err
			}
			password := string(bytePassword)
			err = client.Login(username, password)
			if err != nil {
				return err
			}
			return saveDetails(client.GetTokens())
		},
	},
	{
		Name:  "download",
		Usage: "download a file from scientia",
		Action: func(c *cli.Context) error {
			courseTitle := c.Args().First()
			courses := client.GetCourses()

			for _, course := range courses {
				if course.Title == courseTitle {
					err := downloadCourse(course)
					if err != nil {
						return err
					}
					return nil
				}
			}

			return errors.New("Course does not exist")
		},
		BashComplete: func(c *cli.Context) {
			courses := client.GetCourses()
			for _, course := range courses {
				fmt.Println(course.Title)
			}
		},
	},
}

// downloadCourse downloads all the files for a course
func downloadCourse(course scientia.Course) error {
	files, err := client.ListFiles(course.Code)
	if err != nil {
		return err
	}
	dirName := course.Code + "-" + course.Title

	err = os.Mkdir(dirName, 0777)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}

	bar := progressbar.Default(int64(len(files)), "Downloading files")

	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go func(resource scientia.Resource) {
			defer wg.Done()
			data, err := client.Download(resource)
			if err != nil {
				fmt.Println(err) //TODO: send this to a channel
			}

			filepath := path.Join(dirName, resource.Title)
			err = os.WriteFile(filepath, data, 0777)
			bar.Add(1)
		}(file)
	}

	wg.Wait()
	return nil
}
