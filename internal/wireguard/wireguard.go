package wireguard

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/leetsecure/qryptic-client-cli/internal/models"
)

const wgConfigTemplate = `[Interface]
PrivateKey = {{.InterfaceConfig.ClientPrivateKey}}
Address = {{.InterfaceConfig.AllowedIpAddress}}
DNS = {{.InterfaceConfig.DnsServer}}

[Peer]
PublicKey = {{.PeerConfig.ServerPublicKey}}
AllowedIPs = {{range $index, $ip := .PeerConfig.AllowedIPs}}{{if $index}},{{end}}{{$ip}}{{end}}
Endpoint = {{.PeerConfig.VpnGatewayIP}}:{{.PeerConfig.VpnGatewayPort}}
PersistentKeepalive = {{.PeerConfig.PersistantAlive}}
`

// WireGuardManager manages WireGuard configurations and connections.
type WireGuardManager struct {
	ConfigDir  string
	ConfigPath string
	Interface  string
}

// NewWireGuardManager initializes a new WireGuardManager.
func NewWireGuardManager(configDir, interfaceName string) *WireGuardManager {
	return &WireGuardManager{
		ConfigDir:  configDir,
		ConfigPath: filepath.Join(configDir, "wg0.conf"),
		Interface:  interfaceName,
	}
}

// ApplyConfig applies the WireGuard configuration.
func (wg *WireGuardManager) ApplyConfig(clientConfig models.WGClientConfig) error {
	// Stop any existing VPN
	if err := wg.StopVPN(); err != nil {
		return fmt.Errorf("failed to stop existing VPN: %w", err)
	}

	// Generate configuration file
	if err := wg.generateConfig(clientConfig); err != nil {
		return fmt.Errorf("failed to generate config: %w", err)
	}

	// Start VPN
	if err := wg.StartVPN(); err != nil {
		return fmt.Errorf("failed to start VPN: %w", err)
	}

	return nil
}

// generateConfig generates the WireGuard configuration file from the template.
func (wg *WireGuardManager) generateConfig(clientConfig models.WGClientConfig) error {
	tmpl, err := template.New("wgConfig").Parse(wgConfigTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"InterfaceConfig": clientConfig.WGClientInterfaceConfig,
		"PeerConfig":      clientConfig.WGClientPeerConfig,
	})
	// fmt.Println(buf.String())
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Ensure configuration directory exists
	if _, err := os.Stat(wg.ConfigDir); os.IsNotExist(err) {
		if err := os.MkdirAll(wg.ConfigDir, 0700); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
	}

	// Write configuration file
	if err := os.WriteFile(wg.ConfigPath, buf.Bytes(), 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Println("Config Created")

	return nil
}

// StartVPN brings up the WireGuard interface using the configuration.
func (wg *WireGuardManager) StartVPN() error {
	cmd := exec.Command("wg-quick", "up", wg.ConfigPath)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start VPN: %w", err)
	}
	return nil
}

// StopVPN brings down the WireGuard interface.
func (wg *WireGuardManager) StopVPN() error {
	_, err := os.Stat(wg.ConfigPath)
	if err != nil {
		fmt.Println("Config not present")
		return nil
	}
	isrunning, _, _ := wg.CheckStatus()
	if !isrunning {
		return nil
	}
	cmd := exec.Command("wg-quick", "down", wg.ConfigPath)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop VPN: %w", err)
	}
	return nil
}

// CheckStatus checks the status of the WireGuard interface.
func (wg *WireGuardManager) CheckStatus() (bool, string, error) {
	// cmd := exec.Command("wg", "show", wg.Interface)
	cmd := exec.Command("wg", "show")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return false, "", fmt.Errorf("failed to check VPN status: %w", err)
	}
	return len(output) > 0, string(output), nil
}

// Cleanup removes the current configuration file.
func (wg *WireGuardManager) Cleanup() error {
	err := wg.StopVPN()
	if err != nil {
		return err
	}

	if err := os.Remove(wg.ConfigPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove config file: %w", err)
	}
	return nil
}
