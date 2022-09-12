package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// courseCmd represents the course command
var (
	courseCode string

	courseCmd = &cobra.Command{
		Use:   "course",
		Short: "downloads all file for a course",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			found := false

			for _, course := range courses {
				if (course.Code == courseCode) && course.HasMaterials {
					err := downloadCourse(course, unmodifiedOnly)
					if err != nil {
						fmt.Println(err)
					}
					found = true
					break
				}
			}

			if !found {
				return errors.New("specified course does not exist")
			}

			return nil
		},
	}
)

func init() {
	downloadCmd.AddCommand(courseCmd)
	courseCmd.PersistentFlags().StringVarP(&courseCode, "code", "c", "", "course code")

	courseCmd.MarkFlagRequired("code")
	//TODO: add support
	// courseCmd.PersistentFlags().StringVarP(&courseName, "name", "cn", "", "course name")

	// courseCmd.MarkFlagsMutuallyExclusive("code", "name")
}
