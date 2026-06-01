package openai

import (
	"context"
	"fmt"

	"d0ctr/bldbr-bot/llm/types"
	"d0ctr/bldbr-bot/services"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)


func bareParams(prompt string, model types.Model) responses.ResponseNewParams {
	return responses.ResponseNewParams{
		Instructions: openai.String(prompt),
		Store: openai.Bool(false),
		ToolChoice: responses.ResponseNewParamsToolChoiceUnion{
			OfToolChoiceMode: openai.Opt(responses.ToolChoiceOptionsAuto),
		},
		Reasoning: responses.ReasoningParam{
			Effort: responses.ReasoningEffortLow,
			Summary: openai.ReasoningSummaryAuto,
		},
		Model: model.Name(),
	}
}

func BuildRequest(model types.Model, prompt string, messages []types.Message) responses.ResponseNewParams {
	base := bareParams(prompt, model)
	switch (model.Provider()) {
		case types.MODEL_PROVIDER_OPENAI: {
			return buildRequestGpt(base, messages)
		}
		case types.MODEL_PROVIDER_XAI: {
			return buildRequestGrok(base, messages)
		}
	}
	panic("unreachable")
}

func buildRequestGpt(base responses.ResponseNewParams, messages []types.Message) responses.ResponseNewParams {
	// tools
	base.Tools = []responses.ToolUnionParam{
		responses.ToolParamOfWebSearch(responses.WebSearchToolTypeWebSearch),
	}
	base.ToolChoice = responses.ResponseNewParamsToolChoiceUnion{
		OfToolChoiceMode: openai.Opt(responses.ToolChoiceOptionsAuto),
	}

	// input
	base.Input = mapMessages(messages)

	return base
}

func buildRequestGrok(base responses.ResponseNewParams, messages []types.Message) responses.ResponseNewParams {
	// tools
	xSearchTool := responses.ToolParamOfCustom("x_search")
	xSearchTool.OfCustom.Type = "x_search"

	webSearchTool := responses.ToolParamOfCustom("web_search")
	webSearchTool.OfCustom.Type = "web_search"

	base.Tools = []responses.ToolUnionParam{ xSearchTool, webSearchTool }

	// input
	base.Input = mapMessages(messages)

	return base
}

func SendRequest(model types.Model, request responses.ResponseNewParams) (Response, error) {
	var client *openai.Client

	switch model.Provider() {
		case types.MODEL_PROVIDER_XAI: client = services.XAi();
		case types.MODEL_PROVIDER_OPENAI: client = services.OpenAi()
	}

	if client == nil {
		return Response{}, fmt.Errorf("no client for model [%v]", model.Name())
	}

	oResponse, err := client.Responses.New(context.Background(), request)
	if err != nil {
		err = fmt.Errorf("api error: %v", err)
	}

	return Response{ oResponse }, err
}
