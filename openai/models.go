package openai

// Predefined model constants for common OpenAI-compatible models
const (
	// GPT-4 Models
	ModelGPT4           = "gpt-4"
	ModelGPT4Turbo      = "gpt-4-turbo"
	ModelGPT4o          = "gpt-4o"
	ModelGPT4oMini      = "gpt-4o-mini"

	// GPT-3.5 Models
	ModelGPT35Turbo     = "gpt-3.5-turbo"

	// Azure-specific models
	ModelAzureGPT4      = "gpt-4"
	ModelAzureGPT35Turbo = "gpt-35-turbo"
)

// DefaultModels returns a list of predefined models for a given provider type
func DefaultModels(providerType string) []string {
	switch providerType {
	case "openai":
		return []string{
			ModelGPT4,
			ModelGPT4Turbo,
			ModelGPT4o,
			ModelGPT4oMini,
			ModelGPT35Turbo,
		}
	case "azure":
		return []string{
			ModelAzureGPT4,
			ModelAzureGPT35Turbo,
		}
	case "custom":
		// Custom providers (e.g., Ollama) use their own model names
		return []string{}
	default:
		return []string{
			ModelGPT4,
			ModelGPT35Turbo,
		}
	}
}

// IsValidModel checks if a model is in the predefined list for a provider
func IsValidModel(providerType, model string) bool {
	for _, m := range DefaultModels(providerType) {
		if m == model {
			return true
		}
	}
	// Allow custom model names for any provider
	return model != ""
}
