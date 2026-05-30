# NAS Core Demo 开发文档

## 1. 技术选型

为了让 demo 在当前环境里即开即跑，这里刻意采用：

- `Node.js 原生 http`
- 原生 `ESM`
- 原生前端 `HTML + CSS + JavaScript`
- 零第三方依赖

这样做的原因是：

- 不依赖 `npm install`
- 方便快速阅读
- 方便后续迁移到任意正式框架

## 2. 启动方式

```bash
node src/server.js
```

服务默认监听：

- `127.0.0.1:4000`

## 3. API 说明

### 基础接口

- `GET /api/health`
- `GET /api/overview`
- `GET /api/tasks`

### 数据恢复

- `GET /api/recovery`
- `POST /api/recovery/snapshots`
- `POST /api/recovery/restore-file`
- `POST /api/recovery/rollback-volume`

### 网络配置

- `GET /api/network`
- `POST /api/network/apply`

### 虚拟机管理

- `GET /api/vms`
- `POST /api/vms`
- `POST /api/vms/:id/power`

## 4. 状态与任务流转

任务状态：

- `queued`
- `running`
- `success`
- `failed`
- `rolled_back`

每个任务都包含：

- `id`
- `module`
- `action`
- `status`
- `progress`
- `steps`
- `startedAt`
- `endedAt`

## 5. 模块扩展建议

### RecoveryService

后续可以增加：

- 备份恢复
- 异机恢复
- 恢复演练
- 结果校验器

### NetworkService

后续可以增加：

- VLAN
- Bond / LACP
- Bridge
- MTU
- 静态路由
- IPv6

### VMService

后续可以增加：

- VM 模板
- VM 快照
- VM 克隆
- 镜像导入导出
- PCIe / GPU 直通

## 6. 适合继续演进的方式

如果你后续准备做正式项目，建议这么迁移：

1. 先保留当前领域边界
2. 把内存态 Store 替换为数据库
3. 把服务里的模拟动作替换为真实系统适配器
4. 把任务引擎改成可持久化任务队列
5. 再补认证、权限、监控、告警

## 7. 推荐下一步

最适合继续往下做的是：

1. 给 `NetworkService` 增加真实 bridge / vlan 配置适配器
2. 给 `RecoveryService` 增加快照浏览和恢复点差异展示
3. 给 `VMService` 接入 `libvirt`

