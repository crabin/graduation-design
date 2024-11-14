package request

type HostForm struct {
	targetHost string `json:"target_host"`
	portState  string `json:"port_state"`
	portEnd    string `json:"port_end"`
}
