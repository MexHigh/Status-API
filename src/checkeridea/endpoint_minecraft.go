package checkeridea

type MinecraftEndpoint struct {
	defaultEndpoint
	Config struct {
		URL string `json:"url"`
	} `json:"-"`
}

func (e *MinecraftEndpoint) SetDefaults() {
	return
}

func (e *MinecraftEndpoint) Check() error {
	return nil
}

func (e *MinecraftEndpoint) Status() EndpointStatus {
	return e.status
}