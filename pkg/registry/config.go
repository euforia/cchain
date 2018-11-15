package registry

// Config is the configuration used to initialize a registry
type Config struct {
	Provider string
	Addr     string
	Config   map[string]interface{}
}
