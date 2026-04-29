package llm

import "strings"

type tokenLimitParameter string

const (
	tokenLimitMaxTokens           tokenLimitParameter = "max_tokens"
	tokenLimitMaxCompletionTokens tokenLimitParameter = "max_completion_tokens"
)

type modelPolicy struct {
	TokenLimitParameter     tokenLimitParameter
	SupportsTemperature     bool
	SupportsReasoningEffort bool
	MinOutputTokens         int
	DefaultOutputTokens     int
}

func policyForModel(model string) modelPolicy {
	family := modelFamily(model)
	policy := modelPolicy{
		TokenLimitParameter:     tokenLimitMaxTokens,
		SupportsTemperature:     true,
		SupportsReasoningEffort: true,
	}

	switch family {
	case "gpt-5":
		policy.TokenLimitParameter = tokenLimitMaxCompletionTokens
		policy.SupportsTemperature = false
		policy.SupportsReasoningEffort = false
		policy.MinOutputTokens = 4096
		policy.DefaultOutputTokens = 4096
	}

	return policy
}

func modelFamily(model string) string {
	name := strings.ToLower(strings.TrimSpace(model))
	switch {
	case strings.HasPrefix(name, "gpt-5"):
		return "gpt-5"
	default:
		return ""
	}
}

func (p modelPolicy) outputTokenLimit(requested int) int {
	if requested < p.MinOutputTokens {
		return p.MinOutputTokens
	}
	return requested
}
