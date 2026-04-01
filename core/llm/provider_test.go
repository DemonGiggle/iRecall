package llm

import "testing"

func TestProviderConfigBaseURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		cfg  ProviderConfig
		want string
	}{
		{
			name: "host only",
			cfg: ProviderConfig{
				Host: "localhost",
				Port: 11434,
			},
			want: "http://localhost:11434/v1",
		},
		{
			name: "host with path",
			cfg: ProviderConfig{
				Host: "foo.com/api",
				Port: 994,
			},
			want: "http://foo.com:994/api/v1",
		},
		{
			name: "host with trailing slash path",
			cfg: ProviderConfig{
				Host:  "foo.com/api/",
				Port:  994,
				HTTPS: true,
			},
			want: "https://foo.com:994/api/v1",
		},
		{
			name: "path already includes v1",
			cfg: ProviderConfig{
				Host: "foo.com/api/v1",
				Port: 994,
			},
			want: "http://foo.com:994/api/v1",
		},
		{
			name: "explicit scheme in host input",
			cfg: ProviderConfig{
				Host:  "https://foo.com/api",
				Port:  994,
				HTTPS: false,
			},
			want: "http://foo.com:994/api/v1",
		},
		{
			name: "existing port in host input gets replaced",
			cfg: ProviderConfig{
				Host: "foo.com:1234/api",
				Port: 994,
			},
			want: "http://foo.com:994/api/v1",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.cfg.BaseURL(); got != tt.want {
				t.Fatalf("BaseURL() = %q, want %q", got, tt.want)
			}
		})
	}
}
