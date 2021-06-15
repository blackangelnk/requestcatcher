package config

const (
	defaultClientPort    = 8099
	defaultCatcherPort   = 8098
	StorageTypeMemory    = "memory"
	StorageTypeFile      = "file"
	StorageTypeDummy     = "dummy"
	DefaultStorage       = StorageTypeMemory
	DefaultMemStorageCap = 5
)

type MemStorageConfig struct {
	Cap int
}

type FileStorageConfig struct {
	Path string
}

type Configuration struct {
	ClientPort        int
	CatcherPort       int
	StorageType       string
	FileStorageConfig *FileStorageConfig
	MemStorageConfig  *MemStorageConfig
}

func Init() *Configuration {
	return &Configuration{
		ClientPort:  defaultClientPort,
		CatcherPort: defaultCatcherPort,
		StorageType: DefaultStorage,
		MemStorageConfig: &MemStorageConfig{
			Cap: DefaultMemStorageCap,
		},
	}
}
