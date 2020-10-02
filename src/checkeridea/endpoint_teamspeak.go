package checkeridea

type TeamspeakEndpoint struct {
	defaultEndpoint
	Config struct {
		QueryURL string `json:"query_url"`
	} `json:"-"`
}

func (e *TeamspeakEndpoint) SetDefaults() {
	return
}

func (e *TeamspeakEndpoint) Check() error {
	return nil
}

func (e *TeamspeakEndpoint) Status() EndpointStatus {
	return e.status
}