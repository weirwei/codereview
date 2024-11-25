package cmd

import (
	"fmt"
	"slices"
	"strings"

	"github.com/weirwei/codereview/config"
)

func ReviewCode() string {
	fileNames := filesFromGit()
	var reviewResult string
	fmt.Println("=== code review ===")
	reviewResult = filesReview(fileNames)
	fmt.Println("\rreviewing: 100%")
	fmt.Println("done.")
	return reviewResult
}

func filesFromGit() (fileNames []string) {
	// git fetch
	fmt.Println("=== git fetch ===")
	_, err := ShellExec("git", "fetch")
	if err != nil {
		return
	}
	fmt.Println("done.")
	// git diff
	fmt.Println("=== git diff ===")
	var pkgs []string
	for _, v := range strings.Split(pkg, ",") {
		if len(v) > 0 {
			pkgs = append(pkgs, strings.TrimSpace(v))
		}
	}
	args := []string{
		"diff",
		fmt.Sprintf("%s...%s", config.CRConf.Code.Git.CompareBranch, config.CRConf.Code.Git.ReviewBranch),
		"--name-only",
		"--diff-filter=d",
	}
	args = append(args, pkgs...)
	result, err := ShellExec("git", args...)
	if err != nil {
		return
	}
	fmt.Println("done.")
	arrs := strings.Split(result, "\n")
	// filter files
	for _, v := range arrs {
		if len(config.CRConf.Languages) > 0 && !slices.Contains(config.CRConf.Languages, MatchFileLanguage(v)) {
			continue
		}
		if len(v) > 0 {
			for _, ignore := range config.CRConf.Code.Files.Ignore {
				if regex, ok := regexpM[ignore]; ok {
					// hit ignore files
					if regex.MatchString(v) {
						goto next
					}
				}
			}
			fileNames = append(fileNames, v)
		}
	next:
	}
	return
}

type CodeBlock struct {
	Language string
	FileName string
	Content  string
	Rules    []string
}

type reviewData struct {
	codePatch string
	rules     []string
}

func filesReview(fileNames []string) (result string) {
	// Split code into blocks
	var codeBlockList []CodeBlock
	for _, fileName := range fileNames {
		content, err := ShellExec("git", "diff", fmt.Sprintf("%s...%s", config.CRConf.Code.Git.CompareBranch, config.CRConf.Code.Git.ReviewBranch), fileName)
		if err != nil {
			continue
		}
		language := MatchFileLanguage(fileName)
		var customRules []string
		for _, rule := range config.CRConf.Knowledge.Custom[language] {
			// If custom rule match
			if regex, ok := regexpM[rule.Regexp]; ok && regex.MatchString(content) {
				customRules = append(customRules, rule.Rules...)
			}
		}

		contentToken := EstimateTokens(content)
		promptToken := EstimateTokens(NewPrompt(PromptConfig{
			TreeCustoms: config.CRConf.Knowledge.TreeStandard,
			CodeCustoms: customRules,
		}))
		if contentToken+promptToken > maxToken {
			fmt.Printf("Token overflow, filename:%s，except:%d, actual:%d\n", fileName, maxToken, contentToken+promptToken)
			continue
		}
		codeBlockList = append(codeBlockList, CodeBlock{
			Language: language,
			FileName: fileName,
			Content:  content,
			Rules:    customRules,
		})
	}
	var (
		reviewDataM = make(map[string]*reviewData)
	)
	for _, v := range codeBlockList {
		if _, ok := reviewDataM[v.Language]; !ok {
			reviewDataM[v.Language] = &reviewData{}
		}
		var (
			codePatch string   // Current code block + cached code
			rules     []string // The Rules of Current code block + cached code
		)
		// Merge code blocks
		codePatch = reviewDataM[v.Language].codePatch + v.Content
		rules = RmDuplication(append(reviewDataM[v.Language].rules, v.Rules...))
		if EstimateTokens(codePatch)+EstimateTokens(NewPrompt(PromptConfig{
			TreeCustoms: config.CRConf.Knowledge.TreeStandard,
			CodeCustoms: rules,
		})) > maxToken {
			// If over max token, trigger review
			resultPart, err := Review(reviewDataM[v.Language].codePatch, v.Language, reviewDataM[v.Language].rules)
			if err != nil {
				Error(err.Error())
			}
			if len(resultPart) > 0 && resultPart != "LGTM" {
				if len(resultPart) > 0 {
					result += "\n"
					result += "---\n"
					result += "\n"
					result += resultPart
					result += "\n"
				}
			}
			reviewDataM[v.Language] = &reviewData{
				codePatch: v.Content,
				rules:     v.Rules,
			}
		} else {
			reviewDataM[v.Language].codePatch = codePatch
			reviewDataM[v.Language].rules = rules
		}
	}

	for language, v := range reviewDataM {
		if len(v.codePatch) == 0 {
			continue
		}
		resultPart, err := Review(v.codePatch, language, v.rules)
		if err != nil {
			Error(err.Error())
		}
		if len(resultPart) > 0 && resultPart != "LGTM" {
			if len(resultPart) > 0 {
				result += "\n"
				result += "---\n"
				result += "\n"
				result += resultPart
				result += "\n"
			}
		}
		reviewDataM[language] = &reviewData{}
	}
	return
}
