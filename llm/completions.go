package llm

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

func CreateChatCompletionStream(ctx context.Context, config openai.ClientConfig, request openai.ChatCompletionRequest) (stream *openai.ChatCompletionStream, err error) {
	client := openai.NewClientWithConfig(config)
	return client.CreateChatCompletionStream(ctx, request)
}
