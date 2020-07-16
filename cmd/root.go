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
	Long:  "Little configurable CLI to generate the markdown for your changelogs from release-note blocks found into your project pull requests.",
	Short: "Generate markdown for your changelogs from release-note blocks.",
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
		client := releasenotes.NewClient(opts.Token)
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
	flags.StringVarP(&opts.Milestone, "milestone", "m", opts.Milestone, "the milestone you want to filter by the pull requests")
	flags.StringVarP(&opts.Org, "org", "o", opts.Org, "the github organization")
	flags.StringVarP(&opts.Repo, "repo", "r", opts.Repo, "the github repository name")
	flags.StringVarP(&opts.Branch, "branch", "b", opts.Branch, "the target branch you want to filter by the pull requests")
	flags.StringVarP(&opts.Token, "token", "t", opts.Token, "a GitHub personal API token to perform authenticated requests")
}

func initConfig() {
	// nop
}

// Run ...
func Run() error {
	return program.Execute()
}
