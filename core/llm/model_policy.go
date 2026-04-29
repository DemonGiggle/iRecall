package llm

import "strings"

type tokenLimitParameter string

const (
	tokenLimitMaxTokens           tokenLimitParameter = "max_tokens"
	tokenLimitMaxCompletionTokens tokenLimitParameter = "max_completion_tokens"
)

type modelPolicy struct {
	TokenLimitParameter tokenLimitParameter
	SupportsTemperature bool
}

func policyForModel(model string) modelPolicy {
	family := modelFamily(model)
	policy := modelPolicy{
		TokenLimitParameter: tokenLimitMaxTokens,
		SupportsTemperature: true,
	}

	switch family {
	case "gpt-5":
		policy.TokenLimitParameter = tokenLimitMaxCompletionTokens
		policy.SupportsTemperature = false
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
