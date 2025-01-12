package config

import (
	"os"

	"github.com/leetsecure/qryptic-client-cli/internal/models"
	"github.com/spf13/viper"
)

type Storage struct {
	vip *viper.Viper
}

func NewStorage(vipp *viper.Viper) (*Storage, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	viper.AddConfigPath(home)
	viper.SetConfigType(ConfigFileType)
	viper.SetConfigName(ConfigFileName)
	viper.SafeWriteConfig()
	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return &Storage{
		vip: vipp,
	}, nil
}

func (s *Storage) GetBaseUrl() (string, bool) {
	baseUrl := s.vip.GetString(BaseUrl)
	if baseUrl == "" {
		return "", false
	}
	return baseUrl, true
}

func (s *Storage) SetBaseUrl(baseUrl string) error {
	s.vip.Set(BaseUrl, baseUrl)
	return s.vip.WriteConfig()
}

func (s *Storage) GetAuthToken() (string, bool) {
	authToken := s.vip.GetString(AuthToken)
	if authToken == "" {
		return "", false
	}
	return authToken, true
}

func (s *Storage) SetAuthToken(authToken string) error {
	s.vip.Set(AuthToken, authToken)
	return s.vip.WriteConfig()
}

func (s *Storage) ClearAuthToken() error {
	err := s.SetAuthToken("")
	if err != nil {
		return err
	}
	return s.SetAuthForUrl("")
}

func (s *Storage) GetAuthForUrl() (string, bool) {
	authForUrl := s.vip.GetString(AuthForUrl)
	if authForUrl == "" {
		return "", false
	}
	return authForUrl, true
}

func (s *Storage) SetAuthForUrl(authForUrl string) error {
	s.vip.Set(AuthForUrl, authForUrl)
	return s.vip.WriteConfig()
}

func (s *Storage) GetConnectedToGateway() (string, string, bool) {
	uuid := s.vip.GetString(ConnectedToGatewayUuid)
	name := s.vip.GetString(ConnectedToGatewayName)
	if uuid == "" || name == "" {
		return "", "", false
	}
	return uuid, name, true
}

func (s *Storage) SetConnectedToGateway(uuid, name string) error {
	s.vip.Set(ConnectedToGatewayUuid, uuid)
	s.vip.Set(ConnectedToGatewayName, name)
	return s.vip.WriteConfig()
}

func (s *Storage) ClearConnectedToGateway() error {
	return s.SetConnectedToGateway("", "")
}

func (s *Storage) ClearConfig() error {
	configFilePath := viper.GetViper().ConfigFileUsed()
	err := os.Remove(configFilePath)
	return err
}

func (s *Storage) GetQrypticClient(uuid string) (models.WGClientConfig, error) {
	var qrypticClient models.WGClientConfig
	err := viper.GetViper().UnmarshalKey(uuid, &qrypticClient)
	return qrypticClient, err
}

func (s *Storage) SetQrypticClient(uuid string, clientConfig models.WGClientConfig) error {
	viper.GetViper().Set(uuid, clientConfig)
	return viper.GetViper().WriteConfig()
}

func (s *Storage) GetWireguardSetup() bool {
	return s.vip.GetBool(IsWireguardSetupCompleted)
}

func (s *Storage) SetWireguardSetup(isWireguardSetupCompleted bool) error {
	s.vip.Set(IsWireguardSetupCompleted, isWireguardSetupCompleted)
	return s.vip.WriteConfig()
}
