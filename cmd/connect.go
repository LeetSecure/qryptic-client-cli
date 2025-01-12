/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/leetsecure/qryptic-client-cli/internal/auth"
	"github.com/leetsecure/qryptic-client-cli/internal/client"
	"github.com/leetsecure/qryptic-client-cli/internal/config"
	"github.com/leetsecure/qryptic-client-cli/internal/logger"
	"github.com/leetsecure/qryptic-client-cli/internal/models"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to Qryptic gateway",
	Long:  `Connect to any of the accessible Qryptic gateway`,
	Run: func(cmd *cobra.Command, args []string) {
		listAccessibleGateways()
	},
}

func listAccessibleGateways() {
	log := logger.Default()
	baseUrl, _ := storage.GetBaseUrl()

	isValidUrl := auth.IsURL(baseUrl)
	if !isValidUrl {
		log.Error("Please login first ... ")
		return
	}
	if !auth.IsBaseUrlHealthy() {
		log.Error("Check if the Qryptic service is running at the", "set URL", baseUrl)
		return
	}
	if !auth.IsAuthTokenValid() {
		log.Error("Please authenticate ...")
		return
	}
	authToken, _ := storage.GetAuthToken()
	qrypticClient := client.NewQrypticClient(baseUrl, authToken)
	statusCode, resp, err := qrypticClient.ListAccessibleGateways()
	if err != nil {
		log.Error(err.Error())
		return
	}
	if statusCode == http.StatusOK {
		selectGateway(*resp)
	} else if statusCode == http.StatusUnauthorized {
		log.Error("Please authenticate ...")
	} else {
		log.Error("Server Issue ...")
	}
}

func selectGateway(gateways []models.GatewayResponse) {
	log := logger.Default()
	loginMethodPromptContent := models.PromptContent{
		ErrorMsg: "Please select a valid login method.",
		Label:    "Select a login method.",
	}
	gatewaySelected, index := promptGatewaySelect(loginMethodPromptContent, gateways)
	log.Info("The selected gateway is ", "name", gatewaySelected)

	connectToGateway(gateways[index].Uuid, gateways[index].Name)
}

func clientExisiting(uuid string) (bool, models.WGClientConfig) {
	qrypticClient, _ := storage.GetQrypticClient(uuid)
	if qrypticClient.ClientUuid == "" {
		return false, qrypticClient
	}

	if qrypticClient.ExpiryTime.Before(time.Now().Add(config.QrypticClientRefetchTimeGap)) {
		return false, qrypticClient
	}

	return true, qrypticClient
}

func getGatewayClient(uuid string) (models.WGClientConfig, error) {
	log := logger.Default()
	ifClientExisting, oldQrypticClient := clientExisiting(uuid)
	if ifClientExisting {
		return oldQrypticClient, nil
	}
	baseUrl, _ := storage.GetBaseUrl()
	authToken, _ := storage.GetAuthToken()
	qrypticClient := client.NewQrypticClient(baseUrl, authToken)
	statusCode, clientConfig, err := qrypticClient.GetGatewayClient(uuid)
	if err != nil {
		log.Error(err.Error())
		return *clientConfig, err
	}
	if statusCode == http.StatusOK {
		storage.SetQrypticClient(uuid, *clientConfig)
		return *clientConfig, nil

	} else if statusCode == http.StatusUnauthorized {
		log.Error("Please authenticate ...")
		return *clientConfig, err
	} else {
		log.Error("Server Issue ...")
		return *clientConfig, err
	}
}

func connectToGateway(uuid, name string) {
	log := logger.Default()
	clientConfig, err := getGatewayClient(uuid)
	if err != nil {
		return
	}
	err = wg.ApplyConfig(clientConfig)
	if err != nil {
		log.Error(err.Error())
	}
	storage.SetConnectedToGateway(uuid, name)
	uuid, name, exists := storage.GetConnectedToGateway()
	if exists {
		wgclientConfig, _ := storage.GetQrypticClient(uuid)
		fmt.Printf("Connected to %s gateway at %s:%d\n", name, wgclientConfig.WGClientPeerConfig.VpnGatewayIP, wgclientConfig.WGClientPeerConfig.VpnGatewayPort)
	} else {
		fmt.Println("Connection failed !!")
	}
}

func promptGatewaySelect(pc models.PromptContent, gateways []models.GatewayResponse) (string, int) {
	log := logger.Default()
	items := []string{}
	for _, gateway := range gateways {
		items = append(items, gateway.Name)
	}
	prompt := promptui.Select{
		Label: pc.Label,
		Items: items,
	}
	index, result, err := prompt.Run()

	if err != nil {
		log.Error("Gateway Selection failed %v", "error", err.Error())
		os.Exit(1)
	}
	return result, index
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
