package cli

import (
	"fmt"
	"os"
	"path"
	"sync"
	"syscall"

	"scientia-cli/scientia"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/schollz/progressbar/v3"
	cli "github.com/urfave/cli/v2"
	"golang.org/x/term"
)

var (
	client     scientia.APIClient
	cfg        config
	configPath string
)

func init() {
	client = scientia.NewAPIClient()

	filepath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	// reading the config from the executable's directory
	configPath = path.Dir(filepath) + "/config.json"
	cfg, err = loadConfig(configPath)
	if err != nil {
		panic(err)
	}
	tokens := cfg.tokens()
	client.AddTokens(tokens)
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

			cfg.updateTokens(client.GetTokens())
			cfg.save(configPath)

			return nil
		},
	},
	{
		Name:  "download",
		Usage: "download a file from scientia",
		Action: func(c *cli.Context) error {
			courseTitle := c.Args().First()
			courses, err := client.GetCourses()
			if err != nil {
				return err
			}
			found := false

			var wg sync.WaitGroup

			for _, course := range courses {
				if (course.Title == courseTitle || courseTitle == "all") && course.HasMaterials {
					wg.Add(1)
					go func(course scientia.Course) {
						defer wg.Done()
						err := downloadCourse(course)
						if err != nil {
							fmt.Println(err) //TODO: send this to a channel
						}
					}(course)
					found = true
				}
			}

			wg.Wait()
			if !found {
				return errors.New("Course does not exist")
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			fmt.Println("all")
			courses, _ := client.GetCourses()
			for _, course := range courses {
				fmt.Println(course.Title)
			}
		},
	},
	{
		Name:  "save-dir",
		Usage: "set the directory to save downloaded files",
		Action: func(c *cli.Context) error {
			dir := c.Args().First()
			if dir == "" {
				return errors.New("Please enter a directory")
			}
			cfg.SaveDir = dir
			logrus.Info("Save directory set to ", cfg.SaveDir)
			return cfg.save(configPath)
		},
	},
	{
		Name:  "config",
		Usage: "view the CLI configuration",
		Action: func(c *cli.Context) error {
			fmt.Println("Your save directory is set to:", cfg.SaveDir)
			return nil
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
	saveDir := path.Join(cfg.SaveDir, dirName)
	err = os.Mkdir(saveDir, 0777)

	if err != nil && !errors.Is(err, os.ErrExist) {
		if errors.Is(err, os.ErrNotExist) {
			logrus.Warnf("Directory %s does not exist", saveDir)
			logrus.Println("Trying to create directory", saveDir)
			err = os.Mkdir(cfg.SaveDir, 0777)
			if err != nil {
				logrus.Fatal("Could not create directory", saveDir)
				return err
			}
			os.Mkdir(saveDir, 0777)
		} else {
			return err
		}
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

			filepath := path.Join(saveDir, resource.Title)

			err = os.WriteFile(filepath, data, 0777)
			bar.Add(1)
		}(file)
	}

	wg.Wait()
	return nil
}

//Sad times for this abstraction
// func pFor[T any](tasks []T, function func(task T) []error) []error {

// 	errChannel := make(chan error, 10)
// 	var wg sync.WaitGroup

// 	for _, task := range tasks {
// 		wg.Add(1)
// 		go func(task T) {
// 			defer wg.Done()
// 			for _, err := range function(task) {
// 				if err != nil {
// 					errChannel <- err
// 				}
// 			}
// 		}(task)
// 	}

// 	go func(errChannel chan error) {
// 		wg.Wait()
// 		close(errChannel)
// 	}(errChannel)

// 	allErrors := make([]error, 0)
// 	for err := range errChannel {
// 		allErrors = append(allErrors, err)
// 	}

// 	return allErrors
// }
