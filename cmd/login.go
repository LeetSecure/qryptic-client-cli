/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/leetsecure/qryptic-client-cli/internal/auth"
	"github.com/leetsecure/qryptic-client-cli/internal/client"
	"github.com/leetsecure/qryptic-client-cli/internal/config"
	"github.com/leetsecure/qryptic-client-cli/internal/logger"
	"github.com/leetsecure/qryptic-client-cli/internal/models"
	"github.com/leetsecure/qryptic-client-cli/internal/platform"
	"github.com/leetsecure/qryptic-client-cli/internal/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var URL string
var ForceLogin bool

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login into Qryptic",
	Long:  `Authenticate to your organisation's hosted Qryptic Controller`,
	Run:   loginExecute,
}

func loginExecute(cmd *cobra.Command, args []string) {
	log := logger.Default()
	baseUrl, _ := storage.GetBaseUrl()
	isValidUrl := auth.IsURL(baseUrl)
	if !isValidUrl {
		log.Error("Please make sure the url is in format : http[s]://<domain/subdomain> \n Example : https://qryptic.leetsecure.com")
		return
	}
	if baseUrl[len(baseUrl)-1:] == "/" {
		baseUrl = baseUrl[:len(baseUrl)-1]
		storage.SetBaseUrl(baseUrl)
	}
	if !auth.IsBaseUrlHealthy() {
		log.Error("Given url is unhealthy. Check again if the url is correct. If URL is correct, check with Admin if the Qryptic service is running", "url", baseUrl)
		return
	}

	authForUrl, _ := storage.GetAuthForUrl()
	if !ForceLogin && (baseUrl == authForUrl) {
		isValid := auth.IsAuthTokenValid()
		if isValid {
			log.Info("Already authenticated")
			return
		}
	}
	selectLoginMethod()
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVarP(&URL, "url", "u", "", "Custom url for your organisation's qryptic controller")
	loginCmd.MarkFlagRequired("url")
	viper.GetViper().BindPFlag(config.BaseUrl, loginCmd.Flags().Lookup("url"))
	loginCmd.Flags().BoolVarP(&ForceLogin, "force", "f", false, "Force new login to replace the existing auth credentials with new one")
}

func promptLoginMethodSelect(pc models.PromptContent) (string, int) {
	log := logger.Default()
	items := []string{"Email & Password", "Google SSO"}
	prompt := promptui.Select{
		Label: pc.Label,
		Items: items,
	}
	index, result, err := prompt.Run()

	if err != nil {
		log.Error(fmt.Sprint("Login Failed :", err.Error()))
		os.Exit(1)
	}

	return result, index
}

func selectLoginMethod() {
	log := logger.Default()
	loginMethodPromptContent := models.PromptContent{
		ErrorMsg: "Please select a valid login method.",
		Label:    "Select a login method.",
	}
	loginMethod, index := promptLoginMethodSelect(loginMethodPromptContent)
	log.Info("The selected method is ", "method", loginMethod)
	if index == 0 {
		loginWithEmailAndPassword()
	} else if index == 1 {
		loginWithGoogleSSO()
	}
}

func loginWithGoogleSSO() {
	log := logger.Default()
	codeVerifier := utils.RandomStringGenerator(40)
	codeChallenge := utils.GetCodeChallenge(codeVerifier)
	baseUrl, _ := storage.GetBaseUrl()
	webSSOInitiateUrl :=
		fmt.Sprintf("%s/api/v1/auth/google/web/sso/initiate?code_challenge=%s", baseUrl, codeChallenge)

	err := platform.OpenURL(webSSOInitiateUrl)
	log.Info("Authenticate yourself using - ", "Link", webSSOInitiateUrl)
	if err != nil {
		log.Error(err.Error())
		return
	}
	success := fetchAuthTokenCron(codeVerifier, codeChallenge)
	if success {
		log.Info("Successfully Authenticated")
		return
	}
	log.Info("Authentication Failed. Try again ... ")
}

