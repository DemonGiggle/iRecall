# Web REST API Reference

This document describes the current HTTP endpoints exposed by the iRecall web server.

For the planned bearer-token design for non-browser clients, see [WEB_API_TOKEN_AUTH.md](./WEB_API_TOKEN_AUTH.md).

## Base URL

Use the web server address you started, for example:

- `http://127.0.0.1:9527`
- `http://0.0.0.0:9527` or a LAN host if you explicitly bind it

## Current authentication model

Today, the shipped web server uses browser-session authentication:

1. `POST /api/auth/login` with the web password
2. receive an `HttpOnly` session cookie
3. send that cookie on later authenticated requests

Most `/api/app/*` routes require an authenticated session cookie.

## Planned authentication model

The current code does **not** yet accept `Authorization: Bearer <token>` for REST calls. That design is captured separately in [WEB_API_TOKEN_AUTH.md](./WEB_API_TOKEN_AUTH.md).

## Conventions

- Request and response bodies use JSON unless noted otherwise.
- Error responses use:

```json
{
  "error": "message"
}
```

- `401 Unauthorized` means authentication is required or failed.
- `405 Method Not Allowed` means the route exists but the HTTP method is wrong.
- Field names are **case-sensitive** and currently mixed:
  - some responses use lower camel case, such as `productName`
  - many model-based responses use Go field names such as `ID`, `CreatedAt`, `Provider`

## Common object shapes

### Quote

```json
{
  "ID": 1,
  "GlobalID": "quote-uuid",
  "AuthorUserID": "user-1",
  "AuthorName": "Alice",
  "SourceUserID": "user-1",
  "SourceName": "Alice",
  "SourceBackend": "",
  "SourceNamespace": "",
  "SourceEntityType": "",
  "SourceEntityID": "",
  "SourceLabel": "",
  "SourceURL": "",
  "Content": "Example note",
  "Tags": ["sqlite", "notes"],
  "Version": 1,
  "IsOwnedByMe": true,
  "CreatedAt": "2026-04-22T15:00:00Z",
  "UpdatedAt": "2026-04-22T15:00:00Z"
}
```

### UserProfile

```json
{
  "UserID": "user-1",
  "DisplayName": "Alice",
  "CreatedAt": "2026-04-22T15:00:00Z",
  "UpdatedAt": "2026-04-22T15:00:00Z"
}
```

### Settings

```json
{
  "Provider": {
    "Host": "localhost",
    "Port": 11434,
    "HTTPS": false,
    "APIKey": "",
    "Model": "llama3"
  },
  "Search": {
    "MaxResults": 5,
    "MinRelevance": 0.3
  },
  "Debug": {
    "MockLLM": false
  },
  "Theme": "violet",
  "Web": {
    "Port": 9527
  },
  "RootDir": ""
}
```

## Endpoints

### `GET /api/auth/status`

Returns the current web-auth status.

Authentication: not required.

Parameters: none.

Response:

```json
{
  "runtime": "web",
  "passwordConfigured": true,
  "authenticated": false,
  "currentPort": 9527
}
```

### `POST /api/auth/login`

Validates the configured web password and starts a browser session.

Authentication: not required.

Request body:

```json
{
  "password": "your-password"
}
```

Response:

```json
{
  "ok": true
}
```

Notes:

- success sets the `irecall_session` cookie
- invalid password returns `401`

### `POST /api/auth/logout`

Ends the current browser session.

Authentication: not required, but usually called from an authenticated session.

Parameters: none.

Response:

```json
{
  "ok": true
}
```

### `POST /api/auth/change-password`

Changes the configured web password.

Authentication: required.

Request body:

```json
{
  "current": "old-password",
  "next": "new-password",
  "confirm": "new-password"
}
```

Response:

```json
{
  "ok": true
}
```

### `GET /api/app/bootstrap-state`

Returns the initial application state used to render the frontend shell.

