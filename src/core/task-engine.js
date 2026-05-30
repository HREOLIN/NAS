import { createId, pushTask } from "./store.js";

const wait = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

export class TaskEngine {
  async enqueue({ module, action, payload, run }) {
    const task = {
      id: createId("task"),
      module,
      action,
      payload,
      status: "queued",
      progress: 0,
      startedAt: new Date().toISOString(),
      endedAt: null,
      error: null,
      rollback: false,
      steps: []
    };

    pushTask(task);

    const control = {
      async step(name, detail, delay = 250) {
        task.steps.push({
          name,
          detail,
          at: new Date().toISOString()
        });
        task.progress = Math.min(task.progress + 25, 95);
        await wait(delay);
      },
      fail(message, rollback = false) {
        task.error = message;
        task.rollback = rollback;
        throw new Error(message);
      }
    };

    task.status = "running";

    try {
      const result = await run(control);
      task.progress = 100;
      task.status = task.rollback ? "rolled_back" : "success";
      task.endedAt = new Date().toISOString();
      task.result = result;
      return task;
    } catch (error) {
      task.status = task.rollback ? "rolled_back" : "failed";
      task.endedAt = new Date().toISOString();
      task.error = task.error || error.message;
      return task;
    }
  }
}

