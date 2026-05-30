# NAS 后端架构说明

## 1. 目标

这版 demo 的目标是先把 NAS 核心后端架构搭起来，而不是直接接入真实设备。

重点解决：

1. 三个核心模块怎么拆
2. 高风险操作如何任务化
3. Go 和 C 如何分层
4. 后续怎么替换为真实系统能力

## 2. 总体结构

```text
HTTP API
  -> App
    -> Task Engine
      -> Recovery Service
      -> Network Service
      -> VM Service
        -> Native Adapter (cgo)
          -> C Runtime Adapter
```

## 3. 模块职责

### Recovery Service

负责：

- 快照管理
- 文件恢复
- 卷回滚
- 恢复记录维护

当前通过 C 适配层模拟底层执行结果，后续可替换为真实文件系统能力。

### Network Service

负责：

- 网卡信息查询
- 配置校验
- 网络配置任务执行
- 连通性失败回滚

这部分设计重点是：

- 危险操作必须任务化
- 应用前必须校验
- 应用后必须可回滚

### VM Service

负责：

- VM 列表
- VM 创建
- VM 电源状态切换
- VM 元数据维护

后续适合替换到底层 `libvirt` 或 `qemu` 执行器。

## 4. Go 与 C 的边界

### Go 层负责

- HTTP API
- 参数校验
- 状态管理
- 任务引擎
- 领域模型

### C 层负责

- 模拟底层 NAS 能力执行
- 为后续系统调用保留统一接口入口

当前 C 代码是占位实现，但函数接口已经按真实场景拆开：

- 创建快照
- 恢复文件
- 卷回滚
- 应用网络配置
- 创建虚拟机
- 切换虚拟机电源状态

## 5. 任务引擎

所有关键操作都必须走统一任务引擎。

任务包含：

- `id`
- `module`
- `action`
- `status`
- `progress`
- `steps`
- `error`
- `rollback`
- `startedAt`
- `endedAt`

这样做的好处：

- 前后端都能统一查看任务状态
- 后续容易接入重试、审计、告警
- 可以把危险操作统一纳入回滚框架

## 6. 状态存储

当前使用内存态 Store，保存：

- `tasks`
- `recovery.snapshots`
- `recovery.recoveries`
- `network.interfaces`
- `network.history`
- `vm.vms`

真实项目建议拆为：

- 配置数据库
- 任务数据库
- 审计数据库
- 宿主机状态采集

## 7. 真实项目替换方向

### 数据恢复

建议替换到：

- `Btrfs`
- `ZFS`
- `LVM`

### 网络配置

建议替换到：

- `ip addr`
- `ip route`
- `bridge`
- `bond`
- `vlan`

### 虚拟机管理

建议替换到：

- `QEMU/KVM`
- `libvirt`