Authentication: required.

Parameters: none.

Response:

```json
{
  "productName": "iRecall",
  "greeting": "Hi! Alice",
  "profile": {
    "UserID": "user-1",
    "DisplayName": "Alice",
    "CreatedAt": "2026-04-22T15:00:00Z",
    "UpdatedAt": "2026-04-22T15:00:00Z"
  },
  "settings": {
    "Provider": {
      "Host": "localhost",
      "Port": 11434,
      "HTTPS": false,
      "APIKey": "",
      "Model": ""
    },
    "Search": {
      "MaxResults": 5,
      "MinRelevance": 0
    },
    "Debug": {
      "MockLLM": false
    },
    "Theme": "violet",
    "Web": {
      "Port": 9527
    },
    "RootDir": ""
  },
  "paths": {
    "rootDir": "",
    "dataDir": "/home/user/.local/share/irecall",
    "configDir": "/home/user/.local/share/irecall",
    "stateDir": "/home/user/.local/share/irecall",
    "dbPath": "/home/user/.local/share/irecall/irecall.db",
    "logPath": "/home/user/.local/share/irecall/irecall.log"
  },
  "pages": ["Recall", "History", "Quotes", "Settings"],
  "docs": {
    "uiDesign": "docs/UI_DESIGN.md",
    "desktopMapping": "docs/WAILS_DESKTOP.md"
  }
}
```

### `GET /api/app/list-quotes`

Lists all stored quotes.

Authentication: required.

Parameters: none.

Response:

- `200 OK` with `Quote[]`

### `POST /api/app/add-quote`

Creates a new quote from free-form content.

Authentication: required.

Request body:

```json
{
  "content": "A new note or quote"
}
```

Response:

- `200 OK` with one `Quote`

### `POST /api/app/save-recall-as-quote`

Saves a recall question/response pair as a normal quote.

Authentication: required.

Request body:

```json
{
  "question": "What did I learn about SQLite?",
  "response": "You noted that WAL improves concurrency.",
  "keywords": ["sqlite", "wal"]
}
```

Response:

- `200 OK` with one `Quote`

### `POST /api/app/refine-quote-draft`

Runs quote-draft refinement and returns the refined text.

Authentication: required.

Request body:

```json
{
  "content": "rough draft text"
}
```

Response:

This endpoint returns a JSON string, not an object:

```json
"Refined quote text"
```

### `POST /api/app/update-quote`

Updates an existing quote.

Authentication: required.

Request body:

```json
{
  "id": 1,
  "content": "Updated content"
}
```

Response:

- `200 OK` with one `Quote`

### `POST /api/app/delete-quotes`

Deletes multiple quotes by ID.

Authentication: required.

Request body:

```json
{
  "ids": [1, 2, 3]
}
```

Response:

```json
{
  "ok": true
}
```

### `POST /api/app/preview-quote-export`

Builds a quote export payload and returns it as a JSON string.

Authentication: required.

Request body:

```json
{
  "ids": [1, 2]
}
```

Response:

This endpoint returns a JSON string containing the share-envelope payload:

```json
"{\"schema_version\":2,\"exported_at\":\"2026-04-22T15:00:00Z\",\"quotes\":[...]}"
```

### `POST /api/app/import-quotes-payload`

Imports quotes from a raw JSON share payload string.

Authentication: required.

Request body:

```json
{
  "payload": "{\"schema_version\":2,\"exported_at\":\"2026-04-22T15:00:00Z\",\"quotes\":[...]}"
}
```

Response:

```json
{
  "Inserted": 2,
  "Updated": 0,
  "Duplicates": 1,
  "Stale": 0
}
```

### `POST /api/app/save-user-profile`

Saves the local display name.

Authentication: required.

Request body:

```json
{
  "name": "Alice"
}
```

Response:

- `200 OK` with one `UserProfile`

### `POST /api/app/save-settings`

Saves application settings and returns the persisted settings.

