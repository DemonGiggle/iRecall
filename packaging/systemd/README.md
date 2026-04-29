Install

This unit is intended for a headless local iRecall web API that MCP clients can reach over localhost. The irecall-mcp binary itself is a stdio MCP server and should be launched by an MCP client, not by systemd.

1. Build or install the iRecall web and MCP binaries at the paths used by the examples:

   make build-web build-mcp
   sudo cp bin/irecall-web /usr/local/bin/irecall-web
   sudo cp bin/irecall-mcp /usr/local/bin/irecall-mcp

2. Create a system user and data directory:

   sudo useradd --system --no-create-home --shell /usr/sbin/nologin irecall
   sudo mkdir -p /var/lib/irecall
   sudo chown irecall:irecall /var/lib/irecall

3. (Optional) Create an environment file at /etc/irecall/irecall-web.env to override defaults:

   IRECALL_DATA_PATH=/var/lib/irecall
   IRECALL_WEB_HOST=127.0.0.1

   The packaged unit binds to 127.0.0.1 by default. Only set IRECALL_WEB_HOST=0.0.0.0 if you intentionally expose the API on the network and have an external access-control plan such as a firewall, VPN, or authenticated reverse proxy.

   If you set a different data path, create it first and make sure it is writable by the irecall user:

   sudo mkdir -p /path/to/irecall
   sudo chown irecall:irecall /path/to/irecall

4. Install and start the systemd unit:

   sudo cp packaging/systemd/irecall-web.service /etc/systemd/system/
   sudo systemctl daemon-reload
   sudo systemctl enable --now irecall-web.service
   sudo systemctl status irecall-web.service --no-pager

5. Issue an API token for MCP or other API clients:

   sudo mkdir -p /etc/irecall
   sudo /usr/local/bin/irecall-web auth issue-token \
     --data-path /var/lib/irecall \
     --write-token-file /etc/irecall/api-token
   sudo chown root:irecall /etc/irecall/api-token
   sudo chmod 0640 /etc/irecall/api-token

   Use the generated token as IRECALL_API_TOKEN in the MCP client or other API client environment.

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

4. Configure an MCP client to launch the stdio server:

   Command:

   /usr/local/bin/irecall-mcp

   Environment:

   IRECALL_BASE_URL=http://127.0.0.1:9527
   IRECALL_API_TOKEN=<contents of /etc/irecall/api-token>

   Then call the irecall_health MCP tool. It should report ok=true.

Counting quotes

The REST API does not currently expose a dedicated count endpoint. To count all stored quotes, list quotes with limit=0 and count the returned JSON array:

   TOKEN=$(sudo cat /etc/irecall/api-token)
   curl -fsS \
     -H "Authorization: Bearer $TOKEN" \
     'http://127.0.0.1:9527/api/app/list-quotes?limit=0' \
     | python3 -c 'import json,sys; print(len(json.load(sys.stdin)))'

Notes

- The web service runs with --api-only by default to avoid interactive password prompts on startup. Run the CLI interactively once to configure the web password if desired.
- The irecall-mcp binary is a stdio MCP server that should be launched by an MCP client, not by systemd. Configure IRECALL_BASE_URL and IRECALL_API_TOKEN in the client environment or client config.
