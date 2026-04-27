package irecallapi

import "net/http"

func applyAuth(req *http.Request, token string) {
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "irecall-mcp/0")
}
