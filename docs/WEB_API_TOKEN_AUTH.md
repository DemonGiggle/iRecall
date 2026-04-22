# Web API Token Authentication

This note captures the agreed design for REST access to the iRecall web server.

## Goal

Provide a simple bearer-token mechanism for all REST API requests so clients can talk to the web server from localhost or a remote host with one shared credential.

## Agreed flow

1. The user opens the web app and signs in normally.
2. In `Settings`, the user creates an API token.
3. The server generates a random token and shows it once in a popup dialog.
4. The dialog warns that the token must be copied now because it will not be shown again.
5. The server stores only a hash of the token in the database.
6. Every REST request must send `Authorization: Bearer <token>`.
7. The server hashes the presented token and compares it with the stored hash.
8. A `Renew` action creates a new token, stores the new hash, and immediately invalidates the old token.

## UX notes

- Show the plaintext token only in the one-time popup.
- Keep a short fingerprint or prefix in the settings UI if the user needs to identify the active token later.
- Make the warning explicit: `Copy this token now. It will not be shown again.`

## Security notes

- A bearer token is better than unauthenticated localhost access.
- Treat the token as a secret and avoid logging it.
- A one-way hash is required in storage; a server-side pepper is a good option if SHA-256 is used.
- This token can coexist with the browser session login if the web UI still needs interactive authentication.

## Operational behavior

- The token should apply to every REST endpoint consistently.
- Renewing the token must make the previous token unusable immediately.
- Clients should not depend on query parameters or custom headers for authentication.
