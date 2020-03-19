package cmd

import (
	"fmt"

	"github.com/leodido/rn2md/pkg/releasenotes"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var opts = NewOptions()

var program = &cobra.Command{
	Use:   "rn2md",
	Long:  "...",
	Short: "Obtain a changelog markdown from release-note blocks",
	PersistentPreRun: func(c *cobra.Command, args []string) {
		if c.Name() != "help" {
			if errs := opts.Validate(); errs != nil {
				for _, err := range errs {
					logger.WithError(err).Error("error validating  options")
				}
				logger.Fatal("exiting for validation errors")
			}
		}
	},
	Run: func(c *cobra.Command, args []string) {
		client := releasenotes.NewClient()
		notes, err := client.Get(opts.Org, opts.Repo, opts.Branch, opts.Milestone)
		if err != nil {
			logger.WithError(err).Fatal("error retrieving PRs")
		}
		output, err := releasenotes.Print(opts.Milestone, notes)
		if err != nil {
			logger.WithError(err).Fatal("error printing out release notes")
		}
		fmt.Println(output)
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	// Setup flags before the command is initialized
	flags := program.PersistentFlags()
	flags.StringVarP(&opts.Milestone, "milestone", "m", opts.Milestone, "...")
	flags.StringVarP(&opts.Org, "org", "o", opts.Org, "...")
	flags.StringVarP(&opts.Repo, "repo", "r", opts.Repo, "...")
	flags.StringVarP(&opts.Branch, "branch", "b", opts.Branch, "...")
}

func initConfig() {
	// nop
}

// Run ...
func Run() error {
	return program.Execute()
}