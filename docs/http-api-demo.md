# HTTP API 维护 NAS 核心功能 设计方案与 Demo

你现在想先做的是这一层：

> 通过 HTTP API 来维护 `数据恢复 / 网络配置 / 虚拟机管理`

这是非常对的切入方式，因为它能把 NAS 项目先拆成两个层次：

1. `控制层`
   负责接口、鉴权、任务编排、状态管理
2. `执行层`
   负责真正调用 Linux / 文件系统 / 虚拟化能力

先把控制层做稳，后面接真实底层能力会轻松很多。

## 1. 设计目标

这套 HTTP API 设计，要解决 4 件事：

1. 统一管理三大模块
2. 高风险操作全部任务化
3. 返回结构稳定，方便前端或脚本调用
4. 后续能平滑接真实 Linux 能力

## 2. 推荐架构

```text
Client / Script / Future UI
  -> HTTP API
    -> Task Engine
      -> Recovery Service
      -> Network Service
      -> VM Service
        -> Native Adapter
          -> Linux / FS / Hypervisor
```

这里最关键的是：

- 客户端不要直接碰底层命令
- 所有核心维护动作统一走 API
- 所有危险动作返回任务结果

## 3. 为什么 NAS 适合用 HTTP API 维护

因为 NAS 的核心动作基本都是：

- 有明确输入
- 有明确执行结果
- 需要审计
- 可能需要回滚

比如：

- 恢复一个文件
- 修改一块网卡 IP
- 启动一台 VM

这些都非常适合做成 HTTP API。

## 4. 三大模块的接口设计

### 数据恢复

建议分成两类接口：

1. `查询类`
   - 查看快照
   - 查看恢复记录
2. `执行类`
   - 创建快照
   - 恢复文件
   - 卷回滚

当前 demo 已有：

- `GET /api/recovery`
- `POST /api/recovery/snapshots`
- `POST /api/recovery/restore-file`
- `POST /api/recovery/rollback-volume`

请求示例：

```json
{
  "snapshotId": "snap-system-001",
  "sourcePath": "/srv/nas/team/roadmap.xlsx",
  "targetPath": "/srv/nas/restore/roadmap.xlsx"
}
```

### 网络配置

网络模块一定要任务化，因为它最危险。

当前 demo 已有：

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
  "dns": ["223.5.5.5", "1.1.1.1"],
  "simulateFailure": false
}
```

建议这个模块长期坚持：

- 先校验
- 生成执行计划
- 生成回滚计划
- 应用配置
- 探活
- 必要时回滚

### 虚拟机管理

虚拟机模块最适合先用 API 做生命周期管理。

当前 demo 已有：

- `GET /api/vms`
- `POST /api/vms`
- `POST /api/vms/{id}/power`

创建 VM 示例：

```json
{
  "name": "linux-vm-demo",
  "cpu": 2,
  "memoryGb": 4,
  "diskGb": 50,
  "networkId": "eth0",
  "template": "Ubuntu 24.04 Base"
}
```

电源操作示例：

```json
{
  "action": "start"
}
```

## 5. 统一返回模型建议

NAS 的核心维护动作不建议只返回一句“成功”。

建议统一返回任务对象：

```json
{
  "id": "task-0001",
  "module": "network",
  "action": "apply_profile",
  "status": "success",
  "progress": 100,
  "steps": [],
  "result": {}
}
```

这样做的好处：

- 脚本调用更稳定
- 将来前端也能复用
- 更方便做审计与告警

## 6. 建议的客户端调用方式

如果你现在自己维护 NAS，不建议一开始做大前端。

更推荐：

1. `curl`
2. shell 脚本
3. Python / Go 小客户端

所以我给你补了一套最直接的 HTTP API demo。

## 7. 仓库里的 API Demo

我新增了一套 API demo 文件：

- [scripts/demo-http-api.sh](/C:/Users/19296/Documents/NAS/scripts/demo-http-api.sh:1)
- [examples/http-api-demo.http](/C:/Users/19296/Documents/NAS/examples/http-api-demo.http:1)

### `demo-http-api.sh`

这是给 Linux 虚拟机或类 Unix 环境直接跑的脚本，适合你平时自己维护和测试。

运行方式：

```bash
chmod +x scripts/demo-http-api.sh
./scripts/demo-http-api.sh
```

或者指定服务地址：

```bash
BASE_URL=http://192.168.10.50:8080 ./scripts/demo-http-api.sh
```

### `http-api-demo.http`

这是给你后面用 IDE / HTTP Client 插件手动调接口的示例文件。

## 8. 你后面怎么继续扩

如果你未来真要把这套 NAS 做起来，我建议顺序是：

1. 先把 HTTP API 固定下来
2. 再把 network 接真实 Linux
3. 再把 recovery 接真实快照
4. 最后把 VM 接 `libvirt`

原因是：

- API 先稳，调用方不乱
- 底层能力可以逐步替换
- 团队协作也更顺

## 9. 最小可行 Demo 目标

如果只从“HTTP API 维护 NAS”这个目标看，你的 MVP 其实就是：

1. 能查状态
2. 能发动作
3. 能看任务结果
4. 能处理失败和回滚

你现在这套后端已经有这个雏形了。

## 10. 建议你下一步先做什么

如果你让我给一个最务实的建议：

- 先把 `HTTP API` 当成你的 NAS 控制面
- 先用脚本维护 NAS
- 不急着做界面
- 先把 `network` 接成真实 Linux 命令

因为一旦 network 真实跑起来，你就已经不是“概念 demo”，而是在真正控制一台 Linux NAS 节点了。

