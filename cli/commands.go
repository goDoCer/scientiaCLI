package cli

import (
	"errors"
	"fmt"
	"os"
	"path"
	"syscall"

	"scientia-cli/scientia"

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
		BashComplete: func(c *cli.Context) {
			fmt.Fprintf(c.App.Writer, "--better\n")
		},
	},
	{
		Name:  "download",
		Usage: "download a file from scientia",
		Action: func(c *cli.Context) error {
			courses := client.GetCourses()
			courseTitle := c.Args().First()

			for _, course := range courses {
				if course.Title == courseTitle {
					err := client.DownloadCourse(course)
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
