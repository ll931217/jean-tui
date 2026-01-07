package util

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// PortParser parses port numbers from various configuration files
type PortParser struct {
	worktreePath string
}

// NewPortParser creates a new port parser for a worktree
func NewPortParser(worktreePath string) *PortParser {
	return &PortParser{worktreePath: worktreePath}
}

// ParsePorts parses all port numbers from the worktree
func (p *PortParser) ParsePorts() []int {
	ports := make(map[int]bool)

	// Parse package.json (for Node.js projects)
	if pkgPorts := p.parsePackageJSON(); len(pkgPorts) > 0 {
		for _, port := range pkgPorts {
			ports[port] = true
		}
	}

	// Parse .env files
	if envPorts := p.parseEnvFile(); len(envPorts) > 0 {
		for _, port := range envPorts {
			ports[port] = true
		}
	}

	// Parse docker-compose.yml
	if composePorts := p.parseDockerCompose(); len(composePorts) > 0 {
		for _, port := range composePorts {
			ports[port] = true
		}
	}

	// Convert map to sorted slice
	result := make([]int, 0, len(ports))
	for port := range ports {
		result = append(result, port)
	}

	return result
}

// parsePackageJSON parses ports from package.json scripts
// Looks for common patterns like "PORT=3000", "--port 8080", etc.
func (p *PortParser) parsePackageJSON() []int {
	pkgPath := filepath.Join(p.worktreePath, "package.json")
	data, err := os.ReadFile(pkgPath)
	if err != nil {
		return nil
	}

	var packageJSON struct {
		Scripts map[string]string `json:"scripts"`
	}

	if err := json.Unmarshal(data, &packageJSON); err != nil {
		return nil
	}

	ports := make(map[int]bool)

	// Common port names to look for in scripts
	portPatterns := []string{
		"PORT", "port", "PORT=", "--port", "-p",
		"HOST_PORT", "SERVER_PORT", "API_PORT",
		"VITE_PORT", "NEXT_PORT", "REACT_PORT",
		"3000", "3001", "4000", "5000", "8000", "8080",
	}

	for _, script := range packageJSON.Scripts {
		// Look for port patterns in script commands
		for _, pattern := range portPatterns {
			if strings.Contains(script, pattern) {
				// Try to extract port number after the pattern
				port := p.extractPortFromString(script, pattern)
				if port > 0 {
					ports[port] = true
				}
			}
		}
	}

	// Default ports for common frameworks
	if len(ports) == 0 {
		for _, script := range packageJSON.Scripts {
			// Check for common framework scripts
			if strings.Contains(script, "next") || strings.Contains(script, "remix") {
				ports[3000] = true
			} else if strings.Contains(script, "vite") || strings.Contains(script, "react") {
				ports[5173] = true
			} else if strings.Contains(script, "vue") {
				ports[5173] = true
			} else if strings.Contains(script, "svelte") {
				ports[5173] = true
			} else if strings.Contains(script, "angular") {
				ports[4200] = true
			}
		}
	}

	result := make([]int, 0, len(ports))
	for port := range ports {
		result = append(result, port)
	}

	return result
}

// parseEnvFile parses ports from .env files
// Looks for PORT, HOST_PORT, SERVER_PORT, API_PORT, etc.
func (p *PortParser) parseEnvFile() []int {
	envPaths := []string{".env", ".env.local", ".env.development", ".env.production"}
	ports := make(map[int]bool)

	for _, envFile := range envPaths {
		envPath := filepath.Join(p.worktreePath, envFile)
		data, err := os.ReadFile(envPath)
		if err != nil {
			continue
		}

		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			// Look for PORT variables
			if strings.HasPrefix(line, "PORT=") || strings.HasPrefix(line, "HOST_PORT=") ||
				strings.HasPrefix(line, "SERVER_PORT=") || strings.HasPrefix(line, "API_PORT=") ||
				strings.HasPrefix(line, "VITE_PORT=") || strings.HasPrefix(line, "NEXT_PORT=") {

				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					portStr := strings.Trim(parts[1], `"'`)
					if port, err := strconv.Atoi(portStr); err == nil && port > 0 && port < 65536 {
						ports[port] = true
					}
				}
			}
		}
	}

	result := make([]int, 0, len(ports))
	for port := range ports {
		result = append(result, port)
	}

	return result
}

