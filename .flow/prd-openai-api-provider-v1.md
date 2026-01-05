---
prd:
  version: v1
  feature_name: openai-api-provider
  status: implemented
git:
  branch: feat-llm-integration
  branch_type: feature
  created_at_commit: 04b379b106a80c4e623ac7eb69455f9d36106953
  updated_at_commit: 04b379b106a80c4e623ac7eb69455f9d36106953
worktree:
  is_worktree: true
  name: feat-llm-integration
  path: /home/ll931217/GitHub/jean-tui/.git/worktrees/jean-tui.feat-llm-integration
  repo_root: /home/ll931217/GitHub/jean-tui.feat-llm-integration
metadata:
  created_at: 2026-01-04T13:30:31Z
  updated_at: 2026-01-05T14:15:02Z
  created_by: Andras Bacsai <5845193+andrasbacsai@users.noreply.github.com>
  filename: prd-openai-api-provider-v1.md
beads:
  related_issues: [jean-tui.feat-llm-integration-dmd, jean-tui.feat-llm-integration-6ac, jean-tui.feat-llm-integration-vn1, jean-tui.feat-llm-integration-uy3, jean-tui.feat-llm-integration-pct, jean-tui.feat-llm-integration-d4a, jean-tui.feat-llm-integration-wja, jean-tui.feat-llm-integration-anq, jean-tui.feat-llm-integration-sk5, jean-tui.feat-llm-integration-b1o, jean-tui.feat-llm-integration-gkz, jean-tui.feat-llm-integration-15s, jean-tui.feat-llm-integration-3qk, jean-tui.feat-llm-integration-5mt, jean-tui.feat-llm-integration-eel, jean-tui.feat-llm-integration-zfl, jean-tui.feat-llm-integration-15z]
  related_epics: [jean-tui.feat-llm-integration-gu6, jean-tui.feat-llm-integration-alm, jean-tui.feat-llm-integration-kom, jean-tui.feat-llm-integration-3f7, jean-tui.feat-llm-integration-81u]
priorities:
  enabled: true
  default: P2
  inference_method: ai_inference_with_review
  requirements:
    - id: FR-1
      text: "Replace OpenRouter client with OpenAI-compatible API client"
      priority: P0
      confidence: high
      inferred_from: "Core breaking change requirement"
      user_confirmed: true
    - id: FR-2
      text: "Support multiple providers (OpenAI, Azure, local servers, OpenRouter)"
      priority: P1
      confidence: high
      inferred_from: "Multi-provider support goal"
      user_confirmed: true
    - id: FR-3
      text: "Provider profile management with CRUD operations"
      priority: P1
      confidence: high
      inferred_from: "Provider profile configuration requirement"
      user_confirmed: true
    - id: FR-4
      text: "API key authentication per profile"
      priority: P1
      confidence: high
      inferred_from: "User selected API key authentication"
      user_confirmed: true
    - id: FR-5
      text: "Customizable base URL per provider"
      priority: P2
      confidence: high
      inferred_from: "User selected customizable base URL"
      user_confirmed: true
    - id: FR-6
      text: "Model selection: predefined list + custom model option"
      priority: P2
      confidence: high
      inferred_from: "User selected hybrid approach"
      user_confirmed: true
    - id: FR-7
      text: "Test connection feature for provider profiles"
      priority: P2
      confidence: high
      inferred_from: "User selected test connection feature"
      user_confirmed: true
    - id: FR-8
      text: "Active profile selection"
      priority: P1
      confidence: high
      inferred_from: "User selected active profile feature"
      user_confirmed: true
    - id: FR-9
      text: "Configurable fallback provider on failure"
      priority: P2
      confidence: high
      inferred_from: "User selected configurable fallback"
      user_confirmed: true
    - id: FR-10
      text: "AI commit message generation using new API"
      priority: P1
      confidence: high
      inferred_from: "Existing core feature using AI"
      user_confirmed: true
    - id: FR-11
      text: "AI branch name generation using new API"
      priority: P2
      confidence: high
      inferred_from: "Existing feature using AI"
      user_confirmed: true
    - id: FR-12
      text: "AI PR content generation using new API"
      priority: P2
      confidence: high
      inferred_from: "Existing feature using AI"
      user_confirmed: true
    - id: FR-13
      text: "Clean migration from old OpenRouter configuration"
      priority: P1
      confidence: high
      inferred_from: "User selected clean break from old config"
      user_confirmed: true
    - id: FR-14
      text: "Settings UI updates for new provider system"
      priority: P1
      confidence: high
      inferred_from: "User selected UI simplicity priority"
      user_confirmed: true
