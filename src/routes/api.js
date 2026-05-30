import { json, readJsonBody } from "../utils/http.js";

export function createApiRouter({ recoveryService, networkService, vmService, state }) {
  return async function route(request, response, pathname) {
    if (request.method === "GET" && pathname === "/api/health") {
      return json(response, 200, { ok: true, startedAt: state.startedAt });
    }

    if (request.method === "GET" && pathname === "/api/overview") {
      return json(response, 200, {
        startedAt: state.startedAt,
        modules: {
          snapshots: state.recovery.snapshots.length,
          recoveries: state.recovery.recoveries.length,
          interfaces: state.network.interfaces.length,
          runningVms: state.vm.vms.filter((item) => item.status === "running").length,
          totalVms: state.vm.vms.length,
          tasks: state.tasks.length
        }
      });
    }

    if (request.method === "GET" && pathname === "/api/tasks") {
      return json(response, 200, state.tasks);
    }

    if (request.method === "GET" && pathname === "/api/recovery") {
      return json(response, 200, recoveryService.getSummary());
    }

    if (request.method === "POST" && pathname === "/api/recovery/snapshots") {
      const body = await readJsonBody(request);
      const task = await recoveryService.createSnapshot(body);
      return json(response, 201, task);
    }

    if (request.method === "POST" && pathname === "/api/recovery/restore-file") {
      const body = await readJsonBody(request);
      const task = await recoveryService.restoreFile(body);
      return json(response, 202, task);
    }

    if (request.method === "POST" && pathname === "/api/recovery/rollback-volume") {
      const body = await readJsonBody(request);
      const task = await recoveryService.rollbackVolume(body);
      return json(response, 202, task);
    }

    if (request.method === "GET" && pathname === "/api/network") {
      return json(response, 200, networkService.getSummary());
    }

    if (request.method === "POST" && pathname === "/api/network/apply") {
      const body = await readJsonBody(request);
      const task = await networkService.applyProfile(body);
      const status = task.status === "failed" ? 400 : 202;
      return json(response, status, task);
    }

    if (request.method === "GET" && pathname === "/api/vms") {
      return json(response, 200, vmService.getSummary());
    }

    if (request.method === "POST" && pathname === "/api/vms") {
      const body = await readJsonBody(request);
      const task = await vmService.createVm(body);
      return json(response, 201, task);
    }

    if (request.method === "POST" && pathname.startsWith("/api/vms/") && pathname.endsWith("/power")) {
      const parts = pathname.split("/");
      const id = parts[3];
      const body = await readJsonBody(request);
      const task = await vmService.changePowerState(id, body.action);
      return json(response, 202, task);
    }

    return json(response, 404, { message: "Not found" });
  };
}

