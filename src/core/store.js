import { randomUUID } from "node:crypto";

const now = () => new Date().toISOString();

const state = {
  startedAt: now(),
  tasks: [],
  recovery: {
    snapshots: [
      {
        id: "snap-system-001",
        volume: "volume1",
        name: "system-daily-001",
        createdAt: "2026-05-29T22:00:00.000Z",
        files: 1820
      },
      {
        id: "snap-project-002",
        volume: "volume2",
        name: "project-pre-upgrade",
        createdAt: "2026-05-30T01:15:00.000Z",
        files: 946
      }
    ],
    recoveries: []
  },
  network: {
    interfaces: [
      {
        id: "eth0",
        name: "LAN 1",
        mode: "static",
        ipv4: "192.168.10.20",
        mask: "255.255.255.0",
        gateway: "192.168.10.1",
        dns: ["223.5.5.5", "8.8.8.8"],
        link: "up",
        bridge: "br0"
      },
      {
        id: "eth1",
        name: "LAN 2",
        mode: "dhcp",
        ipv4: "192.168.10.21",
        mask: "255.255.255.0",
        gateway: "192.168.10.1",
        dns: ["223.5.5.5"],
        link: "standby",
        bridge: null
      }
    ],
    history: []
  },
  vm: {
    templates: [
      { id: "tpl-ubuntu-24", name: "Ubuntu 24.04 Base" },
      { id: "tpl-rocky-9", name: "Rocky Linux 9 Base" }
    ],
    vms: [
      {
        id: "vm-demo-001",
        name: "dev-jumpbox",
        cpu: 2,
        memoryGb: 4,
        diskGb: 60,
        networkId: "eth0",
        status: "stopped",
        createdAt: "2026-05-29T09:00:00.000Z",
        template: "Ubuntu 24.04 Base"
      }
    ]
  }
};

export function getState() {
  return state;
}

export function createId(prefix) {
  return `${prefix}-${randomUUID().slice(0, 8)}`;
}

export function pushTask(task) {
  state.tasks.unshift(task);
  state.tasks = state.tasks.slice(0, 30);
}

