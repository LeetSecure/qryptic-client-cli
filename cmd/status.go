/*
Copyright © 2025 Leetsecure hello@leetsecure.com
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var StatusDebug bool

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Current Status of Qryptic",
	Long:  `Check if you are connected to any Qryptic Gateway or not`,
	Run: func(cmd *cobra.Command, args []string) {
		uuid, name, exists := storage.GetConnectedToGateway()
		if exists {
			wgclientConfig, _ := storage.GetQrypticClient(uuid)
			fmt.Printf("Connected to %s gateway at %s:%d\n", name, wgclientConfig.WGClientPeerConfig.VpnGatewayIP, wgclientConfig.WGClientPeerConfig.VpnGatewayPort)
		} else {
			fmt.Println("Qryptic is not running")
		}
		if !StatusDebug {
			return
		}
		isRunning, output, err := wg.CheckStatus()
		if err != nil {
			fmt.Println("Error checking the current status")
			fmt.Println(err)
			return
		}
		if isRunning {
			fmt.Println(output)
			return
		}
		fmt.Println("Qryptic is not running")
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().BoolVarP(&StatusDebug, "debug", "d", false, "Debug actual wireguard logs")
}
