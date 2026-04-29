Install

1. Build or install the iRecall web binary so it is available in PATH, for example:

   make build-web
   sudo cp bin/irecall-web /usr/local/bin/

2. Create a system user and data directory:

   sudo useradd --system --no-create-home --shell /usr/sbin/nologin irecall
   sudo mkdir -p /var/lib/irecall
   sudo chown irecall:irecall /var/lib/irecall

3. (Optional) Create an environment file at /etc/irecall/irecall-web.env to override the default data path:

   IRECALL_DATA_PATH=/var/lib/irecall

4. Install the systemd unit:

   sudo cp packaging/systemd/irecall-web.service /etc/systemd/system/
   sudo systemctl daemon-reload
   sudo systemctl enable --now irecall-web.service

Notes

- The web service runs with --api-only by default to avoid interactive password prompts on startup. Run the CLI interactively once to configure the web password if desired.
- The irecall-mcp binary is a stdio MCP server that should be launched by an MCP client, not by systemd. Configure IRECALL_BASE_URL and IRECALL_API_TOKEN in the client environment or client config.
