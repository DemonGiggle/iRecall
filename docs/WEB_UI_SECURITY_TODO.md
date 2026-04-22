# Web UI Security Follow-Up

This note tracks the remaining security work for the networked web UI after the password-enrollment hardening.

## Already done

- First-time web password setup now happens in the terminal before the HTTP listener starts.
- The remote `/api/auth/setup` path is removed from the web API.
- Web password setup and password changes now enforce a minimum password policy.

## Remaining work

### 1. Add REST API bearer-token auth

- Generate a token from the Settings page and show it once in a popup.
- Store only a hash in the database and compare against `Authorization: Bearer <token>`.
- Add a `Renew` action that replaces the stored hash and invalidates the old token immediately.

### 2. Restrict the default bind address

- Change the default host from `0.0.0.0` to `127.0.0.1`.
- Keep explicit LAN exposure as an opt-in choice via `--host`.

### 3. Harden transport for LAN use

- Document the supported deployment model for non-localhost use.
- Prefer TLS directly in the app or require a trusted reverse proxy.
- Mark session cookies as `Secure` when served over HTTPS.

### 4. Add browser-origin protections

- Reject unexpected `Host` headers.
- Validate `Origin` on state-changing routes.
- Review DNS rebinding protections for browser access over IP/hostnames.

### 5. Add rate limiting and abuse controls

- Throttle login attempts per IP and per session.
- Add backoff or temporary lockout after repeated failed logins.
- Add limits for expensive authenticated routes such as recall/model fetch.

### 6. Replace the default HTTP server

- Stop using bare `http.ListenAndServe`.
- Create an explicit `http.Server` with:
  - `ReadHeaderTimeout`
  - `ReadTimeout`
  - `WriteTimeout`
  - `IdleTimeout`
  - `MaxHeaderBytes`

### 7. Review outbound request controls

- `fetch-models` currently accepts caller-supplied provider host/port.
- Revisit whether arbitrary destinations should be allowed in web mode.
- If retained, consider allowlists or stronger validation to reduce SSRF risk.

### 8. Review security headers

- Add or verify:
  - `Content-Security-Policy`
  - `X-Content-Type-Options`
  - `X-Frame-Options` or equivalent CSP frame control
  - `Referrer-Policy`

## Suggested implementation order

1. Add REST API bearer-token auth.
2. Default bind to loopback.
3. Add explicit `http.Server` timeouts.
4. Add login rate limiting.
5. Add `Host` and `Origin` enforcement.
6. Rework TLS / proxy expectations.
7. Revisit outbound destination controls.
