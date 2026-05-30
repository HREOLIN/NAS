package vm

import (
	"net/http"
	"time"

	"github.com/HREOLIN/NAS/internal/core/store"
	"github.com/HREOLIN/NAS/internal/core/task"
	"github.com/HREOLIN/NAS/internal/native"
)

type Service struct {
	store  *store.Store
	engine *task.Engine
}

type CreateVMRequest struct {
	Name      string `json:"name"`
	CPU       int    `json:"cpu"`
	MemoryGB  int    `json:"memoryGb"`
	DiskGB    int    `json:"diskGb"`
	NetworkID string `json:"networkId"`
	Template  string `json:"template"`
}

type PowerActionRequest struct {
	Action string `json:"action"`
}

type SummaryResponse struct {
	Templates []store.VMTemplate `json:"templates"`
	VMs       []store.VM         `json:"vms"`
}

func NewService(state *store.Store, engine *task.Engine) *Service {
	return &Service{
		store:  state,
		engine: engine,
	}
}

func (s *Service) Summary() SummaryResponse {
	return SummaryResponse{
		Templates: s.store.GetTemplates(),
		VMs:       s.store.GetVMs(),
	}
}

func (s *Service) CreateVM(req CreateVMRequest) (*store.Task, int) {
	if req.Name == "" {
		req.Name = "demo-vm"
	}
	if req.CPU <= 0 {
		req.CPU = 2
	}
	if req.MemoryGB <= 0 {
		req.MemoryGB = 4
	}
	if req.DiskGB <= 0 {
		req.DiskGB = 80
	}
	if req.NetworkID == "" {
		req.NetworkID = "eth0"
	}
	if req.Template == "" {
		req.Template = "Ubuntu 24.04 Base"
	}

	taskResult := s.engine.Run("vm", "create_vm", req, func(ctx *task.Context) (any, error) {
		ctx.Step("reserve-resources", "Reserve CPU and memory")
		ctx.Step("create-disk", "Create virtual disk")
		detail, err := native.CreateVM(req.Name, req.CPU, req.MemoryGB, req.DiskGB, req.NetworkID)
		if err != nil {
			return nil, ctx.Fail(err.Error(), false)
		}
		ctx.Step("attach-network", detail)
		ctx.Step("register-vm", "Register VM definition")

		vm := store.VM{
			ID:        s.store.NextID("vm"),
			Name:      req.Name,
			CPU:       req.CPU,
			MemoryGB:  req.MemoryGB,
			DiskGB:    req.DiskGB,
			NetworkID: req.NetworkID,
			Status:    "stopped",
			Template:  req.Template,
			CreatedAt: time.Now().UTC(),
		}
		s.store.AddVM(vm)
		return vm, nil
	})

	return taskResult, statusCode(taskResult)
}

func (s *Service) ChangePower(id string, req PowerActionRequest) (*store.Task, int) {
	taskResult := s.engine.Run("vm", "power_"+req.Action, req, func(ctx *task.Context) (any, error) {
		vmItem, ok := s.store.GetVM(id)
		if !ok {
			return nil, ctx.Fail("vm not found", false)
		}
		ctx.Step("check-vm", "Check VM state")
		detail, err := native.ChangeVMPower(id, req.Action)
		if err != nil {
			return nil, ctx.Fail(err.Error(), false)
		}
		ctx.Step("invoke-hypervisor", detail)

		switch req.Action {
		case "start", "restart":
			vmItem.Status = "running"
		case "stop":
			vmItem.Status = "stopped"
		default:
			return nil, ctx.Fail("unsupported power action", false)
		}

		if ok := s.store.UpdateVM(vmItem); !ok {
			return nil, ctx.Fail("failed to update vm state", false)
		}

		ctx.Step("refresh-status", "Refresh VM runtime state")
		return vmItem, nil
	})

	return taskResult, statusCode(taskResult)
}

func statusCode(task *store.Task) int {
	if task.Status == "failed" {
		return http.StatusBadRequest
	}
	return http.StatusAccepted
}
