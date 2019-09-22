package main

import (
	"github.com/sirupsen/logrus"

	"github.com/fidellr/jastip/backend/rover/cmd"
)

func main() {
	if err := cmd.RootCMD.Execute(); err != nil {
		logrus.Fatalf("Fail init content Root CMD with error : %s", err.Error())
	}
}
