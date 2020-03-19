package main

import (
	"github.com/leodido/rn2md/cmd"
	logger "github.com/sirupsen/logrus"
)

func main() {
	if err := cmd.Run(); err != nil {
		logger.WithError(err).Fatal("exiting")
	}
}
