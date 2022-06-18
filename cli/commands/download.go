package commands

import (
	"context"
	"fmt"
	"github.com/goDoCer/scientiaCLI/cli/config"
	"github.com/goDoCer/scientiaCLI/scientia"
	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
	"path"
	"sync"
	"time"
)

var Download = &cli.Command{Name: "download",
	Usage: "download a file from scientia",
	Action: func(c *cli.Context) error {
		courseTitle := c.Args().First()
		courses, err := config.Client.GetCourses()
		if err != nil {
			return err
		}
		found := false

		// TODO? make it concurrent
		for _, course := range courses {
			if (course.Title == courseTitle || courseTitle == "all") && course.HasMaterials {
				func(course scientia.Course) {
					err := downloadCourse(course)
					if err != nil {
						fmt.Println(err)
					}
				}(course)
				found = true
			}
		}

		if !found {
			return errors.New("Course does not exist")
		}
		return nil
	},
	BashComplete: func(c *cli.Context) {
		fmt.Println("all")
		courses, _ := config.Client.GetCourses()
		for _, course := range courses {
			fmt.Println(course.Title)
		}
	}}

// downloadCourse downloads all the files for a course
func downloadCourse(course scientia.Course) error {
	// TODO: Add tests for this function
	files, err := config.Client.ListFiles(course.Code)
	if len(files) == 0 {
		log.Info("No files to download for course ", course.FullName())
		return nil
	}

	if err != nil {
		return err
	}
	dirName := course.FullName()
	saveDir := path.Join(config.Cfg.SaveDir, dirName)
	err = os.Mkdir(saveDir, 0777)

	if err != nil && !errors.Is(err, os.ErrExist) {
		if errors.Is(err, os.ErrNotExist) {
			log.Warnf("Directory %s does not exist", saveDir)
			log.Println("Trying to create directory", saveDir)
			err = os.Mkdir(config.Cfg.SaveDir, 0777)
			if err != nil {
				log.Fatal("Could not create directory", saveDir)
				return err
			}
			os.Mkdir(saveDir, 0777)
		} else {
			return err
		}
	}

	bar := progressbar.Default(int64(len(files)), fmt.Sprintf("Downloading files for course: %s", course.Code))

	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go func(resource scientia.Resource) {
			defer wg.Done()
			defer bar.Add(1)
			filepath := path.Join(saveDir, resource.Title)

			fileInfo, err := os.Stat(filepath)
			if err == nil {
				scientiaLastModified, err := config.Client.GetFileLastModified(resource.ID)
				if err != nil || scientiaLastModified.After(fileInfo.ModTime()) {
					log.Warnf("skipping download for file %s because it has not been updated", resource.Title)
					return
				}
			}

			downloadCtx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
			defer cancel()
			data, err := config.Client.Download(downloadCtx, resource.ID)
			if err != nil {
				fmt.Println(err) //TODO: send this to a channel
			}

			err = os.WriteFile(filepath, data, 0777)
		}(file)
	}

	wg.Wait()
	return nil
}
