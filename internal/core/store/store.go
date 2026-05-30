package store

import (
	"fmt"
	"sync"
	"time"
)

type TaskStep struct {
	Name   string    `json:"name"`
	Detail string    `json:"detail"`
	At     time.Time `json:"at"`
}

type Task struct {
	ID        string     `json:"id"`
	Module    string     `json:"module"`
	Action    string     `json:"action"`
	Status    string     `json:"status"`
	Progress  int        `json:"progress"`
	Steps     []TaskStep `json:"steps"`
	Error     string     `json:"error,omitempty"`
	Rollback  bool       `json:"rollback"`
	Payload   any        `json:"payload,omitempty"`
	Result    any        `json:"result,omitempty"`
	StartedAt time.Time  `json:"startedAt"`
	EndedAt   time.Time  `json:"endedAt,omitempty"`
}

type Snapshot struct {
	ID        string    `json:"id"`
	Volume    string    `json:"volume"`
	Name      string    `json:"name"`
	Files     int       `json:"files"`
	CreatedAt time.Time `json:"createdAt"`
}

type RecoveryRecord struct {
	ID               string    `json:"id"`
	Type             string    `json:"type"`
	SourceSnapshotID string    `json:"sourceSnapshotId"`
	SourcePath       string    `json:"sourcePath,omitempty"`
	TargetPath       string    `json:"targetPath,omitempty"`
	Volume           string    `json:"volume,omitempty"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"createdAt"`
}

type NetworkInterface struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Mode    string   `json:"mode"`
	IPv4    string   `json:"ipv4"`
	Mask    string   `json:"mask"`
	Gateway string   `json:"gateway"`
	DNS     []string `json:"dns"`
	Link    string   `json:"link"`
	Bridge  string   `json:"bridge"`
}

type NetworkHistory struct {
	At     time.Time        `json:"at"`
	NicID  string           `json:"nicId"`
	Before NetworkInterface `json:"before"`
}

type VMTemplate struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type VM struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CPU       int       `json:"cpu"`
	MemoryGB  int       `json:"memoryGb"`
	DiskGB    int       `json:"diskGb"`
	NetworkID string    `json:"networkId"`
	Status    string    `json:"status"`
	Template  string    `json:"template"`
	CreatedAt time.Time `json:"createdAt"`
}

type Store struct {
	mu         sync.RWMutex
	idCounter  int
	StartedAt  time.Time
	tasks      []*Task
	snapshots  []Snapshot
	recoveries []RecoveryRecord
	interfaces []NetworkInterface
	history    []NetworkHistory
	templates  []VMTemplate
	vms        []VM
}

func New() *Store {
	return &Store{
		StartedAt: time.Now().UTC(),
		snapshots: []Snapshot{
			{
				ID:        "snap-system-001",
				Volume:    "volume1",
				Name:      "system-daily-001",
				Files:     1820,
				CreatedAt: time.Date(2026, 5, 29, 22, 0, 0, 0, time.UTC),
			},
			{
				ID:        "snap-project-002",
				Volume:    "volume2",
				Name:      "project-pre-upgrade",
				Files:     946,
				CreatedAt: time.Date(2026, 5, 30, 1, 15, 0, 0, time.UTC),
			},
		},
		interfaces: []NetworkInterface{
			{
				ID:      "eth0",
				Name:    "LAN 1",
				Mode:    "static",
				IPv4:    "192.168.10.20",
				Mask:    "255.255.255.0",
				Gateway: "192.168.10.1",
				DNS:     []string{"223.5.5.5", "8.8.8.8"},
				Link:    "up",
				Bridge:  "br0",
			},
			{
				ID:      "eth1",
				Name:    "LAN 2",
				Mode:    "dhcp",
				IPv4:    "192.168.10.21",
				Mask:    "255.255.255.0",
				Gateway: "192.168.10.1",
				DNS:     []string{"223.5.5.5"},
				Link:    "standby",
				Bridge:  "",
			},
		},
		templates: []VMTemplate{
			{ID: "tpl-ubuntu-24", Name: "Ubuntu 24.04 Base"},
			{ID: "tpl-rocky-9", Name: "Rocky Linux 9 Base"},
		},
		vms: []VM{
			{
				ID:        "vm-demo-001",
				Name:      "dev-jumpbox",
				CPU:       2,
				MemoryGB:  4,
				DiskGB:    60,
				NetworkID: "eth0",
				Status:    "stopped",
				Template:  "Ubuntu 24.04 Base",
				CreatedAt: time.Date(2026, 5, 29, 9, 0, 0, 0, time.UTC),
			},
		},
	}
}

func (s *Store) NextID(prefix string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.idCounter++
	return fmt.Sprintf("%s-%04d", prefix, s.idCounter)
}

func (s *Store) AddTask(task *Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks = append([]*Task{task}, s.tasks...)
	if len(s.tasks) > 50 {
		s.tasks = s.tasks[:50]
	}
}

func (s *Store) GetTasks() []*Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Task, len(s.tasks))
	copy(out, s.tasks)
	return out
}

func (s *Store) GetSnapshots() []Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Snapshot, len(s.snapshots))
	copy(out, s.snapshots)
	return out
}

func (s *Store) AddSnapshot(snapshot Snapshot) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.snapshots = append([]Snapshot{snapshot}, s.snapshots...)
}

func (s *Store) GetRecoveries() []RecoveryRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]RecoveryRecord, len(s.recoveries))
	copy(out, s.recoveries)
	return out
}

func (s *Store) AddRecovery(record RecoveryRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.recoveries = append([]RecoveryRecord{record}, s.recoveries...)
}

func (s *Store) GetInterfaces() []NetworkInterface {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]NetworkInterface, len(s.interfaces))
	copy(out, s.interfaces)
	return out
}

func (s *Store) GetInterface(id string) (NetworkInterface, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, item := range s.interfaces {
		if item.ID == id {
			return item, true
		}
	}
	return NetworkInterface{}, false
}

func (s *Store) UpdateInterface(updated NetworkInterface) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.interfaces {
		if s.interfaces[i].ID == updated.ID {
			s.interfaces[i] = updated
			return true
		}
	}
	return false
}

func (s *Store) AddHistory(history NetworkHistory) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.history = append([]NetworkHistory{history}, s.history...)
}

func (s *Store) GetHistory() []NetworkHistory {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]NetworkHistory, len(s.history))
	copy(out, s.history)
	return out
}

func (s *Store) GetTemplates() []VMTemplate {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]VMTemplate, len(s.templates))
	copy(out, s.templates)
	return out
}

func (s *Store) GetVMs() []VM {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]VM, len(s.vms))
	copy(out, s.vms)
	return out
}

func (s *Store) GetVM(id string) (VM, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, item := range s.vms {
		if item.ID == id {
			return item, true
		}
	}
	return VM{}, false
}

func (s *Store) AddVM(vm VM) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.vms = append([]VM{vm}, s.vms...)
}

func (s *Store) UpdateVM(updated VM) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.vms {
		if s.vms[i].ID == updated.ID {
			s.vms[i] = updated
			return true
		}
	}
	return false
}

func (s *Store) CountRunningVMs() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	count := 0
	for _, item := range s.vms {
		if item.Status == "running" {
			count++
		}
	}
	return count
}
