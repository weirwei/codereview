package llm

import (
	"context"
	"testing"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
)

func TestCreateChatCompletionStream(t *testing.T) {
	viper.SetConfigType("ini")
	viper.SetConfigFile("/Users/weirwei/.codereview")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err := viper.WriteConfigAs(viper.ConfigFileUsed())
			if err != nil {
				t.Error(err)
				return
			}
		} else {
			t.Error(err)
			return
		}
	}
	t.Run("test", func(t *testing.T) {
		ctx := context.Background()
		config := openai.DefaultConfig(viper.GetString("llm.sk"))
		config.BaseURL = viper.GetString("llm.base_url")

		stream, err := CreateChatCompletionStream(ctx, config, openai.ChatCompletionRequest{
			Model:    "Qwen/Qwen2.5-Coder-7B-Instruct",
			Messages: []openai.ChatCompletionMessage{{Role: openai.ChatMessageRoleUser, Content: "hello world"}},
		})
		if err != nil {
			t.Fatal(err)
		}
		for {
			response, err := stream.Recv()
			if err != nil {
				t.Fatal(err)
			}
			if response.Choices[0].FinishReason == "stop" {
				break
			}
			t.Log(response.Choices[0].Delta.Content)
		}
	})
}
