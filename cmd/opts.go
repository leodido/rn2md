package cmd

import (
	"fmt"
	"strings"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/leodido/rn2md/validate"
	logger "github.com/sirupsen/logrus"
)

// Options represents the program options.
type Options struct {
	Milestone string `validate:"required,semver" name:"milestone"`
	Branch    string `default:"master" validate:"omitempty,ascii" name:"branch"`
	Repo      string `validate:"required,ascii" name:"repository"`
	Token     string `validate:"omitempty,len=40" name:"token"`
}

// NewOptions create a pointer to an Options instance.
func NewOptions() *Options {
	o := &Options{}
	if err := defaults.Set(o); err != nil {
		logger.WithError(err).Fatal("error setting options")
	}
	return o
}

// Validate ensures Options are valid, otherwise returns all the occurred errors.
func (o *Options) Validate() []error {
	if err := validate.V.Struct(o); err != nil {
		errors := err.(validator.ValidationErrors)
		errArr := []error{}
		for _, e := range errors {
			// Translate each error one at a time
			errArr = append(errArr, fmt.Errorf(e.Translate(validate.T)))
		}
		return errArr
	}
	return nil
}

func (o *Options) SplitRepoOrgName() (string, string, error) {
	names := strings.Split(o.Repo, "/")
	if len(names) != 2 {
		return "", "", fmt.Errorf("provided repo has wrong format, expected org/repo, actual: %s", o.Repo)
	}
	return names[0], names[1], nil
}
