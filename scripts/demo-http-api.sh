#!/usr/bin/env bash

set -euo pipefail

BASE_URL="${BASE_URL:-http://127.0.0.1:8080}"

echo "== api overview =="
curl -s "${BASE_URL}/api/overview"
echo
echo

echo "== network summary =="
curl -s "${BASE_URL}/api/network"
echo
echo

echo "== storage summary =="
curl -s "${BASE_URL}/api/storage"
echo
echo

echo "== recovery summary =="
curl -s "${BASE_URL}/api/recovery"
echo
echo

echo "== create snapshot =="
curl -s -X POST "${BASE_URL}/api/recovery/snapshots" \
  -H "Content-Type: application/json" \
  -d '{
    "volume": "volume1",
    "name": "api-managed-snapshot"
  }'
echo
echo

echo "== apply network profile =="
curl -s -X POST "${BASE_URL}/api/network/apply" \
  -H "Content-Type: application/json" \
  -d '{
    "nicId": "eth0",
    "mode": "static",
    "ipv4": "192.168.10.40",
    "mask": "255.255.255.0",
    "gateway": "192.168.10.1",
    "dns": ["223.5.5.5", "1.1.1.1"],
    "simulateFailure": false
  }'
echo
echo

echo "== create vm =="
VM_RESPONSE="$(curl -s -X POST "${BASE_URL}/api/vms" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "api-managed-vm",
    "cpu": 2,
    "memoryGb": 4,
    "diskGb": 50,
    "networkId": "eth0",
    "template": "Ubuntu 24.04 Base"
  }')"
echo "${VM_RESPONSE}"
echo
echo

VM_ID="$(printf '%s' "${VM_RESPONSE}" | sed -n 's/.*"id":"\([^"]*\)".*/\1/p' | head -n 1)"

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

echo "== task center =="
curl -s "${BASE_URL}/api/tasks"
echo
