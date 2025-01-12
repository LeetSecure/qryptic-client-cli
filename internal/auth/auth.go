package auth

import (
	"net/http"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/leetsecure/qryptic-client-cli/internal/client"
	"github.com/leetsecure/qryptic-client-cli/internal/config"
	"github.com/spf13/viper"
)

func IsAuthTokenValid() bool {
	authToken := viper.GetViper().GetString(config.AuthToken)
	if authToken == "" {
		return false
	}
	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(authToken, claims, nil)
	expiryTime, _ := claims.GetExpirationTime()
	currentTime := time.Now()
	return !expiryTime.Time.Before(currentTime)
}

func IsURL(urlToCheck string) bool {
	_, err := url.ParseRequestURI(urlToCheck)
	return err == nil
}

func IsBaseUrlHealthy() bool {
	// log := logger.Default()
	baseUrl := viper.GetViper().GetString(config.BaseUrl)

	qrypticClient := client.NewQrypticClient(baseUrl, "")
	statusCode, response, err := qrypticClient.ControllerHealthCheck()
	if err != nil {
		// log.Error(err.Error())
		return false
	}
	if statusCode == http.StatusOK && response.Success {
		return true
	}
	return false
}

func GenerateCodeVerifierAndChallenge() (string, string) {
	return "", ""
}
