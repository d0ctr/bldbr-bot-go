package shared

import (
	"encoding/json"
	_ "embed" 
	"log/slog"
)

type ConfigKey = string
const (
	AHEGAO_API ConfigKey = "AHEGAO_API"
	URBAN_API ConfigKey = "URBAN_API"
)

//go:embed config.json
var configBytes []byte

var config map[string]any

func GetValue[T any](key ConfigKey, def T) (T, bool) {
	if config == nil {
		loadConfig()
	}

	v, ok := config[key]
	if (v == nil || !ok) {
		return def, false
	}

	return v.(T), true
}

func loadConfig() {
	logger := slog.Default().With("component", "config")

	if err := json.Unmarshal(configBytes, &config); err != nil {
		logger.Error("failed to load config", err)
	}
}
