package code

import "github.com/weirwei/codereview/llm"

type CodePatch struct {
	Filepaths  []string
	Content    string
	MatchRules llm.PromptConfig
}
