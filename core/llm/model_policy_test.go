package llm

import "testing"

func TestPolicyForModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		model                 string
		wantTokenLimitParam   tokenLimitParameter
		wantTemperaturePolicy bool
		wantMinOutputTokens   int
	}{
		{
			model:                 "llama3.1",
			wantTokenLimitParam:   tokenLimitMaxTokens,
			wantTemperaturePolicy: true,
		},
		{
			model:                 "gpt-4.1-mini",
			wantTokenLimitParam:   tokenLimitMaxTokens,
			wantTemperaturePolicy: true,
		},
		{
			model:                 "gpt-5-nano",
			wantTokenLimitParam:   tokenLimitMaxCompletionTokens,
			wantTemperaturePolicy: false,
			wantMinOutputTokens:   1024,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.model, func(t *testing.T) {
			t.Parallel()
			got := policyForModel(tt.model)
			if got.TokenLimitParameter != tt.wantTokenLimitParam {
				t.Fatalf("TokenLimitParameter = %q, want %q", got.TokenLimitParameter, tt.wantTokenLimitParam)
			}
			if got.SupportsTemperature != tt.wantTemperaturePolicy {
				t.Fatalf("SupportsTemperature = %v, want %v", got.SupportsTemperature, tt.wantTemperaturePolicy)
			}
		})
	}
}
