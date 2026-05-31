package storage

import "github.com/HREOLIN/NAS/internal/core/store"

type Service struct {
	store *store.Store
}

type SummaryResponse struct {
	Disks   []store.Disk   `json:"disks"`
	Volumes []store.Volume `json:"volumes"`
}

func NewService(state *store.Store) *Service {
	return &Service{store: state}
}

func (s *Service) Summary() SummaryResponse {
	return SummaryResponse{
		Disks:   s.store.GetDisks(),
		Volumes: s.store.GetVolumes(),
	}
}
