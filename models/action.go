package models

// Action Struct
type Action struct {
	Type        string            `json:"type"`
	Command     string            `json:"command"`
	OutputToVar bool              `json:"output_to_var"`
	RunOnFail   bool              `json:"run_on_fail"`
	LastConfig  bool              `json:"last_config"`
	Config      map[string]string `json:"config"`
}
