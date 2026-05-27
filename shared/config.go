package shared

import (
	_ "embed"
	"encoding/json"
	"log/slog"

	"github.com/dotenv-org/godotenvvault"
)

var (
	AHEGAO_API = newConfigKey("AHEGAO_API", "")
	URBAN_API  = newConfigKey("URBAN_API", "")
)

type ConfigKey[T any] struct {
	key string
	def T
}

func newConfigKey[T any](key string, def T) ConfigKey[T] {
	return ConfigKey[T]{ key, def }
}

//go:embed config.json
var configBytes []byte

var _CONFIG = func() map[string]any {
	var config map[string]any
	if err := json.Unmarshal(configBytes, &config); err != nil {
		panic(err)
	}

	return config
}()


func (key ConfigKey[T]) Get() (T, bool) {
	v, ok := _CONFIG[key.key]
	if (v == nil || !ok) {
		return key.def, false
	}

	return v.(T), true
}

func init() {
	if err := godotenvvault.Load(); err != nil {
		slog.Error("failed to acquire env variables", err)
		// panic();
	}
}

