# Website Checker

Periodically checks HTTP responses of monitored sites and sends Telegram notifications on status changes.

## How it works

1. **Startup** — the tool reads `domains.json`, checks every site once, and sends one Telegram message with the full status of all sites ("Sites added to monitoring").
2. **Regular checks** — every N seconds (set by `-interval`), each site is probed.
3. **Notifications**:
   - Site goes DOWN (status code outside `expected_codes`): notification sent immediately.
   - Site comes back UP: notification sent immediately.
   - Site stays DOWN: a reminder notification is sent every hour (configurable via `-renotify`).

State is kept in memory — no database or files are needed.

## Quick start
### 1. Configure

```bash
cp .env.example .env
# Edit .env: fill in TELEGRAM_BOT_TOKEN and TELEGRAM_CHAT_ID
```

Edit `domains.json` with the sites you want to monitor:

```json
{
    "example.com": {
        "expected_codes": [200],
        "retries": 3
    },
    "api.example.org": {
        "expected_codes": [200, 301],
        "retries": 2
    }
}
```

Each key is a domain name. `expected_codes` is the list of acceptable HTTP status codes. `retries` is how many times to retry on error before marking the site as DOWN.

### 2. Run locally

```bash
make run
```

Or directly:

```bash
go run ./cmd/checker -config domains.json -interval 30s -renotify 1h
```

## Managing domains with Go helper

Use the provided `domains-helper` binary to add or remove sites from `domains.json` and automatically restart the local checker service.

- Add a site:

```bash
./domains-helper add example.com 200,301
```

- Remove a site:

```bash
./domains-helper remove example.com
```

Notes:

- `domains.json` must already exist when running these commands.
- `add` fails if the site already exists.
- `remove` fails if the site is not present.
- The helper runs `./stop.sh` and then `./start.sh` after a successful change.

You can build the helper with:

```bash
go build -o domains-helper ./cmd/domains-helper
```

## Deployment

### Option A: systemd (recommended for servers)

1. Build for Linux:

   ```bash
   make build-linux
   ```

2. Copy files to the server:

   ```bash
   # Example:
   scp checker start.sh stop.sh .env domains.json user@server:/opt/checker/
   ```

3. Create a systemd unit file `/etc/systemd/system/checker.service`:

   ```ini
   [Unit]
   Description=Website Checker
   After=network.target

   [Service]
   Type=simple
   WorkingDirectory=/opt/checker
   ExecStart=/opt/checker/checker --config /opt/checker/domains.json --interval 30s --renotify 1h
   Restart=on-failure
   RestartSec=10
   User=root
   Group=root

   [Install]
   WantedBy=multi-user.target
   ```

4. Start the service:

   ```bash
   systemctl daemon-reload
   systemctl enable checker
   systemctl start checker
   systemctl status checker
   ```

### Option B: deploy script

```bash
./deploy.sh user@host /opt/checker
```

This builds the binary, copies all needed files via SCP, and prints instructions.

### Option C: manual start

```bash
./start.sh    # starts in background, writes checker.pid
./stop.sh     # stops by PID from checker.pid
```

## Command-line flags

| Flag            | Default   | Description                            |
|-----------------|-----------|----------------------------------------|
| `-config`       | `domains.json` | Path to domains configuration     |
| `-interval`     | `30s`     | Check interval (e.g. `10s`, `1m`)      |
| `-renotify`     | `1h`      | Re-notify interval for DOWN sites      |

## Environment variables

| Variable              | Required | Description                                   |
|-----------------------|----------|-----------------------------------------------|
| `TELEGRAM_BOT_TOKEN`  | yes      | Telegram bot token                            |
| `TELEGRAM_CHAT_ID`    | yes      | Telegram chat ID for notifications            |
| `PROXY`               | no       | SOCKS5 proxy address (e.g. `127.0.0.1:9050`) |

Variables are loaded from a `.env` file in the working directory.
