package config

// Config holds all configuration options for connecting to the server
type Config struct {
    // Whether to connect over TLS or not
    UseTLS bool
    // Server to connect to 
    Server string
    // Port to connect on
    Port int
}

// NewConfig creates a new instance of the config class with the default parameters
func NewConfig() *Config {
    return &Config {
        UseTLS: true,
        Server: "pop.gmail.com",
        Port: 995,
    }
}