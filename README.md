# NAS Core Demo

一个零依赖、可直接运行的 NAS 核心功能框架 demo，覆盖 3 个核心模块：

- `数据恢复`
- `网络配置`
- `虚拟机管理`

这个 demo 的目标不是直接对接真实硬件，而是先把 NAS 核心模块的服务边界、任务编排、状态管理和前后端交互方式搭起来，方便后续继续扩展到真实的 `Btrfs / ZFS / libvirt / Linux bridge` 能力。

## 快速启动

本项目不依赖第三方包，直接用 Node.js 即可启动。

```bash
node src/server.js
```

启动后访问：

- `http://127.0.0.1:4000`

## Demo 包含的能力

### 数据恢复

- 快照列表
- 创建快照
- 文件恢复任务
- 卷回滚任务
- 恢复任务状态跟踪

### 网络配置

- 网卡状态展示
- 应用基础网络配置
- 配置前校验
- 配置后探活模拟
- 失败自动回滚演示

### 虚拟机管理

- VM 列表
- 创建虚拟机
- 启动 / 停止虚拟机
- 资源概览
- 任务化执行

## 目录结构

```text
public/               前端演示页面
src/
  core/               状态存储与任务引擎
  routes/             HTTP API 路由
  services/           领域服务
  utils/              通用 HTTP 工具
docs/                 架构与开发文档
```

## 文档

- [架构说明](./docs/architecture.md)
- [开发文档](./docs/development.md)

## 后续扩展方向

- 把 `NetworkService` 的模拟配置替换为真实 `iproute2 / bridge / bond / vlan` 调用
- 把 `RecoveryService` 的快照操作接到真实文件系统能力
- 把 `VMService` 的生命周期操作接到 `libvirt / qemu-kvm`
- 把内存状态迁移到数据库和宿主机 Agent

