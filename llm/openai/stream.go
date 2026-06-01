package openai

import (
	"context"
	"d0ctr/bldbr-bot/llm/types"
	"d0ctr/bldbr-bot/services"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

type Stream struct {
	Err chan error
	Delta chan string
	Final chan string
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

func BuildStream(model types.Model, prompt string, messages []types.Message) responses.ResponseNewParams {
	base := bareParams(prompt, model)
	switch (model.Provider()) {
		case types.MODEL_PROVIDER_OPENAI: {
			return buildStreamGpt(base, messages)
		}
		case types.MODEL_PROVIDER_XAI: {
			return buildStreamGrok(base, messages)
		}
	}
	panic("unreachable")
}



func CreateStream(model types.Model, request responses.ResponseNewParams) (Stream, error) {
	var client *openai.Client

	switch model.Provider() {
		case types.MODEL_PROVIDER_OPENAI: client = services.OpenAi()
		case types.MODEL_PROVIDER_XAI: client = services.XAi();
	}

	if client == nil {
		return Stream{}, fmt.Errorf("no client for model [%v]", model.Name())
	}

	oStream := client.Responses.NewStreaming(context.Background(), request)

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
		}()

		for oStream.Next() {
			switch cur := oStream.Current().AsAny().(type) {
			case responses.ResponseTextDeltaEvent:
				stream.Delta <- cur.Delta
			case responses.ResponseTextDoneEvent:
				stream.Final <- cur.Text
			}
		}

		err := oStream.Err()
		if err != nil {
			stream.Err <- err
		}
	}()

	return stream, nil
}
