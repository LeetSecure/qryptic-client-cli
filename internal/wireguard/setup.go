package wireguard

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/leetsecure/qryptic-client-cli/internal/config"
	"github.com/spf13/viper"
)

// Supported package managers for Linux
var linuxPackageManagers = []string{"apt-get", "yum", "dnf", "zypper", "pacman"}

// SetupWireGuard ensures tools and directories are properly set up.
func SetupWireGuard() error {
	isSetupCompleted := viper.GetViper().GetBool(config.IsWireguardSetupCompleted)
	if isSetupCompleted {
		fmt.Println("Setup already completed")
		return nil
	}
	fmt.Println("Setting up WireGuard...")

	// Check and install tools if necessary
	if err := checkAndInstallTools(); err != nil {
		return fmt.Errorf("tool setup failed: %w", err)
	}

	// Check and create configuration directory if necessary
	if err := ensureConfigDirectory(); err != nil {
		return fmt.Errorf("directory setup failed: %w", err)
	}

	fmt.Println("WireGuard setup completed successfully.")
	viper.GetViper().Set(config.IsWireguardSetupCompleted, true)
	viper.GetViper().WriteConfig()
	return nil
}

// checkAndInstallTools ensures the required tools are installed.
func checkAndInstallTools() error {
	fmt.Println("Checking for required tools...")

	switch runtime.GOOS {
	case "linux":
		return checkAndInstallLinux()
	case "darwin":
		return checkAndInstallMacOS()
	case "windows":
		return checkAndInstallWindows()
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// checkAndInstallLinux installs tools on Linux if not already installed.
func checkAndInstallLinux() error {
	if err := checkBinary("wg"); err != nil {
		fmt.Println("WireGuard tools not found. Installing...")
		pkgManager, err := detectLinuxPackageManager()
		if err != nil {
			return fmt.Errorf("failed to detect package manager: %w", err)
		}

		installCommand := getInstallCommand(pkgManager, "wireguard-tools")
		cmd := exec.Command("sudo", strings.Fields(installCommand)...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install WireGuard tools on Linux: %w", err)
		}
	}
	fmt.Println("WireGuard tools are installed.")
	return nil
}

// detectLinuxPackageManager detects the appropriate package manager for the system.
func detectLinuxPackageManager() (string, error) {
	for _, manager := range linuxPackageManagers {
		if _, err := exec.LookPath(manager); err == nil {
			fmt.Printf("Detected package manager: %s\n", manager)
			return manager, nil
		}
	}
	return "", fmt.Errorf("no supported package manager found")
}

// getInstallCommand returns the installation command for the given package manager.
func getInstallCommand(pkgManager, packageName string) string {
	switch pkgManager {
	case "apt-get":
		return fmt.Sprintf("apt-get install -y %s", packageName)
	case "yum":
		return fmt.Sprintf("yum install -y %s", packageName)
	case "dnf":
		return fmt.Sprintf("dnf install -y %s", packageName)
	case "zypper":
		return fmt.Sprintf("zypper install -y %s", packageName)
	case "pacman":
		return fmt.Sprintf("pacman -S --noconfirm %s", packageName)
	default:
		return ""
	}
}

// checkAndInstallMacOS installs tools on macOS if not already installed.
func checkAndInstallMacOS() error {
	fmt.Println("Checking architecture on macOS...")
	arch, err := detectArchitecture()
	if err != nil {
		return fmt.Errorf("failed to detect architecture: %w", err)
	}
	fmt.Println("Detected architecture: " + arch)

	if err := checkBinary("wg"); err != nil {
		fmt.Println("Installing WireGuard tools on macOS...")
		var cmd *exec.Cmd
		if arch == "arm64" {
			fmt.Println("Apple Silicon detected. Using Homebrew for arm64.")
			cmd = exec.Command("/opt/homebrew/bin/brew", "install", "wireguard-tools")
		} else {
			fmt.Println("Intel detected. Using Homebrew for x86_64.")
			cmd = exec.Command("brew", "install", "wireguard-tools")
		}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install WireGuard tools on macOS: %w", err)
		}
	}
	fmt.Println("WireGuard tools are installed.")
	return nil
}

// checkAndInstallWindows guides the user to install WireGuard on Windows.
func checkAndInstallWindows() error {
	if err := checkBinary("wireguard.exe"); err != nil {
		fmt.Println("WireGuard is not installed. Please download and install it from: https://www.wireguard.com/install/")
		return fmt.Errorf("WireGuard not installed")
	}
	fmt.Println("WireGuard tools are installed.")
	return nil
}

// detectArchitecture determines the system architecture.
func detectArchitecture() (string, error) {
	cmd := exec.Command("uname", "-m")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to run uname -m: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// checkBinary checks if a binary exists in the system PATH.
func checkBinary(binary string) error {
	if _, err := exec.LookPath(binary); err != nil {
		return fmt.Errorf("%s not found in PATH", binary)
	}
	return nil
}

// ensureConfigDirectory ensures the required configuration directory exists.
func ensureConfigDirectory() error {
	configDir := getConfigDirectory()
	fmt.Printf("Checking if configuration directory %s exists...\n", configDir)

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		fmt.Printf("Creating configuration directory: %s\n", configDir)
		if err := os.MkdirAll(configDir, 0700); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}
	fmt.Println("Configuration directory is set up.")
	return nil
}

// getConfigDirectory returns the appropriate configuration directory based on the platform.
func getConfigDirectory() string {
	switch runtime.GOOS {
	case "windows":
		return "C:\\Program Files\\WireGuard"
	case "darwin":
		return "/etc/wireguard"
	default:
		return "/etc/wireguard"
	}
}
