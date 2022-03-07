package commands

import (
	"github.com/spf13/cobra"
	"github.com/toolfordev/local-cli/application/services"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		err := services.Start("./toolfordev.yaml")
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
