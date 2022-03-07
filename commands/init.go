package commands

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/toolfordev/local-cli/application/services"
	"github.com/toolfordev/local-cli/infrastructure/variables"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		err := services.Init()
		if err != nil {
			panic(err)
		}
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter password: ")
		password, _ := reader.ReadString('\n')
		fmt.Print("Enter password again: ")
		password2, _ := reader.ReadString('\n')
		if password == password2 {
			for i := 0; i < 20; i++ {
				err = variables.SetPasswordEncrypted(password)
				if err == nil {
					break
				}
				time.Sleep(5 * time.Second)
			}
			if err != nil {
				panic(err)
			}
		}
		fmt.Println("Passwords does not match")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
