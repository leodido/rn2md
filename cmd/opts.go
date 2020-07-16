package cmd

import (
	"fmt"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/leodido/rn2md/validate"
	logger "github.com/sirupsen/logrus"
)

// Options represents the program options.
type Options struct {
	Milestone string `validate:"required,semver" name:"milestone"`
	Branch    string `default:"master" validate:"omitempty,ascii" name:"branch"`
	Org       string `validate:"required,ascii" name:"organization"`
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
