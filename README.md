# terminal-interface

An SSH-accessible TUI portfolio. Connect with:

```
ssh tui.arkadiuszjuszczyk.com
```

Built in Go with [Wish](https://github.com/charmbracelet/wish), [Bubble Tea](https://github.com/charmbracelet/bubbletea), and [Lip Gloss](https://github.com/charmbracelet/lipgloss).

## Local development

```bash
go run . --port 2222
```

Then in another terminal:

```bash
ssh -p 2222 localhost
```

A host key is generated on first run under `.ssh/id_ed25519`.

## Deployment

### One-time server setup

1. Provision a small Linux VPS (Hetzner, DigitalOcean, Fly.io — all fine).
2. Point a DNS record at it:
   - **IPv4 VPS**: A record → public IPv4 address
   - **IPv6-only VPS** (e.g. Hetzner CAX IPv6-only): AAAA record → public IPv6 address
   - **Dual-stack**: both A and AAAA records

   The server binds to `::` by default, which accepts IPv6 connections (and IPv4 on dual-stack systems).
3. SSH in as root, create a user:
   ```bash
   adduser arek
   usermod -aG sudo arek
   ```
4. Move admin SSH to a non-standard port so the TUI can own port 22. Edit `/etc/ssh/sshd_config`:
   ```
   Port 2222
   ```
   Then `sudo systemctl restart ssh`.
5. Copy your local public key into `/home/arek/.ssh/authorized_keys`.
6. Copy the systemd unit from `deploy/terminal-interface.service` into `/etc/systemd/system/` on the server, then:
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable terminal-interface
   ```

### Deploying a new version

From your local machine:

```bash
DEPLOY_HOST=tui.arkadiuszjuszczyk.com \
DEPLOY_USER=arek \
./deploy.sh
```

The script builds a statically linked Linux binary, uploads it over SCP, grants it the capability to bind to port 22, and restarts the systemd service.

For ARM servers (Hetzner CAX):

```bash
DEPLOY_ARCH=arm64 DEPLOY_HOST=... DEPLOY_USER=... ./deploy.sh
```

## Flags

```
--host        host to bind to (default: 0.0.0.0)
--port        port to listen on (default: 2222)
--host-key    path to SSH host key (default: .ssh/id_ed25519)
```

## Project structure

```
.
├── main.go              # entry point, Wish server, model/update/view
├── boot.go              # animated boot sequence
├── styles.go            # Lip Gloss styles and color palette
├── content.go           # bio, experience, stack, projects, contacts
├── deploy.sh            # deployment script
└── deploy/
    └── terminal-interface.service  # systemd unit
```
