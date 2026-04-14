# terminal-interface

An SSH-accessible TUI portfolio. Connect with:

```
ssh tui.arkadiuszjuszczyk.com
```

Built in Go with [Wish](https://github.com/charmbracelet/wish), [Bubble Tea](https://github.com/charmbracelet/bubbletea), and [Lip Gloss](https://github.com/charmbracelet/lipgloss).

## Quick start

```bash
go run . --port 2222
```

Then in another terminal:

```bash
ssh -p 2222 localhost
```

A host key is generated on first run under `.ssh/id_ed25519`.

## Flags

```
--host        host to bind to (default: ::, accepts IPv6 and IPv4)
--port        port to listen on (default: 2222)
--host-key    path to SSH host key (default: .ssh/id_ed25519)
```

## Deployment

### One-time server setup

1. Provision a small Linux VPS.
2. Point a DNS record at it — A for IPv4, AAAA for IPv6, or both for dual-stack. The server listens on `::` by default, so it accepts both.
3. Create a non-root user on the server with sudo access and your SSH public key in `~/.ssh/authorized_keys`.
4. **Move admin SSH off port 22** so this service can own it. Edit `/etc/ssh/sshd_config`:
   ```
   Port 2222
   ```
   On Ubuntu 24.04, SSH uses systemd socket activation by default, which ignores the `Port` directive. Disable it:
   ```bash
   sudo systemctl disable --now ssh.socket
   sudo systemctl enable --now ssh.service
   sudo systemctl restart ssh.service
   ```
5. Give the deploy user passwordless sudo for the commands the script runs:
   ```bash
   sudo visudo -f /etc/sudoers.d/terminal-interface-deploy
   ```
   ```
   arek ALL=(ALL) NOPASSWD: ALL
   ```
   (Narrow this down if you prefer — see `deploy.sh` for the exact commands it invokes.)
6. Copy `deploy/terminal-interface.service` to `/etc/systemd/system/` and enable it:
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable terminal-interface
   ```
   The service won't start until the binary is deployed.

### Deploying

From your local machine:

```bash
DEPLOY_HOST=tui.example.com DEPLOY_USER=arek ./deploy.sh
```

The script cross-compiles a statically linked Linux binary, uploads it via `scp`, grants the `cap_net_bind_service` capability so it can bind to port 22 as a non-root user, restarts the systemd service, and verifies it's running.

Supported env vars:

| Variable | Default | Notes |
| --- | --- | --- |
| `DEPLOY_HOST` | *(required)* | Target hostname |
| `DEPLOY_USER` | *(required)* | SSH user on target |
| `DEPLOY_PORT` | `2222` | SSH port for admin connection |
| `DEPLOY_PATH` | `/home/$DEPLOY_USER/terminal-interface` | Remote binary path |
| `DEPLOY_ARCH` | `amd64` | Use `arm64` for ARM servers |
| `SERVICE_NAME` | `terminal-interface` | systemd unit name |

## Project layout

```
.
├── main.go              # entry point, Wish server, tea.Model
├── boot.go              # animated boot sequence
├── styles.go            # Lip Gloss styles and color palette
├── content.go           # bio, experience, stack, projects, contacts
├── deploy.sh            # cross-compile + deploy script
└── deploy/
    └── terminal-interface.service  # systemd unit template
```

## License

MIT