func fetchAuthTokenCron(codeVerifier, codeChallenge string) bool {
	log := logger.Default()
	maxDuration := 2 * time.Minute
	interval := 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), maxDuration)
	defer cancel()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	resultChan := make(chan string)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(resultChan)
				return
			case <-ticker.C:
				success, stop := fetchAuthToken(codeVerifier, codeChallenge)
				if success {
					resultChan <- "success"
				} else if stop {
					resultChan <- "failed"
				} else {
					resultChan <- "continue"
				}

			}
		}
	}()

	// Process the results
	for {
		select {
		case <-ctx.Done():
			// Timeout reached
			return false
		case result, ok := <-resultChan:
			if !ok {
				// Channel closed, end execution
				return false
			}
			// Check conditions
			switch result {
			case "success":
				return true
			case "failed":
				return false
			default:
				// Continue execution
				log.Info("Waiting for confirmation ...")
			}
		}
	}
}

func fetchAuthToken(codeVerifier, codeChallenge string) (bool, bool) {
	log := logger.Default()
	baseUrl, exists := storage.GetBaseUrl()
	if !exists {
		log.Error("Please re-authenticate. BaseURL is missing currently")
		os.Exit(1)
	}
	qrypticClient := client.NewQrypticClient(baseUrl, "")
	statuscode, authResponse, err := qrypticClient.GetWebSSOToken(codeVerifier, codeChallenge)
	if err != nil {
		log.Error(err.Error())
		return false, true
	}
	//success
	if statuscode == http.StatusOK {
		storage.SetAuthToken(authResponse.AuthToken)
		storage.SetAuthForUrl(baseUrl)
		return true, true
	}

	//continue
	if statuscode == http.StatusUnauthorized && authResponse.Error == "" {
		return false, false
	}

	//failed
	return false, true
}

func promptEmailInput(pc models.PromptContent) string {
	log := logger.Default()
	validate := func(input string) error {
		if !utils.IsValidEmailId(input) {
			return errors.New(pc.ErrorMsg)
		}
		return nil
	}
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     pc.Label,
		Templates: templates,
		Validate:  validate,
	}

	result, err := prompt.Run()
	if err != nil {
		log.Error("Login failed", "error", err.Error())
		os.Exit(1)
	}

	return result
}

func promptPasswordInput(pc models.PromptContent) string {
	log := logger.Default()
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New(pc.ErrorMsg)
		}
		return nil
	}
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     pc.Label,
		Templates: templates,
		Validate:  validate,
		Mask:      '*',
	}

	result, err := prompt.Run()
	if err != nil {
		log.Error("Login failed", "error", err.Error())
		os.Exit(1)
	}

	return result
}

func loginWithEmailAndPassword() {
	log := logger.Default()
	emailIdContent := models.PromptContent{
		ErrorMsg: "Please enter valid email id",
		Label:    "Enter your email id",
	}
	emailId := promptEmailInput(emailIdContent)

	passwordContent := models.PromptContent{
		ErrorMsg: "Please enter the password",
		Label:    "Enter your password",
	}
	password := promptPasswordInput(passwordContent)
	baseUrl, exists := storage.GetBaseUrl()
	if !exists {
		log.Error("Please re-authenticate. BaseURL is missing currently")
		os.Exit(1)
	}
	qrypticClient := client.NewQrypticClient(baseUrl, "")

	emailPasswordRequest := models.EmailPasswordLoginRequest{
		Email:    emailId,
		Password: password,
	}
	_, authResponse, err := qrypticClient.EmailPasswordLogin(emailPasswordRequest)
	if err != nil {
		log.Error(err.Error())
		log.Error("Error while login ...")
		os.Exit(1)
	}
	if authResponse.AuthToken == "" {
		log.Error("Invalid credentials")
		os.Exit(1)
	}
	storage.SetAuthToken(authResponse.AuthToken)
	storage.SetAuthForUrl(baseUrl)

	log.Info("Successfully Authenticated")
}
