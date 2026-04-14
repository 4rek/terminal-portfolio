#!/usr/bin/env bash
set -euo pipefail

# Deploy terminal-interface to a remote server over SSH.
#
# Required env vars:
#   DEPLOY_HOST  — hostname or IP (e.g. tui.arkadiuszjuszczyk.com)
#   DEPLOY_USER  — SSH user on the remote (e.g. arek)
#
# Optional env vars:
#   DEPLOY_PORT  — SSH port for deploy connection (default: 2222)
#   DEPLOY_PATH  — remote path for the binary (default: /home/$DEPLOY_USER/terminal-interface)
#   DEPLOY_ARCH  — target architecture (default: amd64, use arm64 for ARM servers)
#   SERVICE_NAME — systemd service name (default: terminal-interface)
#
# Usage:
#   DEPLOY_HOST=tui.arkadiuszjuszczyk.com DEPLOY_USER=arek ./deploy.sh

: "${DEPLOY_HOST:?DEPLOY_HOST is required}"
: "${DEPLOY_USER:?DEPLOY_USER is required}"

DEPLOY_PORT="${DEPLOY_PORT:-2222}"
DEPLOY_PATH="${DEPLOY_PATH:-/home/$DEPLOY_USER/terminal-interface}"
DEPLOY_ARCH="${DEPLOY_ARCH:-amd64}"
SERVICE_NAME="${SERVICE_NAME:-terminal-interface}"

BINARY_NAME="terminal-interface-linux-$DEPLOY_ARCH"

echo "==> Building for linux/$DEPLOY_ARCH"
GOOS=linux GOARCH="$DEPLOY_ARCH" CGO_ENABLED=0 go build -ldflags="-s -w" -o "$BINARY_NAME" .

echo "==> Uploading to $DEPLOY_USER@$DEPLOY_HOST:$DEPLOY_PATH"
scp -P "$DEPLOY_PORT" "$BINARY_NAME" "$DEPLOY_USER@$DEPLOY_HOST:$DEPLOY_PATH.new"

echo "==> Swapping binary and restarting service"
ssh -p "$DEPLOY_PORT" "$DEPLOY_USER@$DEPLOY_HOST" bash <<EOF
set -euo pipefail
chmod +x "$DEPLOY_PATH.new"
sudo setcap 'cap_net_bind_service=+ep' "$DEPLOY_PATH.new"
mv "$DEPLOY_PATH.new" "$DEPLOY_PATH"
sudo systemctl restart "$SERVICE_NAME"
sudo systemctl status "$SERVICE_NAME" --no-pager -l | head -20
EOF

echo "==> Cleaning up local build artifact"
rm "$BINARY_NAME"

echo "==> Done. Test with: ssh $DEPLOY_HOST"
