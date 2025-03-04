package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/weirwei/codereview/code"
	"github.com/weirwei/codereview/config"
	"github.com/weirwei/codereview/log"
	"github.com/weirwei/codereview/review"
)

var rootCmd = &cobra.Command{
	Use:   "codereview",
	Short: "CodeReview is a command line tool for code review",
	Long: `CodeReview is a command line tool that helps developers
perform code reviews more efficiently and effectively.`,
	Run: func(cmd *cobra.Command, args []string) {
		if version {
			versionCmd.Run(cmd, args)
			return
		}

		exec()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

var (
	// output  string // output: file name
	pkg     string // pkg: review dir
	version bool   // version: print version
	debug   bool   // debug: set log level to DEBUG

	userViper = viper.New()
	projViper = viper.New()
)

func flagParse() {
	// rootCmd.PersistentFlags().StringP("output", "o", "", "output filename")
	rootCmd.PersistentFlags().StringP("pkg", "p", "", "review package, split with ','.")
	rootCmd.PersistentFlags().BoolVarP(&version, "version", "v", false, "version")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "set log level to DEBUG")
	// viper.BindPFlag("llm.model", rootCmd.PersistentFlags().Lookup("model"))
}

func init() {
	flagParse()
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	// Search config in home directory with name ".cobra" (without extension).
	userViper.AddConfigPath(home)
	userViper.SetConfigName(".codereview")
	userViper.SetConfigType("ini")
	if err := userViper.ReadInConfig(); err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			// 文件没找到
			createConfig()
			err = userViper.ReadInConfig()
		}
		if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			cobra.CheckErr(err)
		}
	}
	setLogLevel()
	projViper.SetConfigName(".codereview")
	projViper.SetConfigType("yaml")
	projViper.AddConfigPath("./")
	if err := projViper.ReadInConfig(); err != nil {
		if !errors.As(err, &viper.ConfigFileNotFoundError{}) && !errors.As(err, &viper.ConfigParseError{}) {
			cobra.CheckErr(err)
		}
	}
}

func setLogLevel() {
	if debug {
		log.SetLevel(log.LevelDebug)
		return
	}
	switch strings.ToUpper(userViper.GetString("log.level")) {
	case "DEBUG":
		log.SetLevel(log.LevelDebug)
	case "INFO":
		log.SetLevel(log.LevelInfo)
	case "WARN":
		log.SetLevel(log.LevelWarn)
	case "ERROR":
		log.SetLevel(log.LevelError)
	case "FATAL":
		log.SetLevel(log.LevelFatal)
	case "PANIC":
		log.SetLevel(log.LevelPanic)
	default:
		log.SetLevel(log.LevelInfo)
	}
}

func setRegexpF(knowledge *config.Knowledge) {
	for _, rules := range knowledge.Custom {
		for i, rule := range rules {
			rules[i].RegexpF = regexp.MustCompile(rule.Regexp)
		}
	}
}

func createConfig() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please input your llm base url(eg: https://api.anthropic.com/v1): ")
	baseUrl, _ := reader.ReadString('\n')
	fmt.Print("Please input your API sk: ")
	sk, _ := reader.ReadString('\n')
	fmt.Print("Please input your llm model name: ")
	model, _ := reader.ReadString('\n')

	fmt.Print("Please input max token(default 8192): ")
enterMaxToken:
	maxToken, _ := reader.ReadString('\n')
	if len(maxToken) == 0 {
		maxToken = "8192"
	}
	maxTokenInt, err := strconv.ParseInt(maxToken, 10, 64)
	if err != nil {
		fmt.Print("Please input a valid integer for max token: ")
		goto enterMaxToken
	}
	defaultConfig := fmt.Sprintf(`[llm]
base_url=%s
sk=%s
model=%s
max_token=%d

[log]
level=INFO`, baseUrl, sk, model, maxTokenInt)
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	if err := os.WriteFile(home+"/.codereview", []byte(defaultConfig), 0644); err != nil {
		log.Errorf("create config file failed: %v", err)
		return
	}
	fmt.Print("Config file created at: " + home + "/.codereview")
	fmt.Print("You can change your options in config file")
}

func exec() {
	ctx := context.Background()
	var (
		baseUrl  = userViper.GetString("llm.base_url")
		sk       = userViper.GetString("llm.sk")
		model    = userViper.GetString("llm.model")
		maxToken = userViper.GetInt("llm.max_token")

		language      = projViper.GetString("language") // todo: different language custom prompt
		reviewBranch  = projViper.GetString(config.CODE_GIT_REVIEW_BRANCH)
		compareBranch = projViper.GetString(config.CODE_GIT_COMPARE_BRANCH)
		ignore        = projViper.GetStringSlice(config.CODE_FILES_IGNORE)

		knowledge config.Knowledge
	)
	data, _ := projViper.Get(config.KNOWLEDGE).(map[string]any)
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &knowledge,
		TagName: "json",
	})
	if err != nil {
		cobra.CheckErr(err)
	}
	if err := decoder.Decode(data); err != nil {
		cobra.CheckErr(err)
	}
	// set regexp from regexp
	setRegexpF(&knowledge)
	var filepathFilters []regexp.Regexp
	for _, v := range ignore {
		filepathFilters = append(filepathFilters, *regexp.MustCompile(v))
	}
	var pkgs []string
	if len(pkg) > 0 {
		pkgs = strings.Split(pkg, ",")
	}
	codePatch, err := code.NewGit(code.GitCond{
		Pkgs:            pkgs,
		ReviewBranch:    reviewBranch,
		CompareBranch:   compareBranch,
		MaxToken:        maxToken,
		Knowledge:       knowledge,
		FilepathFilters: filepathFilters,
	}).GetCode()
	if err != nil {
		log.Errorf("get code failed, %s", err)
		return
	}
	reviewer := review.NewDefaultReviewer(ctx, baseUrl, sk, model, maxToken, language)
	reviewer.SetHandler(review.GetDefaultHandler(func(data string) {
		fmt.Print(data)
	}))
	for _, v := range codePatch {
		reviewer.SetCodePatch(v)
		if err := reviewer.Exec(); err != nil {
			log.Errorf("exec failed, files: %s\nerr: %s", strings.Join(v.Filepaths, "\n"), err)
		}
	}
}
