package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Languages []string  `yaml:"languages"`
	Code      Code      `yaml:"code"`
	Knowledge Knowledge `yaml:"knowledge"`
}

type CustomRule struct {
	Regexp string   `yaml:"regexp"`
	Rules  []string `yaml:"rules"`
}
type Code struct {
	Git   Git   `yaml:"git"`
	Files Files `yaml:"files"`
}
type Files struct {
	Ignore []string `yaml:"ignore"`
}

type Git struct {
	ReviewBranch  string `yaml:"review_branch"`
	CompareBranch string `yaml:"compare_branch"`
}

type Knowledge struct {
	Custom       map[string][]CustomRule `yaml:"custom"`
	TreeStandard map[string][]string     `yaml:"tree_standard"`
}

var CRConf Config

func LoadConfig() {
	// read config
	yamlFile, err := os.ReadFile(".codereview.yaml")
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			panic(err)
		}
	} else {
		// unmarshal YAML to config
		err = yaml.Unmarshal(yamlFile, &CRConf)
		if err != nil {
			panic(err)
		}
	}

	if len(CRConf.Code.Git.ReviewBranch) == 0 {
		CRConf.Code.Git.ReviewBranch = "HEAD"
	}
	if len(CRConf.Code.Git.CompareBranch) == 0 {
		CRConf.Code.Git.CompareBranch = "origin/master"
	}
}