---

# PRD: OpenAI-Compatible API Provider System

## Overview

Replace the current OpenRouter-specific AI integration with a flexible, OpenAI-compatible API client system. This will support multiple providers (OpenAI, Azure OpenAI, local LLM servers, OpenRouter) through a unified interface, enabling users to configure and switch between providers seamlessly.

## Problem Statement

The current implementation uses OpenRouter as the sole AI provider for commit messages, branch names, and PR content generation. This creates vendor lock-in and prevents users from:
- Using OpenAI directly
- Leveraging Azure OpenAI Service
- Running local LLM servers (Ollama, vLLM, LM Studio)
- Using other OpenAI-compatible endpoints

## Goals

1. Eliminate vendor lock-in to OpenRouter
2. Support multiple OpenAI-compatible API providers
3. Maintain backward compatibility for existing AI features
4. Provide simple, intuitive configuration UI
5. Ensure clean, testable, extensible code architecture

## User Stories

- **As a jean user**, I want to use OpenAI directly so that I can leverage official OpenAI models
- **As a jean user**, I want to configure Azure OpenAI so that I can use my enterprise Azure subscription
- **As a jean user**, I want to use local LLM servers so that I can keep my code private and reduce API costs
- **As a jean user**, I want to save multiple provider profiles so that I can easily switch between them
- **As a jean user**, I want to test my API connection so that I can verify my credentials work before using them
- **As a jean user**, I want to configure a fallback provider so that my AI features continue working if my primary provider fails

## Functional Requirements

| ID   | Requirement                                      | Priority | Notes                    |
| ---- | ------------------------------------------------ | -------- | ------------------------ |
| FR-1 | Replace OpenRouter client with OpenAI-compatible API client | P0 | Core breaking change |
| FR-2 | Support multiple providers (OpenAI, Azure, local servers, OpenRouter) | P1 | Multi-provider goal |
| FR-3 | Provider profile management with CRUD operations | P1 | Profile config |
| FR-4 | API key authentication per profile              | P1 | Bearer token auth |
| FR-5 | Customizable base URL per provider               | P2 | Custom endpoints |
| FR-6 | Model selection: predefined list + custom model option | P2 | Hybrid approach |
| FR-7 | Test connection feature for provider profiles    | P2 | Verification |
| FR-8 | Active profile selection                        | P1 | Profile management |
| FR-9 | Configurable fallback provider on failure       | P2 | Fault tolerance |
| FR-10 | AI commit message generation using new API      | P1 | Existing feature |
| FR-11 | AI branch name generation using new API         | P2 | Existing feature |
| FR-12 | AI PR content generation using new API          | P2 | Existing feature |
| FR-13 | Clean migration from old OpenRouter configuration | P1 | Breaking change |
| FR-14 | Settings UI updates for new provider system     | P1 | UI simplicity |

## Non-Goals (Out of Scope)

- Support for non-OpenAI-compatible APIs (Anthropic direct, Cohere direct, etc.)
- Provider-specific features beyond standard OpenAI API (e.g., DALL-E images)
- API usage tracking or quota management
- Cost estimation or spend limits
- Advanced authentication methods (OAuth, mTLS)

## Assumptions

- All target providers support the OpenAI Chat Completions API format
- Users have valid API keys for their chosen providers
- Local LLM servers expose an OpenAI-compatible endpoint
- The standard OpenAI API client libraries can be used for all providers

## Dependencies

- OpenAI Go SDK (`github.com/openai/openai-go`) or compatible HTTP client
- Existing config system (`config/config.go`)
- Existing TUI modal system (`tui/view.go`, `tui/update.go`)
- Existing notification system

## Acceptance Criteria

**FR-1: Replace OpenRouter client**
- [ ] OpenRouter-specific code removed from codebase
- [ ] New `openai` package created with generic API client
- [ ] All AI features use new client without errors

