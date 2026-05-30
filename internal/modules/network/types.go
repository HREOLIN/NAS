package network

type DNSConfig struct {
	Servers []string `json:"servers"`
}

type InterfaceProfile struct {
	NicID   string    `json:"nicId"`
	Mode    string    `json:"mode"`
	IPv4    string    `json:"ipv4"`
	Mask    string    `json:"mask"`
	Gateway string    `json:"gateway"`
	DNS     DNSConfig `json:"dns"`
}

type ExecutionPlan struct {
	Platform         string   `json:"platform"`
	InterfaceID      string   `json:"interfaceId"`
	Commands         []string `json:"commands"`
	ProbeTargets     []string `json:"probeTargets"`
	RollbackCommands []string `json:"rollbackCommands"`
}

type ApplyResult struct {
	Profile InterfaceProfile `json:"profile"`
	Plan    ExecutionPlan    `json:"plan"`
	Detail  string           `json:"detail"`
}
