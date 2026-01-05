package openai

// ProviderType represents the type of OpenAI-compatible provider
type ProviderType string

const (
	// ProviderOpenAI is the official OpenAI API
	ProviderOpenAI ProviderType = "openai"
	// ProviderAzure is Azure OpenAI Service
	ProviderAzure ProviderType = "azure"
	// ProviderCustom is a custom OpenAI-compatible endpoint (e.g., local server)
	ProviderCustom ProviderType = "custom"
)

// DefaultBaseURL returns the default base URL for a given provider type
func DefaultBaseURL(providerType ProviderType) string {
	switch providerType {
	case ProviderOpenAI:
		return "https://api.openai.com/v1"
	case ProviderAzure:
		// Azure requires custom URL with resource and deployment
		return ""
	case ProviderCustom:
		// Custom providers must specify their own URL
		return ""
	default:
		return "https://api.openai.com/v1"
	}
}

// IsValidProvider checks if a provider type is valid
func IsValidProvider(providerType ProviderType) bool {
	switch providerType {
	case ProviderOpenAI, ProviderAzure, ProviderCustom:
		return true
	default:
		return false
	}
}

// ProviderConfig holds configuration for an OpenAI-compatible provider
type ProviderConfig struct {
	Type     ProviderType `json:"type"`
	BaseURL  string       `json:"base_url"`
	APIKey   string       `json:"api_key"`
	Model    string       `json:"model"`
}

// NewProviderConfig creates a new provider config with defaults
func NewProviderConfig(providerType ProviderType) *ProviderConfig {
	baseURL := DefaultBaseURL(providerType)
	return &ProviderConfig{
		Type:    providerType,
		BaseURL: baseURL,
	}
}

// GetEffectiveBaseURL returns the base URL to use, with fallback to default
func (p *ProviderConfig) GetEffectiveBaseURL() string {
	if p.BaseURL != "" {
		return p.BaseURL
	}
	return DefaultBaseURL(p.Type)
}
