# NAS Backend Demo

一个面向 NAS 核心能力的后端框架 demo，使用 `Go + C` 实现。

这版不做 UI，重点把后端骨架先立起来，覆盖 3 个核心模块：

- `数据恢复`
- `网络配置`
- `虚拟机管理`

## 技术路线

- `Go` 负责 HTTP API、任务编排、状态管理、服务边界
- `C` 负责底层能力适配层示例，后续可以替换为真实系统调用
- `cgo` 负责 Go 与 C 的桥接

## 快速启动

Windows 环境建议：

```powershell
$env:GOCACHE="C:\Users\19296\Documents\NAS\.gocache"
$env:GOMODCACHE="C:\Users\19296\Documents\NAS\.gomodcache"
go run ./cmd/server
```

Linux 环境建议：

```bash
export GOCACHE=/opt/nas-backend/.gocache
export GOMODCACHE=/opt/nas-backend/.gomodcache
go run ./cmd/server
```

默认监听：

- `http://127.0.0.1:8080`

## 当前能力

### 数据恢复

- 获取快照列表
- 创建快照
- 模拟文件恢复
- 模拟卷回滚
- 恢复任务记录

### 网络配置

- 获取网卡状态
- 应用网络配置
- 配置前校验
- 执行计划生成
- 探活与自动回滚演示

### 虚拟机管理

- 获取 VM 列表
- 创建 VM
- 启动 / 停止 / 重启 VM
- 任务化执行

## 目录结构

```text
cmd/server/                  启动入口
internal/app/                应用装配
internal/core/               状态存储与任务引擎
internal/modules/            Recovery / Network / VM 服务
internal/native/             C 适配层与 cgo 包装
docs/                        架构与开发文档
scripts/                     演示与部署脚本
deploy/                      部署文件
```

## 主要接口

- `GET /api/health`
- `GET /api/overview`
- `GET /api/tasks`
- `GET /api/recovery`
- `POST /api/recovery/snapshots`
- `POST /api/recovery/restore-file`
- `POST /api/recovery/rollback-volume`
- `GET /api/network`
- `POST /api/network/apply`
- `GET /api/vms`
- `POST /api/vms`
- `POST /api/vms/{id}/power`

## 文档

- [架构说明](/C:/Users/19296/Documents/NAS/docs/architecture.md)
- [开发文档](/C:/Users/19296/Documents/NAS/docs/development.md)
- [Linux 虚拟机 Demo](/C:/Users/19296/Documents/NAS/docs/linux-vm-demo.md)
- [HTTP API 设计与 Demo](/C:/Users/19296/Documents/NAS/docs/http-api-demo.md)

## 后续扩展建议

- 用真实 `Btrfs / ZFS / LVM` 替换恢复模块的模拟能力
- 用真实 `iproute2 / bridge / bond / vlan` 替换网络模块的适配逻辑
- 用真实 `libvirt / qemu-kvm` 替换虚拟机模块的底层执行器
- 把内存态状态替换为持久化数据库和宿主机 Agent

