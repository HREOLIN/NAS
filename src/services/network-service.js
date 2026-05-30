import { getState } from "../core/store.js";

const isIpv4 = (value) => /^(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}$/.test(value);

export class NetworkService {
  constructor(taskEngine) {
    this.taskEngine = taskEngine;
    this.state = getState();
  }

  getSummary() {
    return this.state.network;
  }

  validateProfile(profile) {
    const errors = [];

    if (!profile.nicId) {
      errors.push("nicId is required");
    }

    if (!["dhcp", "static"].includes(profile.mode)) {
      errors.push("mode must be dhcp or static");
    }

    if (profile.mode === "static") {
      for (const field of ["ipv4", "mask", "gateway"]) {
        if (!profile[field] || !isIpv4(profile[field])) {
          errors.push(`${field} must be a valid IPv4 address`);
        }
      }
    }

    return errors;
  }

  async applyProfile(profile) {
    const errors = this.validateProfile(profile);

    if (errors.length > 0) {
      return {
        id: null,
        module: "network",
        action: "apply_profile",
        status: "failed",
        progress: 0,
        steps: [],
        error: errors.join("; ")
      };
    }

    return this.taskEngine.enqueue({
      module: "network",
      action: "apply_profile",
      payload: profile,
      run: async (task) => {
        const nic = this.state.network.interfaces.find((item) => item.id === profile.nicId);

        if (!nic) {
          task.fail("Network interface not found");
        }

        const before = { ...nic, dns: [...nic.dns] };
        this.state.network.history.unshift({
          at: new Date().toISOString(),
          nicId: nic.id,
          before
        });

        await task.step("backup-config", "Backup current network profile");
        await task.step("apply-config", "Apply new interface config");

        nic.mode = profile.mode;
        nic.ipv4 = profile.mode === "dhcp" ? "192.168.10.88" : profile.ipv4;
        nic.mask = profile.mode === "dhcp" ? "255.255.255.0" : profile.mask;
        nic.gateway = profile.mode === "dhcp" ? "192.168.10.1" : profile.gateway;
        nic.dns = profile.mode === "dhcp" ? ["223.5.5.5"] : profile.dns || ["223.5.5.5"];

        await task.step("probe-connectivity", "Probe management connectivity", 350);

        if (profile.simulateFailure) {
          nic.mode = before.mode;
          nic.ipv4 = before.ipv4;
          nic.mask = before.mask;
          nic.gateway = before.gateway;
          nic.dns = before.dns;
          await task.step("rollback", "Rollback to previous network profile");
          task.fail("Connectivity probe failed, auto rollback applied", true);
        }

        await task.step("confirm", "Confirm management path healthy");
        return nic;
      }
    });
  }
}

