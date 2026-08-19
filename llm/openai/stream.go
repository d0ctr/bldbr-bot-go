package openai

import (
	"context"
	"d0ctr/bldbr-bot/llm/types"
	"d0ctr/bldbr-bot/services"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

type Stream struct {
	Err chan error
	Delta chan string
	Final chan string
	ResponseId chan string
}

func buildStreamGpt(base responses.ResponseNewParams, messages []types.Message) responses.ResponseNewParams {
	// tools
	base.Tools = []responses.ToolUnionParam{
		responses.ToolParamOfWebSearch(responses.WebSearchToolTypeWebSearch),
	}
	base.ToolChoice = responses.ResponseNewParamsToolChoiceUnion{
		OfToolChoiceMode: openai.Opt(responses.ToolChoiceOptionsAuto),
	}

	// streaming
	base.StreamOptions = responses.ResponseNewParamsStreamOptions{}
	

	// input
	base.Input = mapMessages(messages)

	return base
}

func buildStreamGrok(base responses.ResponseNewParams, messages []types.Message) responses.ResponseNewParams {
	// tools
	base.Tools = []responses.ToolUnionParam{
		responses.ToolParamOfWebSearch(responses.WebSearchToolTypeWebSearch),
	}
	base.ToolChoice = responses.ResponseNewParamsToolChoiceUnion{
		OfToolChoiceMode: openai.Opt(responses.ToolChoiceOptionsAuto),
	}

	// streaming
	base.StreamOptions = responses.ResponseNewParamsStreamOptions{}
	

	// input
	base.Input = mapMessages(messages)

	return base
}

func BuildStream(model types.Model, prompt string, messages []types.Message, cursor string) responses.ResponseNewParams {
	params := bareParams(prompt, model)
	switch (model.Provider()) {
		case types.MODEL_PROVIDER_OPENAI: {
			params = buildStreamGpt(params, messages)
		}
		case types.MODEL_PROVIDER_XAI: {
			params = buildStreamGrok(params, messages)
		}
		default: panic("unreachable")
	}

	if cursor != "" {
		params.PreviousResponseID = openai.String(cursor)
	}

	return params
}



func CreateStream(model types.Model, request responses.ResponseNewParams, conversationId string) (Stream, error) {
	var client *openai.Client
	var options []option.RequestOption

	switch model.Provider() {
		case types.MODEL_PROVIDER_OPENAI: {
			client = services.OpenAi()
		}
		case types.MODEL_PROVIDER_XAI: {
			client = services.XAi();
			options = append(options, option.WithHeader("x-grok-conv-id", conversationId))
		}
	}

	if client == nil {
		return Stream{}, fmt.Errorf("no client for model [%v]", model.Name())
	}

	oStream := client.Responses.NewStreaming(context.Background(), request, options...)

	stream := Stream{
		Err: make(chan error, 1),
		Delta: make(chan string),
		Final: make(chan string),
	}

	go func() {
		defer func() {
			oStream.Close()

			close(stream.Delta)
			close(stream.Err)
			close(stream.Final)
			close(stream.ResponseId)
		}()

		for oStream.Next() {
			switch cur := oStream.Current().AsAny().(type) {
				case responses.ResponseTextDeltaEvent: stream.Delta <- cur.Delta
				case responses.ResponseTextDoneEvent: stream.Final <- cur.Text
			}

			responseId := oStream.Current().AsResponseCompleted().Response.ID
			if responseId != "" {
				stream.ResponseId <- responseId
			}
		}

		err := oStream.Err()
		if err != nil {
			stream.Err <- err
		}
	}()

	return stream, nil
}
