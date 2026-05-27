package shared

import (
	"fmt"
	"os"
)

type EnvKey string
const (
	REDIS_URL       EnvKey = "REDIS_URL"
	TELEGRAM_TOKEN  EnvKey = "TELEGRAM_TOKEN"
	ENV             EnvKey = "ENV"
	OPENAI_TOKEN    EnvKey = "OPENAI_TOKEN"
	XAI_TOKEN       EnvKey = "XAI_TOKEN"
	DISCORD_TOKEN   EnvKey = "DISCORD_TOKEN"
	DISCORD_APP_ID  EnvKey = "DISCORD_APP_ID"
	DOMAIN_URL      EnvKey = "DOMAIN_URL"
)

func (key EnvKey) Get() (string, bool) {
	v := os.Getenv(string(key))
	if (v == "") {
		return v, false
	}

	return v, true
}

func (key EnvKey) MustGet() (string) {
	if v, ok := key.Get(); !ok {
		panic(fmt.Errorf("failed to get value of variable '%v' from environment", key))
	} else {
		return v
	}
}