Authentication: required.

Request body:

```json
{
  "Provider": {
    "Host": "localhost",
    "Port": 11434,
    "HTTPS": false,
    "APIKey": "",
    "Model": "llama3"
  },
  "Search": {
    "MaxResults": 5,
    "MinRelevance": 0.3
  },
  "Debug": {
    "MockLLM": false
  },
  "Theme": "violet",
  "Web": {
    "Port": 9527
  },
  "RootDir": ""
}
```

Response:

- `200 OK` with `Settings`

### `POST /api/app/fetch-models`

Fetches available models from the configured OpenAI-compatible provider.

Authentication: required.

Request body:

```json
{
  "Host": "localhost",
  "Port": 11434,
  "HTTPS": false,
  "APIKey": "",
  "Model": ""
}
```

Response:

```json
[
  "llama3",
  "qwen2.5"
]
```

### `POST /api/app/run-recall`

Runs the full recall flow: keyword extraction, quote search, answer generation, and history save.

Authentication: required.

Request body:

```json
{
  "question": "What did I learn about SQLite?"
}
```

Response:

```json
{
  "question": "What did I learn about SQLite?",
  "keywords": ["sqlite", "wal"],
  "quotes": [
    {
      "ID": 1,
      "GlobalID": "quote-uuid",
      "AuthorUserID": "user-1",
      "AuthorName": "Alice",
      "SourceUserID": "user-1",
      "SourceName": "Alice",
      "SourceBackend": "",
      "SourceNamespace": "",
      "SourceEntityType": "",
      "SourceEntityID": "",
      "SourceLabel": "",
      "SourceURL": "",
      "Content": "WAL improves concurrency.",
      "Tags": ["sqlite", "wal"],
      "Version": 1,
      "IsOwnedByMe": true,
      "CreatedAt": "2026-04-22T15:00:00Z",
      "UpdatedAt": "2026-04-22T15:00:00Z"
    }
  ],
  "response": "You noted that WAL improves concurrency."
}
```

### `GET /api/app/list-recall-history`

Lists saved recall-history summaries.

Authentication: required.

Parameters: none.

Response:

```json
[
  {
    "ID": 1,
    "Question": "What did I learn about SQLite?",
    "Response": "You noted that WAL improves concurrency.",
    "CreatedAt": "2026-04-22T15:00:00Z"
  }
]
```

### `GET /api/app/get-recall-history?id=<id>`

Returns one saved recall-history entry.

Authentication: required.

Query parameters:

- `id` (`int64`, required): history entry ID

Response:

```json
{
  "ID": 1,
  "Question": "What did I learn about SQLite?",
  "Response": "You noted that WAL improves concurrency.",
  "CreatedAt": "2026-04-22T15:00:00Z",
  "Quotes": [
    {
      "ID": 1,
      "GlobalID": "quote-uuid",
      "AuthorUserID": "user-1",
      "AuthorName": "Alice",
      "SourceUserID": "user-1",
      "SourceName": "Alice",
      "SourceBackend": "",
      "SourceNamespace": "",
      "SourceEntityType": "",
      "SourceEntityID": "",
      "SourceLabel": "",
      "SourceURL": "",
      "Content": "WAL improves concurrency.",
      "Tags": ["sqlite", "wal"],
      "Version": 1,
      "IsOwnedByMe": true,
      "CreatedAt": "2026-04-22T15:00:00Z",
      "UpdatedAt": "2026-04-22T15:00:00Z"
    }
  ]
}
```

### `POST /api/app/delete-recall-history`

Deletes multiple recall-history entries by ID.

Authentication: required.

Request body:

```json
{
  "ids": [1, 2]
}
```

Response:

```json
{
  "ok": true
}
```

## Non-REST helper route

### `GET /bridge.js`

Returns the JavaScript bridge used by the shipped frontend. This is not a public REST endpoint, but it is part of the current web runtime.
