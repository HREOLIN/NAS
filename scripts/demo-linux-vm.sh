#!/usr/bin/env bash

set -euo pipefail

BASE_URL="${BASE_URL:-http://127.0.0.1:8080}"

echo "== health =="
curl -s "${BASE_URL}/api/health"
echo
echo

echo "== overview =="
curl -s "${BASE_URL}/api/overview"
echo
echo

echo "== create snapshot =="
curl -s -X POST "${BASE_URL}/api/recovery/snapshots" \
  -H "Content-Type: application/json" \
  -d '{
    "volume": "volume1",
    "name": "linux-vm-nightly"
  }'
echo
echo

echo "== restore file =="
curl -s -X POST "${BASE_URL}/api/recovery/restore-file" \
  -H "Content-Type: application/json" \
  -d '{
    "snapshotId": "snap-system-001",
    "sourcePath": "/srv/nas/team/roadmap.xlsx",
    "targetPath": "/srv/nas/restore/roadmap.xlsx"
  }'
echo
echo

echo "== apply safe network profile =="
curl -s -X POST "${BASE_URL}/api/network/apply" \
  -H "Content-Type: application/json" \
  -d '{
    "nicId": "eth0",
    "mode": "static",
    "ipv4": "192.168.10.30",
    "mask": "255.255.255.0",
    "gateway": "192.168.10.1",
    "dns": ["223.5.5.5", "1.1.1.1"],
    "simulateFailure": false
  }'
echo
echo

echo "== apply risky network profile with rollback =="
curl -s -X POST "${BASE_URL}/api/network/apply" \
  -H "Content-Type: application/json" \
  -d '{
    "nicId": "eth0",
    "mode": "static",
    "ipv4": "192.168.20.30",
    "mask": "255.255.255.0",
    "gateway": "192.168.20.1",
    "dns": ["223.5.5.5"],
    "simulateFailure": true
  }'
echo
echo

echo "== create vm =="
CREATE_VM_RESPONSE="$(curl -s -X POST "${BASE_URL}/api/vms" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "linux-vm-demo",
    "cpu": 2,
    "memoryGb": 4,
    "diskGb": 50,
    "networkId": "eth0",
    "template": "Ubuntu 24.04 Base"
  }')"
echo "${CREATE_VM_RESPONSE}"
echo
echo

VM_ID="$(printf '%s' "${CREATE_VM_RESPONSE}" | sed -n 's/.*"id":"\([^"]*\)".*/\1/p' | head -n 1)"

if [[ -n "${VM_ID}" ]]; then
  echo "== start vm ${VM_ID} =="
  curl -s -X POST "${BASE_URL}/api/vms/${VM_ID}/power" \
    -H "Content-Type: application/json" \
    -d '{
      "action": "start"
    }'
  echo
  echo
fi

echo "== tasks =="
curl -s "${BASE_URL}/api/tasks"
echo

