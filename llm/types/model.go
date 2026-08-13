package types

import (
	"maps"
	"slices"
	"strings"
)

type ModelProvider = string

var (
	MODEL_PROVIDER_OPENAI ModelProvider = "openai"
	MODEL_PROVIDER_XAI ModelProvider = "xai"
)

type Model struct {
	name string
	provider ModelProvider
}

func (m Model) Name() string {
	return m.name
}

func (m Model) Provider() string {
	return m.provider
}

// predifined list although not exclusive

var MODEL_GPT_5_4_NANO Model = Model{ "gpt-5.4-nano", MODEL_PROVIDER_OPENAI }

var MODEL_GROK_4_3 Model = Model{ "grok-4.3", MODEL_PROVIDER_XAI }
var MODEL_GROK_4_6 Model = Model{ "grok-4.6", MODEL_PROVIDER_XAI }

var models = map[string]Model {
	MODEL_GPT_5_4_NANO.name: MODEL_GPT_5_4_NANO,
	MODEL_GROK_4_6.name: MODEL_GROK_4_6,
	MODEL_GROK_4_3.name: MODEL_GROK_4_3,
}

var defModel = MODEL_GPT_5_4_NANO

func GetModelNames() []string {
	return slices.AppendSeq([]string{}, maps.Keys(models))
}

func GetOrDefault(name string) (Model, bool) {
	if name == "" {
		return defModel, true
	}

	model, ok := models[name]

	if ok {
		return model, true
	}

	prefix, _, ok := strings.Cut(name, "-")
	if !ok {
		return Model{}, false
	}

	switch (prefix) {
		case "gpt": return Model{ name, MODEL_PROVIDER_OPENAI }, true;
		case "grok": return Model{ name, MODEL_PROVIDER_XAI }, true;
	}

	return Model{}, false;
}
