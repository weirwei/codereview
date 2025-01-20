package review

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/sashabaranov/go-openai"
	"github.com/weirwei/codereview/code"
	"github.com/weirwei/codereview/llm"
	"github.com/weirwei/codereview/log"
	"github.com/weirwei/codereview/utils"
)

type implLLM struct {
	ctx       context.Context
	codePatch code.CodePatch
	language  string
	maxToken  int
	model     string
	sk        string
	baseUrl   string
	handler   func(context.Context, string) error
}

func (r *implLLM) SetHandler(handler func(context.Context, string) error) {
	r.handler = handler
}
func (r *implLLM) SetCodePatch(codePatch code.CodePatch) {
	r.codePatch = codePatch
}

func NewDefaultReviewer(ctx context.Context, baseUrl, sk, model string, maxToken int, language string) *implLLM {
	return newReviewer(ctx, code.CodePatch{}, language, maxToken, model, sk, baseUrl, nil)
}

func newReviewer(ctx context.Context, codePatch code.CodePatch, language string, maxToken int, model string, sk string, baseUrl string, handler func(context.Context, string) error) *implLLM {
	return &implLLM{
		ctx:       ctx,
		codePatch: codePatch,
		language:  language,
		maxToken:  maxToken,
		model:     model,
		sk:        sk,
		baseUrl:   baseUrl,
		handler:   handler,
	}
}

func (r *implLLM) getPrompt() string {
	return llm.NewPrompt(r.codePatch.MatchRules)
}

func (r *implLLM) check() error {
	prompt := r.getPrompt()
	codePatchTokens := utils.EstimateTokens(r.codePatch.Content)
	promptTokens := utils.EstimateTokens(prompt)
	if codePatchTokens+promptTokens > r.maxToken {
		return fmt.Errorf("token over, except:%d, codePatch:%d, llm.prompt:%d", r.maxToken, codePatchTokens, promptTokens)
	}
	if r.handler == nil {
		return fmt.Errorf("handler is nil")
	}
	return nil
}

func (r *implLLM) Exec() error {
	if err := r.check(); err != nil {
		return err
	}
	config := openai.DefaultConfig(r.sk)
	config.BaseURL = r.baseUrl
	stream, err := llm.CreateChatCompletionStream(r.ctx, config, openai.ChatCompletionRequest{
		Model: r.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: r.getPrompt()},
			{Role: openai.ChatMessageRoleUser, Content: r.codePatch.Content},
		},
	})
	if err != nil {
		return err
	}
	for {
		response, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Errorf("stream error:%v", err)
			break
		}
		if len(response.Choices) == 0 {
			break
		}
		if response.Choices[0].FinishReason == "stop" {
			break
		}
		if err := r.handler(r.ctx, response.Choices[0].Delta.Content); err != nil {
			log.Errorf("handler error:%v", err)
			break
		}
	}
	return nil
}
