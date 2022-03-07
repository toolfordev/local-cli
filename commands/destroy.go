package commands

import (
	"github.com/spf13/cobra"
	"github.com/toolfordev/local-cli/application/services"
)

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		err := services.Destroy()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
}
