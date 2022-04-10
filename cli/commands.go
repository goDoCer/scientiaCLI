package cli

import (
	"fmt"
	"os"
	"path"
	"sync"
	"syscall"

	"github.com/goDoCer/scientiaCLI/scientia"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/schollz/progressbar/v3"
	cli "github.com/urfave/cli/v2"
	"golang.org/x/term"
)

var commands = []*cli.Command{
	{
		Name:      "login",
		Usage:     "login to the scientia API",
		ArgsUsage: "shortcode",
		Action: func(c *cli.Context) error {
			var shortcode string

			if c.NArg() < 1 {
				log.Warn("No shortcode provided")
				fmt.Print("Please enter your shortcode: ")
				fmt.Scanln(&shortcode)
			} else {
				shortcode = c.Args().First()
				log.Infof("shortcode: %s", shortcode)
			}

			if shortcode == "" {
				return errors.New("Please enter your shortcode")
			}

			fmt.Print("Enter password (It will be hidden): ")
			bytePassword, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				return err
			}
			password := string(bytePassword)
			err = client.Login(shortcode, password)
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
			log.Info("Save directory set to ", cfg.SaveDir)
			return cfg.save(configPath)
		},
	},
	{
		Name:  "config",
		Usage: "view the CLI configuration",
		Action: func(c *cli.Context) error {
			fmt.Println("Save directory: ", cfg.SaveDir)
			fmt.Println("Access token: ", cfg.AccessToken)
			fmt.Println("Refresh token: ", cfg.RefreshToken)
			return nil
		},
	},
}

// downloadCourse downloads all the files for a course
func downloadCourse(course scientia.Course) error {
	files, err := client.ListFiles(course.Code)
	if len(files) == 0 {
		log.Info("No files to download for course ", course.FullName())
		return nil
	}

	if err != nil {
		return err
	}
	dirName := course.FullName()
	saveDir := path.Join(cfg.SaveDir, dirName)
	err = os.Mkdir(saveDir, 0777)

	if err != nil && !errors.Is(err, os.ErrExist) {
		if errors.Is(err, os.ErrNotExist) {
			log.Warnf("Directory %s does not exist", saveDir)
			log.Println("Trying to create directory", saveDir)
			err = os.Mkdir(cfg.SaveDir, 0777)
			if err != nil {
				log.Fatal("Could not create directory", saveDir)
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
			defer bar.Add(1)
			filepath := path.Join(saveDir, resource.Title)

			fileInfo, err := os.Stat(filepath)
			if err == nil {
				scientiaLastModified, err := client.GetFileLastModified(resource.ID)
				if err != nil || scientiaLastModified.After(fileInfo.ModTime()) {
					log.Warnf("skipping download for file %s because it has not been updated", resource.Title)
					return
				}
			}
			data, err := client.Download(resource.ID)
			if err != nil {
				fmt.Println(err) //TODO: send this to a channel
			}

			err = os.WriteFile(filepath, data, 0777)
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
