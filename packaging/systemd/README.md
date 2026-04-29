Install

1. Build or install the iRecall binaries so they are available in PATH, for example:

   make build-web build-mcp
   sudo cp bin/irecall-web /usr/local/bin/
   sudo cp bin/irecall-mcp /usr/local/bin/

2. Create a system user and data directory:

   sudo useradd --system --no-create-home --shell /usr/sbin/nologin irecall
   sudo mkdir -p /var/lib/irecall
   sudo chown irecall:irecall /var/lib/irecall

3. (Optional) Create an environment file at /etc/irecall/irecall.env or /etc/irecall/irecall-web.env to set variables:

   IRECALL_DATA_PATH=/var/lib/irecall
   IRECALL_BASE_URL=http://127.0.0.1:9527
   IRECALL_API_TOKEN=your_api_token_here

4. Install the systemd units:

   sudo cp packaging/systemd/irecall-web.service /etc/systemd/system/
   sudo cp packaging/systemd/irecall-mcp.service /etc/systemd/system/
   sudo systemctl daemon-reload
   sudo systemctl enable --now irecall-web.service irecall-mcp.service

Notes

- The web service runs with --api-only by default to avoid interactive password prompts on startup. Run the CLI interactively once to configure the web password if desired.
- The MCP service expects the web API to be available; set IRECALL_BASE_URL and IRECALL_API_TOKEN in the environment file.
