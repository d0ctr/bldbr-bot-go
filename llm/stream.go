package llm

import (
	"d0ctr/bldbr-bot/llm/openai"
	"d0ctr/bldbr-bot/llm/types"
	"d0ctr/bldbr-bot/shared"
	"log/slog"
	"strings"
	"time"
)

func CreateStream(model types.Model, messages []types.Message, cursor string, conversationId string) (chan types.Message, chan string, chan error) {
	// disable conversation-based routing and take any server
	if cursor == "" {
		conversationId = ""
	}

	req := openai.BuildStream(model, prompt, messages, cursor)
	responses := make(chan types.Message)
	errs := make(chan error, 1)
	cursorChan := make(chan string, 1)

	closeAll := func() {
		close(errs)
		close(responses)
		close(cursorChan)
	}

	if stream, err := openai.CreateStream(model, req, conversationId); err != nil {
		errs <- err
		closeAll()
	} else {
		go func() {
			defer closeAll()
			partialBuilder := strings.Builder{}

			end := 0
			for end < 2 {
				var text string

				timer := time.NewTimer(10 * time.Second)
				select {
				case text = <-stream.Final:
					cursorChan <- <- stream.ResponseId
					end += 1
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
				case cursor = <- stream.ResponseId:
					cursorChan <- cursor
					end += 1
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

	return responses, cursorChan, errs
}
