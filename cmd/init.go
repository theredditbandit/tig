package cmd

import (
	"errors"
	"fmt"
	"path/filepath"
	"tig/pkg/repo"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new, empty repository.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("init command needs one argument to specify where to create the repo")
		}
		path := args[0]
		_, err := repo.InitRepo(path)
		if err != nil {
			return err
		}
		absPath, _ := filepath.Abs(path)
		projPath := filepath.Join(absPath, ".git/")
		fmt.Printf("Initialized empty Git reposiroty in %s\n", projPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
