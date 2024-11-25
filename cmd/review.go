package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/weirwei/codereview/config"
)

var (
	TokenOver = errors.New("token over")
)

func Review(codePatch string, language string, rules []string) (string, error) {
	Debugf("TreeCustoms:%s", ToJson(config.CRConf.Knowledge.TreeStandard))
	prompt := NewPrompt(PromptConfig{
		TreeCustoms: config.CRConf.Knowledge.TreeStandard,
		CodeCustoms: rules,
	})
	if EstimateTokens(codePatch)+EstimateTokens(prompt) > maxToken {
		return "", fmt.Errorf("token over, except:%d, codePatch:%d, prompt:%d", maxToken, EstimateTokens(codePatch), EstimateTokens(prompt))
	}
	var result string
	// todo request llm api

	// req := CompletionsReq{
	// 	AppId:     "chatatp",
	// 	Model:     model,
	// 	MaxTokens: maxToken,
	// 	Stream:    true,
	// 	Messages: []Message{
	// 		{
	// 			Role:    "system",
	// 			Content: prompt,
	// 		}, {
	// 			Role:    "user",
	// 			Content: codePatch,
	// 		},
	// 	},
	// }
	// Debugf("review req:%s", utils.ToJson(req))
	// result, err = Completions(&req)
	// if err != nil {
	// 	return "", err
	// }
	Debugf("review result:%s", result)
	return ExtractHtmlTagContent("output", result)
}
