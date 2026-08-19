package openai

import (
	"context"
	"fmt"

	"d0ctr/bldbr-bot/llm/types"
	"d0ctr/bldbr-bot/services"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)


func bareParams(prompt string, model types.Model, cursor string) responses.ResponseNewParams {
	params := responses.ResponseNewParams{
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

	if cursor != "" {
		params.PreviousResponseID = openai.String(cursor)
	} else {
		params.Instructions = openai.String(prompt)
	}

	return params
}

func BuildRequest(model types.Model, prompt string, messages []types.Message, cursor string) responses.ResponseNewParams {
	params := bareParams(prompt, model, cursor)
	
	switch (model.Provider()) {
		case types.MODEL_PROVIDER_OPENAI: {
			params = buildRequestGpt(params, messages)
		}
		case types.MODEL_PROVIDER_XAI: {
			params = buildRequestGrok(params, messages)
		}
		default: panic("unreachable")
	}

	return params
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

func SendRequest(model types.Model, request responses.ResponseNewParams, conversationId string) (Response, error) {
	var client *openai.Client

	var options []option.RequestOption

	switch model.Provider() {
		case types.MODEL_PROVIDER_XAI: {
			client = services.XAi();
			options = append(options, option.WithHeader("x-grok-conv-id", conversationId))
		}

		case types.MODEL_PROVIDER_OPENAI: client = services.OpenAi()
	}

	if client == nil {
		return Response{}, fmt.Errorf("no client for model [%v]", model.Name())
	}

	oResponse, err := client.Responses.New(context.Background(), request, options...)
	if err != nil {
		err = fmt.Errorf("api error: %v", err)
	}

	return Response{ oResponse }, err
}
