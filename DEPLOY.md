# Deployment

How to host an instance of `terminal-interface` on your own domain.

## One-time server setup

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

## Deploying

From your local machine:

```bash
DEPLOY_HOST=tui.example.com DEPLOY_USER=arek ./deploy.sh
```

The script cross-compiles a statically linked Linux binary, uploads it via `scp`, grants the `cap_net_bind_service` capability so it can bind to port 22 as a non-root user, restarts the systemd service, and verifies it's running.

### Configuration

| Variable | Default | Notes |
| --- | --- | --- |
| `DEPLOY_HOST` | *(required)* | Target hostname |
| `DEPLOY_USER` | *(required)* | SSH user on target |
| `DEPLOY_PORT` | `2222` | SSH port for admin connection |
| `DEPLOY_PATH` | `/home/$DEPLOY_USER/terminal-interface` | Remote binary path |
| `DEPLOY_ARCH` | `amd64` | Use `arm64` for ARM servers |
| `SERVICE_NAME` | `terminal-interface` | systemd unit name |

## Troubleshooting

**Service fails to start with "permission denied" on port 22**
The `setcap` capability didn't apply. Re-run the deploy — the script sets it on each deploy.

**Colors don't render over SSH**
The server uses a per-session Lip Gloss renderer that respects each client's terminal. If colors are stripped, the client terminal likely reports `TERM=dumb`. Check with `echo $TERM` on the client.

**"Host key has changed" warning after redeploy**
Happens if the host key file was rotated. Users clear it with:
```bash
ssh-keygen -R tui.example.com
```
