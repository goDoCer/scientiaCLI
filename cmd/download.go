package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/goDoCer/scientiaCLI/scientia"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var (
	unmodifiedOnly bool
	courses        []scientia.Course
	downloadCmd    = &cobra.Command{
		Use: "download",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			courses, err = client.GetCourses()
			if err != nil {
				return err
			}
			return nil
		},
		Short: "download files from scientia",
	}
)

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.PersistentFlags().BoolVarP(&unmodifiedOnly, "unmodified-only", "n", false, "do not overwrite existing files")
}

// TODO!
// UPDATE downloadCourse to send logs to a channel and then print them together
// UPDATE downloadCourse to send errots to a channel and then print them together
// downloadCourse downloads all the files for a course
func downloadCourse(course scientia.Course, unmodifiedOnly bool) error {
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
				log.Fatal("Could not create directory", saveDir, " because of ", err)
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
				// File already exists
				if unmodifiedOnly {
					log.Warnf("skipping download for file %s because it already exists and the unmodified-only flag is on", resource.Title)
					return
				}

				fileLastModified, err := client.GetFileLastModified(resource.ID)
				if err != nil || fileLastModified.After(fileInfo.ModTime()) {
					log.Warnf("skipping download for file %s because it has not been updated", resource.Title)
					return
				}
			}

			downloadCtx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
			defer cancel()
			data, err := client.Download(downloadCtx, resource.ID)
			if err != nil {
				fmt.Println(err) //TODO: send this to a channel
			}

			err = os.WriteFile(filepath, data, 0777)
			if err != nil {
				fmt.Println(err) //TODO: send this to a channel
			}
		}(file)
	}

	wg.Wait()
	return nil
}
