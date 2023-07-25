package internal

import "github.com/caarlos0/env"

type Configuration struct {
	AppStoreAppVersionsSaveChunkSize int    `env:"APP_STORE_APPLICATION_VERSIONS_SAVE_CHUNK_SIZE" envDefault:"20"`
	ChartProviderId                  string `env:"CHART_PROVIDER_ID" envDefault:"*"`
	IsOCIRegistry                    bool   `env:"IS_OCI_REGISTRY" envDefault:"true"`
}

func ParseConfiguration() (*Configuration, error) {
	cfg := &Configuration{}
	err := env.Parse(cfg)
	return cfg, err
}
