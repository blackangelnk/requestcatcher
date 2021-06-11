package config

const (
	defaultClientPort  = 8099
	defaultCatcherPort = 8098
)

type Configuration struct {
	ClientPort  int
	CatcherPort int
}

func Init() *Configuration {
	return &Configuration{
		ClientPort:  defaultClientPort,
		CatcherPort: defaultCatcherPort,
	}
}
