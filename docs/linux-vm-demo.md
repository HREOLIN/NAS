# Linux 虚拟机作为 NAS 后端服务器 Demo

这份 demo 不是讲抽象设计，而是直接告诉你：

- 如果你有一台自己的 Linux 虚拟机
- 想把它当成 NAS 后端服务节点
- 自己维护数据恢复、网络配置、虚拟机管理这些能力

应该怎么搭起来。

## 1. 适合你的使用方式

你现在的思路很适合先这样做：

1. 准备一台 Linux 虚拟机
2. 在虚拟机里运行当前这个 Go 后端服务
3. 通过 HTTP API 来管理 NAS 核心功能
4. 后续逐步把 demo 里的模拟能力替换成真实 Linux 系统能力

也就是说，这台 Linux 虚拟机先不是“文件共享界面”，而是你的：

- NAS 控制后端
- 存储恢复实验环境
- 网络配置实验环境
- 虚拟机管理实验环境

## 2. 推荐环境

建议 Linux 虚拟机配置：

- `Ubuntu 22.04 / 24.04` 或 `Debian 12`
- `2 vCPU`
- `4 GB RAM`
- `40 GB+` 系统盘
- 如果要做恢复实验，额外挂一块数据盘更好

建议安装：

- `git`
- `curl`
- `rsync`
- `build-essential`
- `golang`

Ubuntu / Debian 示例：

```bash
sudo apt update
sudo apt install -y git curl rsync build-essential golang
```

## 3. 项目放置方式

建议把项目部署到：

```bash
/opt/nas-backend
```

你可以把仓库代码拷到虚拟机，或者直接在虚拟机里 clone。

## 4. 快速启动

在 Linux 虚拟机里执行：

```bash
export GOCACHE=/opt/nas-backend/.gocache
export GOMODCACHE=/opt/nas-backend/.gomodcache
cd /opt/nas-backend
go run ./cmd/server
```

默认监听：

```text
http://127.0.0.1:8080
```

如果你想给局域网其他机器访问，可以后面把监听地址改成 `0.0.0.0:8080`。

## 5. 一键部署为系统服务

仓库里已经准备好了：

- `deploy/linux/nas-backend.service`
- `scripts/setup-linux-vm.sh`

在 Linux 虚拟机里执行：

```bash
sudo bash scripts/setup-linux-vm.sh
```

它会做这些事：

1. 把项目同步到 `/opt/nas-backend`
2. 创建 `.gocache` 和 `.gomodcache`
3. 安装 `systemd` 服务
4. 启动 `nas-backend.service`

查看状态：

```bash
sudo systemctl status nas-backend.service
```

查看日志：

```bash
sudo journalctl -u nas-backend.service -f
```

## 6. 直接跑一套维护 demo

仓库里准备了一个完整演示脚本：

```bash
scripts/demo-linux-vm.sh
```

它会依次演示：

1. 健康检查
2. 查看总览
3. 创建恢复快照
4. 模拟文件恢复
5. 应用安全网络配置
6. 模拟失败网络配置并触发回滚
7. 创建虚拟机
8. 启动虚拟机
9. 查看任务列表

运行方式：

```bash
chmod +x scripts/demo-linux-vm.sh
./scripts/demo-linux-vm.sh
```

如果服务不在本机，而是在另一台 Linux 虚拟机：

```bash
BASE_URL=http://192.168.10.50:8080 ./scripts/demo-linux-vm.sh
```

## 7. 你可以怎么“自己维护 NAS 功能”

这部分是重点。

你可以先把这台 Linux 虚拟机当作“NAS 控制面”，逐步接入真实能力。

### 数据恢复

你现在先用当前接口做流程打通：

- 创建快照
- 模拟恢复
- 记录恢复任务

后续替换方向：

- `Btrfs subvolume snapshot`
- `ZFS snapshot / rollback`
- `LVM snapshot`

你后面真正要做的，就是把：

- [internal/modules/recovery/service.go](/C:/Users/19296/Documents/NAS/internal/modules/recovery/service.go:1)
- [internal/native/native.c](/C:/Users/19296/Documents/NAS/internal/native/native.c:1)

里的模拟逻辑，换成真实的快照和恢复命令。

### 网络配置

你已经有了一版更接近真实 NAS 的结构：

- 先生成执行计划
- 再生成回滚计划
- 再执行
- 再探活
- 失败回滚

重点代码：

- [internal/modules/network/service.go](/C:/Users/19296/Documents/NAS/internal/modules/network/service.go:1)
- [internal/modules/network/types.go](/C:/Users/19296/Documents/NAS/internal/modules/network/types.go:1)

你后面可以接：

- `ip addr`
- `ip route`
- `bridge`
- `bond`
- `vlan`

这样你的 Linux 虚拟机就真的开始具备 NAS 网络配置能力了。

### 虚拟机管理

如果你的这台 Linux 虚拟机所在宿主机支持虚拟化，你可以继续往这条路扩：

- `libvirt`
- `qemu-kvm`

当前代码位置：

- [internal/modules/vm/service.go](/C:/Users/19296/Documents/NAS/internal/modules/vm/service.go:1)

你现在先用它把：

- 创建 VM
- 启停 VM
- 任务记录

这条控制流走通。

以后再把 native 层换成真实的 `virsh` / `libvirt` 调用。

## 8. 最推荐你的实践顺序

如果你是自己维护 NAS 功能，我建议你按这个顺序做：

1. 先把这套后端在 Linux 虚拟机里跑起来
2. 先跑 `demo-linux-vm.sh` 看三大模块的接口行为
3. 先把 `network` 模块改成真实 Linux 命令执行
4. 再把 `recovery` 模块接真实快照
5. 最后再接 `vm` 模块到 `libvirt`

原因很简单：

- `network` 最容易先做出真实效果
- `recovery` 是 NAS 真正核心价值
- `vm` 最后接，技术门槛和环境依赖更高

## 9. 你实际怎么用这套 demo

最简单的方式就是：

1. Linux 虚拟机跑服务
2. 你自己写脚本调用接口
3. 先把“控制逻辑”做对
4. 再把底层能力逐步接真

也就是说，这个 demo 现在最适合你做两件事：

- 当 NAS 后端框架原型
- 当你自己维护 NAS 功能的实验控制台

## 10. 下一步建议

如果你愿意，我下一步可以继续直接帮你做两种 demo 里的一个：

1. 把 `network` 接成真实 Linux 命令版本
2. 把 `recovery` 接成真实 `Btrfs` 快照版本

如果是为了最先看到真实效果，我建议先做 `network` 真实版。  
如果是为了更像 NAS 本体，我建议先做 `recovery + Btrfs`。

