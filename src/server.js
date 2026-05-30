import http from "node:http";
import path from "node:path";
import { argv } from "node:process";
import { fileURLToPath } from "node:url";
import { getState } from "./core/store.js";
import { TaskEngine } from "./core/task-engine.js";
import { RecoveryService } from "./services/recovery-service.js";
import { NetworkService } from "./services/network-service.js";
import { VMService } from "./services/vm-service.js";
import { createApiRouter } from "./routes/api.js";
import { json, serveStatic } from "./utils/http.js";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const publicDir = path.resolve(__dirname, "../public");

export function createServer() {
  const state = getState();
  const taskEngine = new TaskEngine();
  const recoveryService = new RecoveryService(taskEngine);
  const networkService = new NetworkService(taskEngine);
  const vmService = new VMService(taskEngine);
  const apiRouter = createApiRouter({
    recoveryService,
    networkService,
    vmService,
    state
  });

  return http.createServer(async (request, response) => {
    try {
      const url = new URL(request.url, "http://127.0.0.1");

      if (url.pathname.startsWith("/api/")) {
        return await apiRouter(request, response, url.pathname);
      }

      const served = await serveStatic(response, publicDir, url.pathname);

      if (!served) {
        return json(response, 404, { message: "Page not found" });
      }
    } catch (error) {
      return json(response, 500, {
        message: "Internal server error",
        error: error.message
      });
    }
  });
}

export function startServer(port = 4000, host = "127.0.0.1") {
  const server = createServer();
  server.listen(port, host, () => {
    console.log(`NAS Core Demo is running at http://${host}:${port}`);
  });
  return server;
}

if (argv[1] && path.resolve(argv[1]) === __filename) {
  startServer();
}
