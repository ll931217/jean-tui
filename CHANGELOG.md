# Changelog

All notable changes to jean will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added - OpenAI-Compatible API Provider System

#### AI Integration Migration
- **BREAKING CHANGE**: Replaced OpenRouter with OpenAI-compatible API provider system
- **Multi-Provider Support**: Configure multiple AI provider profiles
  - OpenAI (official API)
  - Azure OpenAI Service
  - Custom endpoints (Ollama, vLLM, LM Studio, local LLM servers)
- **Provider Profiles**: Save and manage multiple AI configurations
  - Profile name, type, base URL, API key, and model
  - Create, edit, delete, and test provider profiles via TUI
  - Set active provider for AI features
  - Configure fallback provider for automatic failover
- **Test Connection**: Validate API credentials and connectivity before use
- **Custom AI Prompts**: Override default prompts globally
  - Custom prompts for commit messages, branch names, PR content
  - Supports placeholders: `{status}`, `{diff}`, `{branch}`, `{log}`

#### New TUI Features
- AI Providers modal (`s` → AI Providers)
  - List all provider profiles with status indicators
  - Create new provider profiles with type selection
  - Edit existing provider profiles
  - Delete provider profiles
  - Test provider connection
  - Set active and fallback providers
- Provider Profile Edit Modal
  - Name field for display name
  - Type selector (openai, azure, custom)
  - Base URL field with provider-specific defaults
  - API key field with secure storage
  - Model selection with predefined options

#### Configuration Changes
- **Removed** `openrouter_api_key` and `openrouter_model` fields
- **Added** `ai_provider` configuration per repository:
  ```json
  {
    "ai_provider": {
      "profiles": {
        "profile-name": {
          "name": "Display Name",
          "type": "openai|azure|custom",
          "base_url": "https://...",
          "api_key": "...",
          "model": "..."
        }
      },
      "active_profile": "profile-name",
      "fallback_profile": "fallback-profile-name"
    }
  }
  ```
- **Added** global `ai_prompts` configuration:
  ```json
  {
    "ai_prompts": {
      "commit_message": "Custom prompt...",
      "branch_name": "Custom prompt...",
      "pr_content": "Custom prompt..."
    }
  }
  ```

#### Internal Changes
- New `openai/` package:
  - `client.go`: Generic OpenAI-compatible HTTP client
  - `provider.go`: Provider types and configuration
  - `models.go`: Predefined model lists
  - `prompts.go`: Default prompt templates
  - `errors.go`: Structured error handling
- Updated `config/config.go`:
  - Provider profile CRUD operations
  - Active/fallback provider management
  - 93.5% test coverage for new provider system
- Updated all AI generation to use new provider system
  - Commit message generation with fallback
  - Branch name generation with fallback
  - PR content generation with fallback

### Migration Guide

If you were using the old OpenRouter configuration:

1. **Open jean** and press `s` → AI Providers
2. **Create a new provider profile**:
   - Choose type: `custom` (for OpenRouter compatibility)
   - Set base URL: `https://openrouter.ai/api/v1`
   - Enter your OpenRouter API key
   - Select your preferred model
3. **Set as active provider** to restore AI functionality

Or migrate to OpenAI directly:
1. Create a new profile with type `openai`
2. Set base URL to `https://api.openai.com/v1`
3. Enter your OpenAI API key
4. Select model (e.g., `gpt-4`, `gpt-3.5-turbo`)

### Fixed
- Improved error handling for AI API failures
- Better fallback logic when primary provider fails
- More informative error messages for provider configuration issues

### Changed
- AI Settings modal now opens AI Providers management
- Provider profiles are repository-specific (not global)
- Default prompts are now customizable globally

## [0.1.15] - Previous Release
- Previous features and changes
