/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/leetsecure/qryptic-client-cli/internal/config"
	"github.com/leetsecure/qryptic-client-cli/internal/logger"
	"github.com/leetsecure/qryptic-client-cli/internal/platform"
	"github.com/leetsecure/qryptic-client-cli/internal/wireguard"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var storage *config.Storage
var wg *wireguard.WireGuardManager

var rootCmd = &cobra.Command{
	Use:   "qryptic",
	Short: "Client CLI for Qryptic",
	Long:  `Qryptic Client CLI will help you in connecting to Qryptic gateways`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log := logger.Default()
		log.Error(err.Error())
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	if os.Geteuid() != 0 {
		fmt.Println("Please run as a root user or with sudo permission")
		os.Exit(1)
	}
}

func initConfig() {
	storageRes, err := config.NewStorage(viper.GetViper())
	if err != nil {
		fmt.Println("Error in setting up/accessing the config file")
		os.Exit(1)
	}
	storage = storageRes
	wg = wireguard.NewWireGuardManager(platform.GetConfigDirectory(), platform.GetDefaultInterfaceName())

}
