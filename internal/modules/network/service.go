package network

import (
	"net"
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

type ApplyProfileRequest struct {
	NicID           string   `json:"nicId"`
	Mode            string   `json:"mode"`
	IPv4            string   `json:"ipv4"`
	Mask            string   `json:"mask"`
	Gateway         string   `json:"gateway"`
	DNS             []string `json:"dns"`
	SimulateFailure bool     `json:"simulateFailure"`
}

type SummaryResponse struct {
	Interfaces []store.NetworkInterface `json:"interfaces"`
	History    []store.NetworkHistory   `json:"history"`
}

func NewService(state *store.Store, engine *task.Engine) *Service {
	return &Service{
		store:  state,
		engine: engine,
	}
}

func (s *Service) Summary() SummaryResponse {
	return SummaryResponse{
		Interfaces: s.store.GetInterfaces(),
		History:    s.store.GetHistory(),
	}
}

func (s *Service) ApplyProfile(req ApplyProfileRequest) (*store.Task, int) {
	if err := validate(req); err != nil {
		return &store.Task{
			Module:   "network",
			Action:   "apply_profile",
			Status:   "failed",
			Progress: 0,
			Error:    err.Error(),
		}, http.StatusBadRequest
	}

	taskResult := s.engine.Run("network", "apply_profile", req, func(ctx *task.Context) (any, error) {
		current, ok := s.store.GetInterface(req.NicID)
		if !ok {
			return nil, ctx.Fail("network interface not found", false)
		}

		before := current
		before.DNS = cloneStrings(current.DNS)
		s.store.AddHistory(store.NetworkHistory{
			At:     time.Now().UTC(),
			NicID:  current.ID,
			Before: before,
		})

		ctx.Step("backup-config", "Backup current network profile")
		ctx.Step("apply-config", "Apply new interface config")

		current.Mode = req.Mode
		if req.Mode == "dhcp" {
			current.IPv4 = "192.168.10.88"
			current.Mask = "255.255.255.0"
			current.Gateway = "192.168.10.1"
			current.DNS = []string{"223.5.5.5"}
		} else {
			current.IPv4 = req.IPv4
			current.Mask = req.Mask
			current.Gateway = req.Gateway
			current.DNS = cloneStrings(req.DNS)
		}

		detail, err := native.ApplyNetworkProfile(req.NicID, req.Mode, current.IPv4, current.Gateway, req.SimulateFailure)
		if err != nil {
			_ = s.store.UpdateInterface(before)
			ctx.Step("rollback", "Rollback to previous network profile")
			return nil, ctx.Fail(err.Error(), true)
		}

		if ok := s.store.UpdateInterface(current); !ok {
			return nil, ctx.Fail("failed to update network state", false)
		}

		ctx.Step("probe-connectivity", detail)
		ctx.Step("confirm", "Confirm management path healthy")
		return current, nil
	})

	return taskResult, statusCode(taskResult)
}

func validate(req ApplyProfileRequest) error {
	if req.NicID == "" {
		return errText("nicId is required")
	}
	if req.Mode != "dhcp" && req.Mode != "static" {
		return errText("mode must be dhcp or static")
	}
	if req.Mode == "static" {
		if net.ParseIP(req.IPv4) == nil {
			return errText("ipv4 must be a valid IPv4 address")
		}
		if net.ParseIP(req.Mask) == nil {
			return errText("mask must be a valid IPv4 address")
		}
		if net.ParseIP(req.Gateway) == nil {
			return errText("gateway must be a valid IPv4 address")
		}
	}
	return nil
}

func statusCode(task *store.Task) int {
	if task.Status == "failed" {
		return http.StatusBadRequest
	}
	return http.StatusAccepted
}

func cloneStrings(input []string) []string {
	out := make([]string, len(input))
	copy(out, input)
	return out
}

type errText string

func (e errText) Error() string {
	return string(e)
}