**FR-2: Support multiple providers**
- [ ] At least 4 providers work: OpenAI, Azure, local server, OpenRouter
- [ ] Each provider uses correct base URL format
- [ ] Provider type is stored in profile configuration

**FR-3: Provider profile CRUD**
- [ ] User can create new provider profiles (name, base URL, API key, model)
- [ ] User can view list of saved profiles
- [ ] User can edit existing profiles
- [ ] User can delete profiles (with confirmation if active)
- [ ] Profiles are persisted in config

**FR-4: API key authentication**
- [ ] API key is stored per profile
- [ ] Bearer token is sent in Authorization header
- [ ] API key is validated on save (optional test)

**FR-5: Customizable base URL**
- [ ] User can override default base URL for any provider
- [ ] Base URL supports custom ports and paths
- [ ] URL validation is performed

**FR-6: Model selection**
- [ ] Predefined list of common models (gpt-4, gpt-3.5-turbo, etc.)
- [ ] User can enter custom model name
- [ ] Selected model is stored per profile

**FR-7: Test connection**
- [ ] User can test API connection with current profile
- [ ] Test sends minimal request (e.g., list models or simple completion)
- [ ] Success/failure is displayed in notification
- [ ] Test shows actionable error messages

**FR-8: Active profile selection**
- [ ] User can select one profile as active
- [ ] Active profile is used for all AI operations
- [ ] Active profile indicator is shown in settings
- [ ] Active profile persists across restarts

**FR-9: Configurable fallback**
- [ ] User can designate a fallback profile
- [ ] If primary fails, fallback is attempted automatically
- [ ] User is notified when fallback is used
- [ ] Optional: disable fallback behavior

**FR-10, FR-11, FR-12: AI features**
- [ ] Commit message generation works with new API
- [ ] Branch name generation works with new API
- [ ] PR content generation works with new API
- [ ] All features respect active profile configuration

**FR-13: Clean migration**
- [ ] Old `openrouter_api_key`, `openrouter_model` fields removed from config
- [ ] Old AI settings modal removed
- [ ] Existing users prompted to create new profile on first use

**FR-14: Settings UI**
- [ ] Provider profiles accessible via Settings → AI Provider
- [ ] Profile management modal with intuitive navigation
- [ ] Help text explains required fields
- [ ] Current profile shown with status indicator

## Technical Considerations

### Configuration Structure

New config structure in `~/.config/jean/config.json`:

```json
{
  "repositories": {
    "<repo-path>": {
      "ai_provider": {
        "profiles": [
          {
            "name": "OpenAI",
            "type": "openai",
            "base_url": "https://api.openai.com/v1",
            "api_key": "sk-...",
            "model": "gpt-4"
          },
          {
            "name": "Local Ollama",
            "type": "custom",
            "base_url": "http://localhost:11434/v1",
            "api_key": "optional",
            "model": "llama2"
          }
        ],
        "active_profile": "OpenAI",
        "fallback_profile": "Local Ollama"
      },
      "ai_commit_enabled": true,
      "ai_branch_name_enabled": true,
      "ai_pr_content_enabled": true
    }
  }
}
```

### Provider Types

| Type      | Default Base URL                           | Example Custom URL               |
| --------- | ------------------------------------------ | -------------------------------- |
| openai    | https://api.openai.com/v1                  | (use default)                    |
| azure     | https://<resource>.openai.azure.com/openai/deployments/<deployment> | Custom endpoint |
| custom    | (user must provide)                        | http://localhost:11434/v1        |

### API Client Package

New package structure:
```
openai/
  ├── client.go       # Generic OpenAI-compatible client
  ├── provider.go     # Provider types and defaults
  ├── models.go       # Model list definitions
  └── errors.go       # Error handling and wrapping
```

### Code Changes Required

1. **Remove**: `openrouter/client.go` (replace with new `openai/` package)
2. **Update**: `config/config.go` - new provider config structure
3. **Update**: `tui/model.go` - new AI provider modal state
4. **Update**: `tui/update.go` - new modal handlers
5. **Update**: `tui/view.go` - new modal rendering
6. **Update**: `tui/model.go` - AI generation commands use new client

## Architecture Patterns

### SOLID Principles

