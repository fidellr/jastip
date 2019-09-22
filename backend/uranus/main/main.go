package main

import (
	"github.com/fidellr/jastip/backend/uranus/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := cmd.RootCMD.Execute(); err != nil {
		logrus.Fatal("Fail init Root CMD with error")
	}
}
