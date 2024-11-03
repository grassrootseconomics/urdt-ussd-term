package config

import (
	urdtconfig "git.grassecon.net/urdt/ussd/config"
	"git.grassecon.net/urdt/ussd/initializers"
)

var (
	JetstreamURL string
	JetstreamClientName string
)


func LoadConfig() {
	urdtconfig.LoadConfig()

	JetstreamURL = initializers.GetEnv("NATS_JETSTREAM_URL", "localhost:4222")
	JetstreamClientName = initializers.GetEnv("NATS_JETSTREAM_CLIENT_NAME", "omnom")
}
