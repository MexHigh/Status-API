package checkeridea

import (
	"encoding/json"
)

type defaultEndpoint struct {
	FriendlyURL string         `json:"friendly_url"`
	Protocol string `json:"protocol"`
	ConfigRaw json.RawMessage `json:"config"`
	status      EndpointStatus `json:"-"`
}

// EndpointStatus -
type EndpointStatus map[string]string

// JSON return the json endcoded version of EndpointStatus
func (es *EndpointStatus) JSON() ([]byte, error) {
	json, err := json.MarshalIndent(*es, "", "    ")
	if err != nil {
		return nil, err
	}
	return json, nil
}

// Endpoint -
type Endpoint interface {
	// Check should check if the service is online
	// (it usually fills the Status field of defaultEndpoint, too)
	Check() error
	// SetDefaults sets default values for ommited fields
	SetDefaults()
	// Status returns the endpoints status
	Status() EndpointStatus
}
