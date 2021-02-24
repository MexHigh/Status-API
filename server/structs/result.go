package structs

type Status string // can be "up" or "down"

type Misc map[string]string // other informations

type Result struct {
	status Status `json:"status"`
	url    string `json:"url"`
	Misc   `json:",omitempty"`
}