// parseDockerCompose parses ports from docker-compose.yml files
func (p *PortParser) parseDockerCompose() []int {
	composePaths := []string{"docker-compose.yml", "docker-compose.yaml", "compose.yml", "compose.yaml"}
	ports := make(map[int]bool)

	for _, composeFile := range composePaths {
		composePath := filepath.Join(p.worktreePath, composeFile)
		data, err := os.ReadFile(composePath)
		if err != nil {
			continue
		}

		content := string(data)
		lines := strings.Split(content, "\n")

		for _, line := range lines {
			line = strings.TrimSpace(line)

			// Look for port mappings like "3000:3000" or "- 3000" or "ports: - '3000:3000'"
			if strings.Contains(line, ":") && (strings.Contains(line, "3000") ||
				strings.Contains(line, "3001") || strings.Contains(line, "4000") ||
				strings.Contains(line, "5000") || strings.Contains(line, "8000") ||
				strings.Contains(line, "8080") || strings.Contains(line, "9000")) {

				// Extract port numbers from the line
				port := p.extractPortFromDockerComposeLine(line)
				if port > 0 {
					ports[port] = true
				}
			}
		}
	}

	result := make([]int, 0, len(ports))
	for port := range ports {
		result = append(result, port)
	}

	return result
}

// extractPortFromString extracts a port number from a string after a pattern
func (p *PortParser) extractPortFromString(s, pattern string) int {
	idx := strings.Index(s, pattern)
	if idx == -1 {
		return 0
	}

	// Get the substring after the pattern
	substr := s[idx+len(pattern):]
	substr = strings.TrimSpace(substr)

	// Look for a port number (4 digits)
	for i := 0; i < len(substr)-3; i++ {
		c := substr[i]
		if c >= '0' && c <= '9' {
			// Try to parse a 4-5 digit number as port
			portStr := ""
			for j := i; j < len(substr) && j < i+5; j++ {
				if substr[j] >= '0' && substr[j] <= '9' {
					portStr += string(substr[j])
				} else {
					break
				}
			}
			if port, err := strconv.Atoi(portStr); err == nil && port > 0 && port < 65536 {
				return port
			}
		}
	}

	return 0
}

// extractPortFromDockerComposeLine extracts a port from a docker-compose line
func (p *PortParser) extractPortFromDockerComposeLine(line string) int {
	// Remove quotes
	line = strings.ReplaceAll(line, "'", "")
	line = strings.ReplaceAll(line, `"`, "")
	line = strings.ReplaceAll(line, " ", "")

	// Look for patterns like "3000:3000" or "-3000" or "ports:3000"
	if strings.Contains(line, ":") {
		parts := strings.Split(line, ":")
		for _, part := range parts {
			if port, err := strconv.Atoi(part); err == nil && port > 0 && port < 65536 {
				return port
			}
		}
	}

	// Look for standalone port numbers
	words := strings.Fields(line)
	for _, word := range words {
		if port, err := strconv.Atoi(word); err == nil && port > 0 && port < 65536 {
			return port
		}
	}

	return 0
}

// GetPortURLs returns clickable localhost URLs for the given ports
func GetPortURLs(ports []int) []string {
	urls := make([]string, 0, len(ports))
	for _, port := range ports {
		urls = append(urls, fmt.Sprintf("http://localhost:%d", port))
	}
	return urls
}

// GetPortServiceName returns a likely service name for a given port
func GetPortServiceName(port int) string {
	serviceNames := map[int]string{
		3000:  "Dev Server",
		3001:  "Dev Server (Alt)",
		4200:  "Angular",
		5173:  "Vite/React/Vue/Svelte",
		8000:  "Python/Go Server",
		8080:  "API Server",
		9000:  "Proxy Server",
		5432:  "PostgreSQL",
		3306:  "MySQL",
		6379:  "Redis",
		27017: "MongoDB",
	}

	if name, ok := serviceNames[port]; ok {
		return name
	}

	return fmt.Sprintf("Port %d", port)
}
