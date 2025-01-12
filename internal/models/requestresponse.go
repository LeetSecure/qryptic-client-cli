package models

import "time"

// ErrorResponse represents an error response from the API.
// type ErrorResponse struct {
// 	Message string `json:"message"`
// 	Code    int    `json:"code"`
// }

type EmailPasswordLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	AuthToken string `json:"authToken"`
	Error     string `json:"error"`
	Message   string `json:"message"`
}

type GatewayResponse struct {
	Domain          string `json:"domain"`
	IpAddress       string `json:"ipAddress"`
	Name            string `json:"name"`
	Port            int    `json:"port"`
	ServerPublicKey string `json:"serverPublicKey"`
	Uuid            string `json:"uuid"`
}

type HealthCheckResponse struct {
	Success bool `json:"success"`
}

type PromptContent struct {
	ErrorMsg string
	Label    string
}

type WGClientInterfaceConfig struct {
	ClientPrivateKey string `json:"privateKey"`
	AllowedIpAddress string `json:"ipAddress"`
	DnsServer        string `json:"dnsServer"`
}

type WGClientPeerConfig struct {
	AllowedIPs       []string `json:"allowedIPs"`
	ServerPublicKey  string   `json:"publicKey"`
	PresharedKey     string   `json:"presharedKey"`
	PersistantAlive  int      `json:"persistantAlive"`
	VpnGatewayDomain string   `json:"vpnGatewayDomain"`
	VpnGatewayIP     string   `json:"vpnGatewayIP"`
	VpnGatewayPort   int      `json:"vpnGatewayPort"`
}

type WGClientConfig struct {
	ClientUuid              string                  `json:"clientUuid"`
	WGClientInterfaceConfig WGClientInterfaceConfig `json:"clientInterfaceConfig"`
	WGClientPeerConfig      WGClientPeerConfig      `json:"clientPeerConfig"`
	ExpiryTime              time.Time               `json:"expiryTime"`
}

/*


{
	"clientInterfaceConfig": {
	  "dnsServer": "string",
	  "ipAddress": "string",
	  "privateKey": "string"
	},
	"clientPeerConfig": {
	  "allowedIPs": [
		"string"
	  ],
	  "persistantAlive": 0,
	  "presharedKey": "string",
	  "publicKey": "string",
	  "vpnGatewayDomain": "string",
	  "vpnGatewayIP": "string",
	  "vpnGatewayPort": 0
	},
	"clientUuid": "string",
	"expiryTime": "string"
  }
*/
