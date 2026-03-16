package shared

import (
	"os"
	"encoding/json"
	_ "embed" 
	"fmt"

	"github.com/dotenv-org/godotenvvault"
)

var (
	AHEGAO_API = newConfigKey("AHEGAO_API", "")
	URBAN_API = newConfigKey("URBAN_API", "")
)

type ConfigKey[T any] struct {
	key string
	def T
}

func newConfigKey[T any](key string, def T) ConfigKey[T] {
	return ConfigKey[T]{
		key,
		def,
	}
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

var _ = func () any {
	if err := godotenvvault.Load(".env.staging"); err != nil {
		panic(fmt.Errorf("failed to acquire env variables: %w", err));
	}

	return nil
}()

type EnvKey string
const (
	REDIS_URL EnvKey = "REDIS_URL"
	TELEGRAM_TOKEN EnvKey = "TELEGRAM_TOKEN"
	ENV EnvKey = "ENV"
)

func (key EnvKey) Get() (string, bool) {
	v := os.Getenv(string(key))
	if (v == "") {
		return v, false
	}

	return v, true
}

