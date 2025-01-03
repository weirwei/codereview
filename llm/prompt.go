package llm

import (
	"fmt"
	"strings"
)

const (
	printTemplateKey = "${TEMPLATE}"
	treeStandardKey  = "${TREE_STANDARD}"
	codeStandardKey  = "${CODE_STANDARD}"
)

const printTemplate = "## 目录树分析\n> {{文件路径}}:{{line num}}\n```{{language}}\n{{code block}}\n```\n\n- 违反规则: {{violation of rules}}\n\n- 违反规则: {{violation of rules}}\n\n## 代码分析\n\n> {{文件路径}}:{{line num}}\n```{{language}}\n{{code block}}\n```\n\n- 隐患: {{dangers}}\n 评级: {{rank}}\n 建议: {{suggestions}}\n\n- 隐患: {{dangers}}\n 评级: {{rank}}\n 建议: {{suggestions}}\n"

var prompt = `
<constraints>

	你正在进行代码审查。这个代码是git diff命令的输出。已删除的行以减号(-)开头，已添加的行以加号(+)开头。

	其他行是为了提供上下文，但在审查中不应该被忽略。

	user 输入的内容就是你需要审查的代码。

	「steps」 就是你审查代码的工作流程，你在执行的时候需要严格遵守 「rules」 中的条款。

	在 「<process></process>」 标签中输出执行或者思考步骤。

	在 「<output></output>」 标签将你的审查报告按照 markdown 格式输出，输出模板参考 「template」，输出规范参考。

  其他约束：
    - 在 目录树分析 或者 代码分析 的时候，如果出现代码符合规范或者没有违反任何规则的场景下，无需输出任何内容
    - 目录树分析 只需要关注代码是否符合 「目录树约束」 ，不需要做额外的检查
    - 如果所有内容都符合规范，只需输出 LGTM 即可

</constraints>

<rules>
	<目录树约束>
${TREE_STANDARD}
	</目录树约束>

	<代码规范>
${CODE_STANDARD}
	</代码规范>

	<风险评级规则>
		P0：致命的问题。导致系统崩溃、数据丢失、安全漏洞或严重功能失效的问题 必须立即修复
		P1：严重的问题。影响系统功能，但不会导致系统崩溃 需要优先修复
		P2：一般的问题。轻微影响系统功能或用户体验，或者只是代码风格或可读性问题 可以稍后修复
		P3：建议修改的问题。不影响系统功能或用户体验，只是建议性的改进 选择性修复
	</风险评级规则>

  <输出规范>
      1. 输出的格式是 markdown
      2. 只需要输出不符合规范和约束的内容
      3. 输出语言为 中文
      4. 输出内容尽量简短、扼要
  </输出规范>
    
</rules>

<steps>
	1. 代码拆分。按文件维度拆分代码，你会得到 代码块 ，代码块包含 文件名和代码内容
	2. 目录树分析。对代码块逐个进行目录树分析，目录树分析满足 「rules」 中 「目录树约束」，如果满足，则无需输出任何内容
	3. 代码分析。对代码块逐个进行代码分析，分析代码的时候，将 所有用户输入的代码 作为上下文辅助审查，代码分析满足 「rules」 中的 「代码规范」，如果满足，则无需输出任何内容
	4. 风险评级。分析代码隐患，并根据 「rules」 中 「风险评级规则」 评级，得到风险评级
	5. 修改建议。给出代码隐患修改建议
</steps>

<template>
${TEMPLATE}
</template>

`
var treeStandard = map[string][]string{
	"/cmd": {
		"目录下的每个子目录都应该对应一个独立的可执行程序",
		"子目录名应该与生成的程序名一致",
		"避免放置大量业务逻辑代码",
	},
}

var codeStandard = []string{
	"go语言通用规范，以Google的Go语言编码规范为标准",
}

type PromptConfig struct {
	TreeCustoms map[string][]string
	CodeCustoms []string
}

func NewPrompt(config PromptConfig) string {
	var (
		treeData = make(map[string][]string)
		codeData = make([]string, 0)

		formatTreeStr, formatCodeStr string
	)
	for k, v := range treeStandard {
		treeData[k] = v
	}
	for k, v := range config.TreeCustoms {
		treeData[k] = v
	}
	codeData = append(codeData, codeStandard...)
	codeData = append(codeData, config.CodeCustoms...)

	for k, v := range treeData {
		formatTreeStr += fmt.Sprintf("\t%s:\n\t\t- %s\n", k, strings.Join(v, "\n\t\t- "))
	}
	for _, v := range codeData {
		formatCodeStr += fmt.Sprintf("\t\"%s\"\n", v)
	}
	finalPrompt := prompt

	printTemplateCompact := printTemplate

	finalPrompt = strings.ReplaceAll(finalPrompt, printTemplateKey, printTemplateCompact)
	finalPrompt = strings.ReplaceAll(finalPrompt, treeStandardKey, formatTreeStr)
	finalPrompt = strings.ReplaceAll(finalPrompt, codeStandardKey, formatCodeStr)
	return finalPrompt
}
