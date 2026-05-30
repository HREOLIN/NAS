# NAS 后端开发文档

## 1. 运行环境

当前 demo 依赖：

- `Go 1.22+`
- `gcc`
- `cgo`

这台机器已经具备 `go` 和 `gcc`，可以直接编译运行。

## 2. 启动方式

```bash
go run ./cmd/server
```

默认监听：

- `127.0.0.1:8080`

## 3. 接口说明

### 健康与总览

- `GET /api/health`
- `GET /api/overview`
- `GET /api/tasks`

### 数据恢复

- `GET /api/recovery`
- `POST /api/recovery/snapshots`
- `POST /api/recovery/restore-file`
- `POST /api/recovery/rollback-volume`

请求示例：

```json
{
  "volume": "volume1",
  "name": "nightly-snapshot"
}
```

### 网络配置

- `GET /api/network`
- `POST /api/network/apply`

请求示例：

```json
{
  "nicId": "eth0",
  "mode": "static",
  "ipv4": "192.168.10.30",
  "mask": "255.255.255.0",
  "gateway": "192.168.10.1",
  "dns": ["223.5.5.5"],
  "simulateFailure": false
}
```

### 虚拟机管理

- `GET /api/vms`
- `POST /api/vms`
- `POST /api/vms/{id}/power`

创建 VM 请求示例：

```json
{
  "name": "analytics-node",
  "cpu": 4,
  "memoryGb": 8,
  "diskGb": 120,
  "networkId": "eth0",
  "template": "Rocky Linux 9 Base"
}
```

电源操作请求示例：

```json
{
  "action": "start"
}
```

## 4. 代码结构说明

### `cmd/server`

启动入口，负责启动 HTTP 服务。

### `internal/app`

负责组装 Store、TaskEngine、Service 和 HTTP 路由。

### `internal/core`

负责：

- 任务引擎
- 状态存储
- 通用数据模型

### `internal/modules`

负责具体业务模块：

- `recovery`
- `network`
- `vm`

### `internal/native`

负责：

- C 头文件
- C 实现
- cgo 包装

## 5. 继续开发建议

建议下一阶段按这个顺序做：

1. 把 `network` 模块替换为真实 Linux 网络配置适配器
2. 给 `recovery` 模块增加备份恢复和恢复校验
3. 给 `vm` 模块接入 `libvirt`
4. 把当前同步任务引擎改为可持久化异步任务系统

## 6. 风险点

### 网络模块

最容易把设备配失联，所以必须长期坚持：

- 先校验
- 后应用
- 再探活
- 失败回滚

### 恢复模块

不能只关注“恢复成功”，还要关注：

- 恢复前影响分析
- 恢复后校验
- 审计记录

### VM 模块

不能只完成创建和启动，还要注意：

- 资源竞争
- 存储映射
- 网络绑定
- 电源状态一致性

