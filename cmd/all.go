package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// allCmd represents the all command
var allCmd = &cobra.Command{
	Use:   "all",
	Short: "download all files from scientia",
	Long:  `TODO`,
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, course := range courses {
			if course.HasMaterials {
				err := downloadCourse(course, newOnly)
				if err != nil {
					fmt.Println(err)
				}
			}
		}

		return nil
	},
}

func init() {
	downloadCmd.AddCommand(allCmd)
}
