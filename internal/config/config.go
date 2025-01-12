package config

import "time"

var BaseUrl = "baseUrl"
var AuthForUrl = "authForUrl"
var AuthToken = "authToken"
var ConnectedToGateway = "connectedToGateway"
var ConnectedToGatewayUuid = "connectedToGateway.uuid"
var ConnectedToGatewayName = "connectedToGateway.name"

var ConfigFileName = ".qryptic"
var ConfigFileType = "yaml"
var QrypticClientRefetchTimeGap = 30 * time.Minute
var IsWireguardSetupCompleted = "isWireguardSetupCompleted"
