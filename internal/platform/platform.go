package platform

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func GetPlatform() string {
	switch runtime.GOOS {
	case "linux":
		return "linux"
	case "darwin":
		return "macos"
	case "windows":
		return "windows"
	default:
		return "unknown"
	}
}

func GetDefaultInterfaceName() string {
	switch runtime.GOOS {
	case "linux":
		return "wg0"
	case "darwin":
		return "utun0"
	case "windows":
		return "wg0"
	default:
		return "unknown"
	}
}

func PlatformSpecificOperation() {
	platform := GetPlatform()
	switch platform {
	case "linux":
		fmt.Println("Running on Linux")
		// Linux-specific logic here
	case "macos":
		fmt.Println("Running on macOS")
		// macOS-specific logic here
	case "windows":
		fmt.Println("Running on Windows")
		// Windows-specific logic here
	default:
		fmt.Println("Unsupported platform")
	}
}

func OpenURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default: // "linux", "freebsd", "openbsd", "netbsd"
		// Check if running under WSL
		if isWSL() {
			// Use 'cmd.exe /c start' to open the URL in the default Windows browser
			cmd = "cmd.exe"
			args = []string{"/c", "start", url}
		} else {
			// Use xdg-open on native Linux environments
			cmd = "xdg-open"
			args = []string{url}
		}
	}
	if len(args) > 1 {
		// args[0] is used for 'start' command argument, to prevent issues with URLs starting with a quote
		args = append(args[:1], append([]string{""}, args[1:]...)...)
	}
	return exec.Command(cmd, args...).Start()
}

// isWSL checks if the Go program is running inside Windows Subsystem for Linux
func isWSL() bool {
	releaseData, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(releaseData)), "microsoft")
}

func GetConfigDirectory() string {
	switch runtime.GOOS {
	case "linux":
		return "/etc/wireguard"
	case "darwin":
		return "/etc/wireguard"
	case "windows":
		return filepath.Join(os.Getenv("LOCALAPPDATA"), "WireGuard")
	default:
		return "/etc/wireguard"
	}
}
