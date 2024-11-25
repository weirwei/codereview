package cmd

import (
	"flag"
	"regexp"

	"github.com/spf13/cobra"
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
		cmd.Help()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

var (
	output  string // output: file name
	pkg     string // pkg: review dir
	model   string // model: llm model
	debug   bool   // debug: debug mode
	version bool   // version: print version

	regexpM  = map[string]*regexp.Regexp{}
	maxToken = 1024 * 16
)

func init() {
	flagParse()

}

func flagParse() {
	flag.StringVar(&output, "o", "", "output filename")
	flag.StringVar(&pkg, "p", "", "review package, split with ','.")
	flag.StringVar(&model, "m", "claude-3-5-sonnet-20240620", "specified model")
	flag.BoolVar(&debug, "d", false, "debug mode")
	flag.BoolVar(&version, "v", false, "version")
	flag.Parse()
}
