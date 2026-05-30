#!/usr/bin/env bash

set -euo pipefail

APP_DIR="${APP_DIR:-/opt/nas-backend}"
SERVICE_NAME="${SERVICE_NAME:-nas-backend}"
REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "[1/5] creating application directory"
mkdir -p "${APP_DIR}"

echo "[2/5] syncing repository"
rsync -a --delete \
  --exclude ".git" \
  --exclude ".gocache" \
  --exclude ".gomodcache" \
  "${REPO_DIR}/" "${APP_DIR}/"

echo "[3/5] preparing go caches"
mkdir -p "${APP_DIR}/.gocache" "${APP_DIR}/.gomodcache"

echo "[4/5] installing systemd unit"
install -m 0644 "${APP_DIR}/deploy/linux/nas-backend.service" "/etc/systemd/system/${SERVICE_NAME}.service"
sed -i "s#/opt/nas-backend#${APP_DIR}#g" "/etc/systemd/system/${SERVICE_NAME}.service"

echo "[5/5] reloading and enabling service"
systemctl daemon-reload
systemctl enable --now "${SERVICE_NAME}.service"
systemctl status "${SERVICE_NAME}.service" --no-pager

echo
echo "service installed successfully"
echo "demo script: ${APP_DIR}/scripts/demo-linux-vm.sh"

