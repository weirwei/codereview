package config

import (
	"regexp"
)

const (
	LANGUAGES = "languages"

	CODE_GIT_REVIEW_BRANCH  = "code.git.review_branch"
	CODE_GIT_COMPARE_BRANCH = "code.git.compare_branch"

	CODE_FILES_IGNORE = "code.files.ignore"

	KNOWLEDGE = "knowledge"
)

type CustomRule struct {
	Regexp  string        `yaml:"regexp" json:"regexp"`
	RegexpF regexp.Regexp `yaml:"-" json:"-"`
	Rules   []string      `yaml:"rules" json:"rules"`
}

type Knowledge struct {
	Custom       map[string][]CustomRule `yaml:"custom"`
	TreeStandard map[string][]string     `yaml:"tree_standard"`
}
