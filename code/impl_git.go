package code

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/weirwei/codereview/config"
	"github.com/weirwei/codereview/llm"
	"github.com/weirwei/codereview/log"
	"github.com/weirwei/codereview/utils"
)

type implGit struct {
	pkgs          []string // files in packages
	reviewBranch  string   // branch to review, if empty, use HEAD
	compareBranch string   // branch to compare with, if empty, use origin/master

	maxToken        int              // max token for single request. default 4096
	knowledge       config.Knowledge // knowledge
	filepathFilters []regexp.Regexp  // exclude filepaths
}

type GitCond struct {
	Pkgs            []string
	ReviewBranch    string
	CompareBranch   string
	MaxToken        int
	Knowledge       config.Knowledge
	FilepathFilters []regexp.Regexp
}

func NewGit(cond GitCond) *implGit {
	if cond.MaxToken <= 0 {
		cond.MaxToken = 4096
	}
	if len(cond.CompareBranch) == 0 {
		cond.CompareBranch = "origin/master"
	}

	return &implGit{
		pkgs:            cond.Pkgs,
		reviewBranch:    cond.CompareBranch,
		compareBranch:   cond.ReviewBranch,
		maxToken:        cond.MaxToken,
		knowledge:       cond.Knowledge,
		filepathFilters: cond.FilepathFilters,
	}
}

func (i *implGit) GetCode() ([]CodePatch, error) {
	var (
		maxToken       = i.maxToken
		getFilepath    = i.getFilepath
		getFileContent = i.getFileContent
		getMatchRules  = i.getMatchRules

		codePatch []CodePatch
	)

	filepathList, err := getFilepath()
	if err != nil {
		return nil, err
	}
	if len(filepathList) == 0 {
		return nil, fmt.Errorf("no files found")
	}

	var cp CodePatch
	for _, v := range filepathList {
		contentItem, err := getFileContent(v)
		if err != nil {
			log.Errorf("failed to get file content, file path: %s", v)
			continue
		}

		content := cp.Content + contentItem + "\n"
		rules := getMatchRules(v, content)
		contentToken := utils.EstimateTokens(content)
		promptToken := utils.EstimateTokens(utils.ToJson(rules))
		if contentToken+promptToken > maxToken {
			codePatch = append(codePatch, cp)
			cp = CodePatch{}
		} else {
			cp.Content += contentItem + "\n"
			cp.MatchRules = rules
			cp.Filepaths = append(cp.Filepaths, v)
		}
	}
	if len(cp.Content) > 0 {
		codePatch = append(codePatch, cp)
		cp = CodePatch{}
	}
	return codePatch, nil
}

func (i *implGit) getFilepath() ([]string, error) {
	var (
		pkgs            = i.pkgs
		reviewBranch    = i.reviewBranch
		compareBranch   = i.compareBranch
		filepathFilters = i.filepathFilters
	)

	_, err := utils.ShellExec("git", "fetch")
	if err != nil {
		return nil, err
	}
	args := []string{
		"diff",
		fmt.Sprintf("%s...%s", reviewBranch, compareBranch),
		"--name-only",
		"--diff-filter=d",
	}
	args = append(args, pkgs...)
	result, err := utils.ShellExec("git", args...)
	if err != nil {
		return nil, err
	}
	var filepathList []string
	for _, v := range strings.Fields(result) {
		for _, filter := range filepathFilters {
			if filter.MatchString(v) {
				goto next
			}
		}
		filepathList = append(filepathList, v)
	next:
	}
	return filepathList, nil
}

func (i *implGit) getFileContent(filepath string) (string, error) {
	content, err := utils.ShellExec("git", "diff", fmt.Sprintf("%s...%s", i.compareBranch, i.reviewBranch), filepath)
	return content, err
}

func (i *implGit) getMatchRules(filepath string, content string) llm.PromptConfig {
	var (
		knowledge = i.knowledge
	)
	var pconf llm.PromptConfig
	customRules := knowledge.Custom[utils.GetLangByFilepath(filepath)]
	for _, v := range customRules {
		if v.RegexpF.MatchString(content) {
			pconf.CodeCustoms = append(pconf.CodeCustoms, v.Rules...)
		}
	}
	pconf.TreeCustoms = knowledge.TreeStandard
	return pconf
}
