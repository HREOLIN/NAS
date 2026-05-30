import { createId, getState } from "../core/store.js";

export class VMService {
  constructor(taskEngine) {
    this.taskEngine = taskEngine;
    this.state = getState();
  }

  getSummary() {
    return this.state.vm;
  }

  async createVm(input = {}) {
    return this.taskEngine.enqueue({
      module: "vm",
      action: "create_vm",
      payload: input,
      run: async (task) => {
        await task.step("reserve-resources", "Reserve CPU and memory");
        await task.step("create-disk", "Create virtual disk");
        await task.step("attach-network", "Attach VM to selected network");
        await task.step("register-vm", "Register VM definition");

        const vm = {
          id: createId("vm"),
          name: input.name || `demo-vm-${Date.now()}`,
          cpu: Number(input.cpu || 2),
          memoryGb: Number(input.memoryGb || 4),
          diskGb: Number(input.diskGb || 80),
          networkId: input.networkId || "eth0",
          status: "stopped",
          createdAt: new Date().toISOString(),
          template: input.template || "Ubuntu 24.04 Base"
        };

        this.state.vm.vms.unshift(vm);
        return vm;
      }
    });
  }

  async changePowerState(id, action) {
    return this.taskEngine.enqueue({
      module: "vm",
      action: `power_${action}`,
      payload: { id, action },
      run: async (task) => {
        const vm = this.state.vm.vms.find((item) => item.id === id);

        if (!vm) {
          task.fail("VM not found");
        }

        await task.step("check-vm", "Check VM state");
        await task.step("invoke-hypervisor", `Invoke ${action} action`);

        if (action === "start") {
          vm.status = "running";
        } else if (action === "stop") {
          vm.status = "stopped";
        } else if (action === "restart") {
          vm.status = "running";
        } else {
          task.fail("Unsupported power action");
        }

        await task.step("refresh-status", "Refresh VM runtime state");
        return vm;
      }
    });
  }
}

