/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/leetsecure/qryptic-client-cli/internal/logger"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout, Cleanup and Reset",
	Long:  `This will reset your Qrytic CLI as new one and remove all the saved data along with logging you out.`,
	Run: func(cmd *cobra.Command, args []string) {
		log := logger.Default()
		err := wg.StopVPN()
		if err != nil {
			fmt.Printf("Error in stopping the running qryptic client \n %s \n", err.Error())
			fmt.Printf("disconnect before logging out ...")
			return
		}
		storage.ClearConfig()

		log.Info("Logged out ....")
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
