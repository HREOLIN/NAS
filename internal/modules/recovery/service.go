package recovery

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

type CreateSnapshotRequest struct {
	Volume string `json:"volume"`
	Name   string `json:"name"`
	Files  int    `json:"files"`
}

type RestoreFileRequest struct {
	SnapshotID string `json:"snapshotId"`
	SourcePath string `json:"sourcePath"`
	TargetPath string `json:"targetPath"`
}

type RollbackVolumeRequest struct {
	SnapshotID string `json:"snapshotId"`
	Volume     string `json:"volume"`
}

type SummaryResponse struct {
	Snapshots  []store.Snapshot       `json:"snapshots"`
	Recoveries []store.RecoveryRecord `json:"recoveries"`
}

func NewService(state *store.Store, engine *task.Engine) *Service {
	return &Service{
		store:  state,
		engine: engine,
	}
}

func (s *Service) Summary() SummaryResponse {
	return SummaryResponse{
		Snapshots:  s.store.GetSnapshots(),
		Recoveries: s.store.GetRecoveries(),
	}
}

func (s *Service) CreateSnapshot(req CreateSnapshotRequest) (*store.Task, int) {
	if req.Volume == "" {
		req.Volume = "volume1"
	}
	if req.Name == "" {
		req.Name = "manual-snapshot"
	}

	taskResult := s.engine.Run("recovery", "create_snapshot", req, func(ctx *task.Context) (any, error) {
		ctx.Step("lock-volume", "Lock volume metadata")
		ctx.Step("flush-cache", "Flush pending writes")
		detail, err := native.CreateSnapshot(req.Volume, req.Name)
		if err != nil {
			return nil, ctx.Fail(err.Error(), false)
		}
		ctx.Step("create-snapshot", detail)

		snapshot := store.Snapshot{
			ID:        s.store.NextID("snap"),
			Volume:    req.Volume,
			Name:      req.Name,
			Files:     max(req.Files, 1200),
			CreatedAt: time.Now().UTC(),
		}
		s.store.AddSnapshot(snapshot)
		return snapshot, nil
	})

	return taskResult, statusCode(taskResult)
}

func (s *Service) RestoreFile(req RestoreFileRequest) (*store.Task, int) {
	snapshots := s.store.GetSnapshots()
	if req.SnapshotID == "" && len(snapshots) > 0 {
		req.SnapshotID = snapshots[0].ID
	}
	if req.SourcePath == "" {
		req.SourcePath = "/team/design.docx"
	}
	if req.TargetPath == "" {
		req.TargetPath = "/restore/design.docx"
	}

	taskResult := s.engine.Run("recovery", "restore_file", req, func(ctx *task.Context) (any, error) {
		ctx.Step("validate-snapshot", "Check restore point")
		ctx.Step("prepare-target", "Prepare restore target")
		detail, err := native.RestoreFile(req.SnapshotID, req.SourcePath, req.TargetPath)
		if err != nil {
			return nil, ctx.Fail(err.Error(), false)
		}
		ctx.Step("restore-file", detail)
		ctx.Step("verify-result", "Verify restored file")

		record := store.RecoveryRecord{
			ID:               s.store.NextID("recovery"),
			Type:             "file",
			SourceSnapshotID: req.SnapshotID,
			SourcePath:       req.SourcePath,
			TargetPath:       req.TargetPath,
			Status:           "success",
			CreatedAt:        time.Now().UTC(),
		}
		s.store.AddRecovery(record)
		return record, nil
	})

	return taskResult, statusCode(taskResult)
}

func (s *Service) RollbackVolume(req RollbackVolumeRequest) (*store.Task, int) {
	snapshots := s.store.GetSnapshots()
	if req.SnapshotID == "" && len(snapshots) > 0 {
		req.SnapshotID = snapshots[0].ID
	}
	if req.Volume == "" {
		req.Volume = "volume1"
	}

	taskResult := s.engine.Run("recovery", "rollback_volume", req, func(ctx *task.Context) (any, error) {
		ctx.Step("validate-impact", "Estimate rollback impact")
		ctx.Step("quiesce-services", "Stop dependent services")
		detail, err := native.RollbackVolume(req.SnapshotID, req.Volume)
		if err != nil {
			return nil, ctx.Fail(err.Error(), false)
		}
		ctx.Step("rollback-volume", detail)
		ctx.Step("post-check", "Check volume health")

		record := store.RecoveryRecord{
			ID:               s.store.NextID("recovery"),
			Type:             "volume",
			SourceSnapshotID: req.SnapshotID,
			Volume:           req.Volume,
			Status:           "success",
			CreatedAt:        time.Now().UTC(),
		}
		s.store.AddRecovery(record)
		return record, nil
	})

	return taskResult, statusCode(taskResult)
}

func statusCode(task *store.Task) int {
	if task.Status == "failed" {
		return http.StatusBadRequest
	}
	return http.StatusAccepted
}

func max(value, fallback int) int {
	if value > 0 {
		return value
	}
	return fallback
}
