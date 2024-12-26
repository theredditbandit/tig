package cmd

import (
	"tig/pkg/repo"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var test = &cobra.Command{
	Use:   "test",
	Short: "test random things , remove this later",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, err := repo.FindGitRepo("", true)
		if err != nil {
			return err
		}
		if repo != nil {
			log.Info("this is inside a repo")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(test)
}
