package shared

import (
	"encoding/json"
	_ "embed" 
)

type ConfigKey = string
const (
	AHEGAO_API ConfigKey = "AHEGAO_API"
	URBAN_API  ConfigKey = "URBAN_API"
)

//go:embed config.json
var _CONFIG_B []byte
var _CONFIG   map[string]any = loadConfig()

func GetValue[T any](key ConfigKey, def T) (T, bool) {
	v, ok := _CONFIG[key]
	if (v == nil || !ok) {
		return def, false
	}

	return v.(T), true
}

func loadConfig() map[string]any {
	var config map[string]any
	if err := json.Unmarshal(_CONFIG_B, &config); err != nil {
		panic(err)
	}

	return config
}
