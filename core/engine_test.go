package core

import (
	"reflect"
	"testing"
)

func TestParseJSONStringArray(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{
			name:  "plain json array",
			input: `["emmc","flash memory","partition"]`,
			want:  []string{"emmc", "flash memory", "partition"},
		},
		{
			name:  "markdown fenced json array",
			input: "```json\n[\"emmc\", \"flash memory\"]\n```",
			want:  []string{"emmc", "flash memory"},
		},
		{
			name:  "extra prose before array",
			input: "Here you go: [\"alpha\", \"beta\"]",
			want:  []string{"alpha", "beta"},
		},
		{
			name:  "comma fallback",
			input: `"Alpha", beta, gamma`,
			want:  []string{"alpha", "beta", "gamma"},
		},
		{
			name:    "empty response",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseJSONStringArray(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("parseJSONStringArray() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
