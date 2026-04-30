Install

This unit is intended for a headless local iRecall web API that MCP clients can reach over localhost. The irecall-mcp binary itself is a stdio MCP server and should be launched by an MCP client, not by systemd.

1. Build or install the iRecall web and MCP binaries at the paths used by the examples:

   make build-web build-mcp
   sudo cp bin/irecall-web /usr/local/bin/irecall-web
   sudo cp bin/irecall-mcp /usr/local/bin/irecall-mcp

2. Create a system user, data directory, and local configuration directory:

   sudo useradd --system --no-create-home --shell /usr/sbin/nologin irecall
   sudo mkdir -p /var/lib/irecall /etc/irecall
   sudo chown irecall:irecall /var/lib/irecall
   sudo chmod 0750 /var/lib/irecall

3. (Optional) Create an environment file at /etc/irecall/irecall-web.env to override defaults:

   IRECALL_DATA_PATH=/var/lib/irecall
   IRECALL_WEB_HOST=127.0.0.1
   IRECALL_PROVIDER_HOST=localhost
   IRECALL_PROVIDER_PORT=11434
   IRECALL_PROVIDER_HTTPS=false
   IRECALL_PROVIDER_API_KEY_PATH=/etc/irecall/provider-api-key
   IRECALL_PROVIDER_MODEL=

   For OpenAI, use:

   IRECALL_PROVIDER_HOST=api.openai.com
   IRECALL_PROVIDER_PORT=443
   IRECALL_PROVIDER_HTTPS=true
   IRECALL_PROVIDER_API_KEY_PATH=/etc/irecall/openai-api-key
   IRECALL_PROVIDER_MODEL=gpt-4.1-mini

   Store the provider API key in the referenced file and make it readable by the irecall service user:

   sudo install -m 0640 -o root -g irecall /path/to/openai-api-key /etc/irecall/openai-api-key

   The packaged unit binds to 127.0.0.1 by default. Only set IRECALL_WEB_HOST=0.0.0.0 if you intentionally expose the API on the network and have an external access-control plan such as a firewall, VPN, or authenticated reverse proxy.

   If you set a different data path, create it first and make sure it is writable by the irecall user:

   sudo mkdir -p /path/to/irecall
   sudo chown irecall:irecall /path/to/irecall

4. Install and start the systemd unit:

   sudo cp packaging/systemd/irecall-web.service /etc/systemd/system/
   sudo systemctl daemon-reload
   sudo systemctl enable --now irecall-web.service
   sudo systemctl status irecall-web.service --no-pager

   After changing /etc/irecall/irecall-web.env or the provider API key file, restart the service:

   sudo systemctl restart irecall-web.service
   sudo journalctl -u irecall-web.service -n 50 --no-pager

5. Issue an API token for MCP or other API clients and store it at the path used by the smoke tests:

   sudo /usr/local/bin/irecall-web auth issue-token \
     --data-path /var/lib/irecall \
     --write-token-file /etc/irecall/api-token
   sudo chown root:irecall /etc/irecall/api-token
   sudo chmod 0640 /etc/irecall/api-token

   The command writes the full token to /etc/irecall/api-token and prints only the token prefix. Prefer pointing `irecall-mcp --token-file` at this path instead of copying the plaintext token into client config.

End-to-end smoke test

1. Confirm the web API service is running:

   systemctl is-active irecall-web.service
   sudo journalctl -u irecall-web.service -n 20 --no-pager

2. Confirm the REST API accepts the issued token:

   TOKEN=$(sudo cat /etc/irecall/api-token)
   curl -fsS \
     -H "Authorization: Bearer $TOKEN" \
     http://127.0.0.1:9527/api/app/bootstrap-state

3. Confirm quote storage works through the REST API:

   TOKEN=$(sudo cat /etc/irecall/api-token)
   curl -fsS \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"content":"iRecall systemd smoke test"}' \
     http://127.0.0.1:9527/api/app/add-quote

   curl -fsS \
     -H "Authorization: Bearer $TOKEN" \
     'http://127.0.0.1:9527/api/app/list-quotes?limit=20&offset=0'

4. Configure an MCP client to launch the stdio server without a shell wrapper:

   Command:

   /usr/local/bin/irecall-mcp

   Args:

   --token-file /etc/irecall/api-token

   Environment:

   IRECALL_BASE_URL=http://127.0.0.1:9527

   Do not use `bash -lc`, `sh -c`, or inline `cat /etc/irecall/api-token` snippets to inject `IRECALL_API_TOKEN`. Shell wrappers can fail before the bridge starts under restricted MCP hosts, for example with `spawn /bin/sh EACCES`, and can expose tokens in process metadata or logs. Let `irecall-mcp` read the protected token file directly.

   Then call the irecall_health MCP tool. It should report ok=true.

Counting quotes

Use the dedicated count endpoint:

    TOKEN=$(sudo cat /etc/irecall/api-token)
    curl -fsS \
      -H "Authorization: Bearer $TOKEN" \
      http://127.0.0.1:9527/api/app/count-quotes

Notes

- The web service runs with --api-only by default to avoid interactive password prompts on startup. Run the CLI interactively once to configure the web password if desired.
- The irecall-mcp binary is a stdio MCP server that should be launched by an MCP client, not by systemd. Configure IRECALL_BASE_URL plus `--token-file` (or `IRECALL_API_TOKEN_FILE`) in the client launch config.
