package irecallapi

import "testing"

func TestAPIErrorErrorMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  APIError
		want string
	}{
		{
			name: "status and message",
			err: APIError{
				StatusCode: 401,
				Message:    "bad auth",
			},
			want: "iRecall API returned 401: bad auth",
		},
		{
			name: "default message without status",
			err:  APIError{},
			want: "unexpected API error",
		},
		{
			name: "trimmed message without status",
			err: APIError{
				Message: "  upstream failed  ",
			},
			want: "upstream failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Fatalf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}
