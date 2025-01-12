/*
Copyright © 2025 Leetsecure hello@leetsecure.com
*/
package cmd

import (
	"fmt"

	"github.com/leetsecure/qryptic-client-cli/internal/logger"
	"github.com/spf13/cobra"
)

// disconnectCmd represents the disconnect command
var disconnectCmd = &cobra.Command{
	Use:   "disconnect",
	Short: "Disconnect from Qryptic Gateway",
	Long:  `Disconnect from any connection with Qryptic Gateway`,
	Run: func(cmd *cobra.Command, args []string) {
		log := logger.Default()
		storage.GetConnectedToGateway()
		err := wg.StopVPN()
		if err != nil {
			fmt.Printf("Error in stopping the running qryptic client \n %s \n", err.Error())

		}
		log.Info("Qryptic disconnected")
	},
}

func init() {
	rootCmd.AddCommand(disconnectCmd)

}
