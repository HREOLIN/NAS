package app

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/HREOLIN/NAS/internal/core/store"
	"github.com/HREOLIN/NAS/internal/core/task"
	"github.com/HREOLIN/NAS/internal/modules/network"
	"github.com/HREOLIN/NAS/internal/modules/recovery"
	"github.com/HREOLIN/NAS/internal/modules/storage"
	"github.com/HREOLIN/NAS/internal/modules/vm"
)

type Server struct {
	mux             *http.ServeMux
	state           *store.Store
	recoveryService *recovery.Service
	networkService  *network.Service
	storageService  *storage.Service
	vmService       *vm.Service
}

func NewServer(addr string) *http.Server {
	state := store.New()
	engine := task.NewEngine(state)

	server := &Server{
		mux:             http.NewServeMux(),
		state:           state,
		recoveryService: recovery.NewService(state, engine),
		networkService:  network.NewService(state, engine),
		storageService:  storage.NewService(state),
		vmService:       vm.NewService(state, engine),
	}

	server.routes()

	return &http.Server{
		Addr:    addr,
		Handler: server.mux,
	}
}

func (s *Server) routes() {
	s.mux.HandleFunc("/api/health", s.handleHealth)
	s.mux.HandleFunc("/api/overview", s.handleOverview)
	s.mux.HandleFunc("/api/tasks", s.handleTasks)
	s.mux.HandleFunc("/api/recovery", s.handleRecovery)
	s.mux.HandleFunc("/api/recovery/snapshots", s.handleCreateSnapshot)
	s.mux.HandleFunc("/api/recovery/restore-file", s.handleRestoreFile)
	s.mux.HandleFunc("/api/recovery/rollback-volume", s.handleRollbackVolume)
	s.mux.HandleFunc("/api/network", s.handleNetwork)
	s.mux.HandleFunc("/api/network/apply", s.handleApplyNetwork)
	s.mux.HandleFunc("/api/storage", s.handleStorage)
	s.mux.HandleFunc("/api/vms", s.handleVMs)
	s.mux.HandleFunc("/api/vms/", s.handleVMPower)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":        true,
		"startedAt": s.state.StartedAt,
	})
}

func (s *Server) handleOverview(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"startedAt": s.state.StartedAt,
		"modules": map[string]any{
			"snapshots":   len(s.state.GetSnapshots()),
			"recoveries":  len(s.state.GetRecoveries()),
			"interfaces":  len(s.state.GetInterfaces()),
			"disks":       len(s.state.GetDisks()),
			"volumes":     len(s.state.GetVolumes()),
			"runningVms":  s.state.CountRunningVMs(),
			"totalVms":    len(s.state.GetVMs()),
			"tasks":       len(s.state.GetTasks()),
			"vmTemplates": len(s.state.GetTemplates()),
		},
	})
}

func (s *Server) handleTasks(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.state.GetTasks())
}

func (s *Server) handleRecovery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"message": "method not allowed"})
		return
	}
	writeJSON(w, http.StatusOK, s.recoveryService.Summary())
}

func (s *Server) handleCreateSnapshot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"message": "method not allowed"})
		return
	}
	var req recovery.CreateSnapshotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}
	task, status := s.recoveryService.CreateSnapshot(req)
	writeJSON(w, status, task)
}

func (s *Server) handleRestoreFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"message": "method not allowed"})
		return
	}
	var req recovery.RestoreFileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}
	task, status := s.recoveryService.RestoreFile(req)
	writeJSON(w, status, task)
}

func (s *Server) handleRollbackVolume(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"message": "method not allowed"})
		return
	}
	var req recovery.RollbackVolumeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}
	task, status := s.recoveryService.RollbackVolume(req)
	writeJSON(w, status, task)
}

func (s *Server) handleNetwork(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"message": "method not allowed"})
		return
	}
	writeJSON(w, http.StatusOK, s.networkService.Summary())
}

func (s *Server) handleApplyNetwork(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"message": "method not allowed"})
		return
	}
	var req network.ApplyProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}
	task, status := s.networkService.ApplyProfile(req)
	writeJSON(w, status, task)
}

func (s *Server) handleStorage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"message": "method not allowed"})
		return
	}
	writeJSON(w, http.StatusOK, s.storageService.Summary())
}

func (s *Server) handleVMs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, s.vmService.Summary())
	case http.MethodPost:
		var req vm.CreateVMRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
			return
		}
		task, status := s.vmService.CreateVM(req)
		writeJSON(w, status, task)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"message": "method not allowed"})
	}
}

func (s *Server) handleVMPower(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"message": "method not allowed"})
		return
	}
	trimmed := strings.TrimPrefix(r.URL.Path, "/api/vms/")
	parts := strings.Split(trimmed, "/")
	if len(parts) != 2 || parts[1] != "power" {
		writeJSON(w, http.StatusNotFound, map[string]string{"message": "not found"})
		return
	}
	var req vm.PowerActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}
	task, status := s.vmService.ChangePower(parts[0], req)
	writeJSON(w, status, task)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
