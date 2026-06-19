package llm

import (
	"d0ctr/bldbr-bot/llm/openai"
	"d0ctr/bldbr-bot/llm/types"
	"d0ctr/bldbr-bot/shared"
	"log/slog"
	"strings"
	"time"
)

func CreateStream(model types.Model, messages []types.Message) (chan types.Message, chan error) {
	req := openai.BuildStream(model, prompt, messages)
	responses := make(chan types.Message)
	errs := make(chan error, 1)

	closeAll := func() {
		close(errs)
		close(responses)
	}

	if stream, err := openai.CreateStream(model, req); err != nil {
		errs <- err
		closeAll()
	} else {
		go func() {
			defer closeAll()
			partialBuilder := strings.Builder{}

			end := false
			for !end {
				var text string

				timer := time.NewTimer(10 * time.Second)
				select {
				case text = <-stream.Final:
					end = true
				case err := <- stream.Err:
					errs <- err
					return
				case delta := <-stream.Delta:
					_, err := partialBuilder.WriteString(delta)
					if err != nil {
						slog.Error("failed to write partial string", shared.ErrAttr(err))
						text = ""
					} else {
						text = partialBuilder.String()
					}
				case <- timer.C:
					slog.Debug("timeout waiting for response in llm/stream.go")
					continue
				}

				if text != "" {
					msg := types.FromText("", "", types.MESSAGE_ROLE_ASSISTANT, text)

					responses <- msg
				} else {
					slog.Debug("empty text, end={}", shared.TemplateAttr(end))
				}
			}
		}()
	}

	return responses, errs
}