**Single Responsibility Principle:**
- `client.go` - API communication only
- `provider.go` - Provider configuration and defaults
- `models.go` - Model list management
- `config.go` - Configuration persistence

**Open/Closed Principle:**
- New providers can be added by extending provider list
- No modification needed to client for new providers

**Dependency Inversion:**
- AI features depend on `AIClient` interface
- Concrete implementations are injected

### Creational Patterns

**Factory Pattern:**
- `NewProvider(type, config)` - Create provider client from type
- `NewClient(profile)` - Create API client from profile

**Builder Pattern:**
- `ProfileBuilder` - Construct profile config step-by-step for validation

### Structural Patterns

**Adapter Pattern:**
- Each provider type adapts to OpenAI API interface
- Azure-specific headers handled by adapter

### Design Recommendations

- **Simple Feature**: Use SOLID + minimal patterns (Factory for provider creation)
- **Avoid Over-Engineering**: No IoC/DI container needed for this scope
- **Interface Segregation**: Small focused interfaces (AIClient with only needed methods)

## Design Considerations

### UI Flow

**AI Provider Settings Modal:**
```
┌─────────────────────────────────────────┐
│ AI Provider Profiles                    │
├─────────────────────────────────────────┤
│ ▸ OpenAI (active)                       │
│   • gpt-4                               │
│                                         │
│ ▸ Local Ollama (fallback)               │
│   • llama2                              │
│                                         │
│ ▸ Azure OpenAI                          │
│   • gpt-35-turbo                        │
│                                         │
│ [New] [Edit] [Delete] [Test] [Close]    │
└─────────────────────────────────────────┘
```

**Profile Edit Modal:**
```
┌─────────────────────────────────────────┐
│ Edit Profile                            │
├─────────────────────────────────────────┤
│ Name:        [OpenAI           ]       │
│ Type:        [OpenAI            ▼]     │
│ Base URL:    [https://api...   ]       │
│ API Key:     [sk-...            ]       │
│ Model:       [gpt-4             ▼]     │
│                                         │
│ Fallback:    [☐] Set as fallback       │
│                                         │
│ [Test] [Save] [Cancel]                  │
└─────────────────────────────────────────┘
```

### Model List

Predefined models for dropdown:
- GPT-4: `gpt-4`, `gpt-4-turbo`, `gpt-4o`
- GPT-3.5: `gpt-3.5-turbo`
- Azure: `gpt-35-turbo`, `gpt-4`
- Custom: (user enters model name)

## Risks & Mitigations

| Risk                          | Impact | Mitigation                                |
| ----------------------------- | ------ | ----------------------------------------- |
| Azure API format differences  | Medium | Test with Azure, document requirements    |
| Local server reliability      | Low    | Fallback to cloud provider                |
| Breaking existing workflows   | High   | Clear migration message in UI             |
| API key exposure in logs      | High   | Redact keys in debug output               |
| Provider rate limits          | Medium | Retry with exponential backoff            |
| Config schema migration       | Medium | Add migration logic on startup            |

## Success Metrics

- All existing AI features work with new system
- User can switch between providers in <5 seconds
- Connection test provides actionable feedback within 3 seconds
- No regression in AI generation quality
- Code coverage for new openai package >80%

## Priority & Timeline

**Priority**: High (Breaking change, core feature dependency)

**Target**: Complete before next minor release

## Open Questions

1. Should we migrate existing OpenRouter config to a profile automatically, or force manual setup?
   - **Decision**: Clean break - user creates new profile

2. Should we store API keys encrypted in config?
   - **Recommendation**: Out of scope for now - rely on file system permissions

3. Should profile be per-repository or global?
   - **Recommendation**: Per-repository (matches current pattern)

## Glossary

- **OpenAI-compatible**: Any API that implements the OpenAI Chat Completions endpoint format
- **Profile**: Saved configuration for a provider (name, URL, key, model)
- **Active profile**: The currently selected provider for AI operations
- **Fallback profile**: Secondary provider to use if primary fails
- **Base URL**: The API endpoint URL for a provider

## Changelog

| Version | Date             | Summary of Changes                    |
| ------- | ---------------- | ------------------------------------- |
| 1       | 2026-01-04 13:30 | Initial PRD approved                  |
