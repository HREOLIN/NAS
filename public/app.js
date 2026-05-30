const state = {
  overview: null,
  recovery: null,
  network: null,
  vm: null,
  tasks: []
};

async function api(path, options = {}) {
  const response = await fetch(path, {
    headers: {
      "Content-Type": "application/json"
    },
    ...options
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || error.message || "Request failed");
  }

  return response.json();
}

function setHtml(id, html) {
  document.getElementById(id).innerHTML = html;
}

function fmt(value) {
  return new Date(value).toLocaleString("zh-CN", { hour12: false });
}

function badge(status) {
  const normalized = status === "queued" ? "running" : status;
  return `<span class="badge ${normalized}">${status}</span>`;
}

function renderOverview() {
  if (!state.overview) {
    return;
  }

  document.getElementById("hero-task-count").textContent = state.overview.modules.tasks;

  const items = [
    ["快照数", state.overview.modules.snapshots],
    ["恢复记录", state.overview.modules.recoveries],
    ["网络接口", state.overview.modules.interfaces],
    ["运行中 VM", state.overview.modules.runningVms],
    ["总任务", state.overview.modules.tasks]
  ];

  setHtml(
    "stats-grid",
    items
      .map(
        ([label, value]) => `
          <article class="stat-card">
            <p>${label}</p>
            <strong>${value}</strong>
          </article>
        `
      )
      .join("")
  );
}

function renderRecovery() {
  const snapshots = state.recovery?.snapshots || [];

  if (snapshots.length === 0) {
    setHtml("snapshot-list", `<div class="empty">还没有快照数据。</div>`);
    return;
  }

  setHtml(
    "snapshot-list",
    snapshots
      .map(
        (snapshot) => `
          <article class="list-item">
            <strong>${snapshot.name}</strong>
            <div class="meta-line">卷: ${snapshot.volume}</div>
            <div class="meta-line">文件数: ${snapshot.files}</div>
            <div class="meta-line">创建时间: ${fmt(snapshot.createdAt)}</div>
          </article>
        `
      )
      .join("")
  );
}

function renderNetwork() {
  const interfaces = state.network?.interfaces || [];

  setHtml(
    "network-list",
    interfaces
      .map(
        (item) => `
          <article class="list-item">
            <strong>${item.name} (${item.id})</strong>
            <div class="meta-line">模式: ${item.mode}</div>
            <div class="meta-line">IPv4: ${item.ipv4}</div>
            <div class="meta-line">网关: ${item.gateway}</div>
            <div class="meta-line">Bridge: ${item.bridge || "未绑定"}</div>
            <div class="meta-line">链路: ${item.link}</div>
          </article>
        `
      )
      .join("")
  );
}

function renderVms() {
  const vms = state.vm?.vms || [];

  setHtml(
    "vm-list",
    vms
      .map(
        (vm) => `
          <article class="list-item">
            <strong>${vm.name}</strong>
            <div class="meta-line">状态: ${vm.status}</div>
            <div class="meta-line">CPU / 内存 / 磁盘: ${vm.cpu} / ${vm.memoryGb}GB / ${vm.diskGb}GB</div>
            <div class="meta-line">网络接口: ${vm.networkId}</div>
            <div class="meta-line">模板: ${vm.template}</div>
          </article>
        `
      )
      .join("")
  );
}

function renderTasks() {
  if (!state.tasks.length) {
    setHtml("task-list", `<div class="empty">还没有任务，点击上面的操作按钮可以触发一条任务流。</div>`);
    return;
  }

  setHtml(
    "task-list",
    state.tasks
      .slice(0, 8)
      .map(
        (task) => `
          <article class="task-item">
            <div class="task-head">
              <strong>${task.module} / ${task.action}</strong>
              ${badge(task.status)}
            </div>
            <div class="meta-line">任务 ID: ${task.id}</div>
            <div class="meta-line">开始时间: ${fmt(task.startedAt)}</div>
            <div class="progress-bar">
              <div class="progress-fill" style="width:${task.progress || 0}%"></div>
            </div>
            <div class="meta-line">${task.error || `${task.steps.length} steps recorded`}</div>
          </article>
        `
      )
      .join("")
  );
}

async function refresh() {
  const [overview, recovery, network, vm, tasks] = await Promise.all([
    api("/api/overview"),
    api("/api/recovery"),
    api("/api/network"),
    api("/api/vms"),
    api("/api/tasks")
  ]);

  state.overview = overview;
  state.recovery = recovery;
  state.network = network;
  state.vm = vm;
  state.tasks = tasks;

  renderOverview();
  renderRecovery();
  renderNetwork();
  renderVms();
  renderTasks();
}

async function act(action) {
  try {
    await action();
    await refresh();
  } catch (error) {
    window.alert(error.message);
  }
}

document.getElementById("create-snapshot-btn").addEventListener("click", () =>
  act(() =>
    api("/api/recovery/snapshots", {
      method: "POST",
      body: JSON.stringify({
        volume: "volume1",
        name: `manual-snapshot-${Date.now()}`
      })
    })
  )
);

document.getElementById("restore-file-btn").addEventListener("click", () =>
  act(() =>
    api("/api/recovery/restore-file", {
      method: "POST",
      body: JSON.stringify({
        sourcePath: "/team/roadmap.xlsx",
        targetPath: "/restore/roadmap.xlsx"
      })
    })
  )
);

document.getElementById("rollback-volume-btn").addEventListener("click", () =>
  act(() =>
    api("/api/recovery/rollback-volume", {
      method: "POST",
      body: JSON.stringify({
        volume: "volume2"
      })
    })
  )
);

document.getElementById("apply-safe-network-btn").addEventListener("click", () =>
  act(() =>
    api("/api/network/apply", {
      method: "POST",
      body: JSON.stringify({
        nicId: "eth0",
        mode: "static",
        ipv4: "192.168.10.30",
        mask: "255.255.255.0",
        gateway: "192.168.10.1",
        dns: ["223.5.5.5", "1.1.1.1"]
      })
    })
  )
);

document.getElementById("apply-risky-network-btn").addEventListener("click", () =>
  act(() =>
    api("/api/network/apply", {
      method: "POST",
      body: JSON.stringify({
        nicId: "eth0",
        mode: "static",
        ipv4: "192.168.50.20",
        mask: "255.255.255.0",
        gateway: "192.168.50.1",
        dns: ["223.5.5.5"],
        simulateFailure: true
      })
    })
  )
);

document.getElementById("create-vm-btn").addEventListener("click", () =>
  act(() =>
    api("/api/vms", {
      method: "POST",
      body: JSON.stringify({
        name: `analytics-node-${Date.now().toString().slice(-4)}`,
        cpu: 4,
        memoryGb: 8,
        diskGb: 120,
        networkId: "eth0",
        template: "Rocky Linux 9 Base"
      })
    })
  )
);

document.getElementById("start-first-vm-btn").addEventListener("click", () =>
  act(async () => {
    const first = state.vm?.vms?.[0];

    if (!first) {
      throw new Error("No VM available");
    }

    return api(`/api/vms/${first.id}/power`, {
      method: "POST",
      body: JSON.stringify({ action: "start" })
    });
  })
);

document.getElementById("stop-first-vm-btn").addEventListener("click", () =>
  act(async () => {
    const first = state.vm?.vms?.[0];

    if (!first) {
      throw new Error("No VM available");
    }

    return api(`/api/vms/${first.id}/power`, {
      method: "POST",
      body: JSON.stringify({ action: "stop" })
    });
  })
);

refresh();
setInterval(refresh, 2500);
