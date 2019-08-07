package cmd

import (
	"github.com/spf13/cobra"
)

var RootCMD = &cobra.Command{
	Use:   "rover",
	Short: "Rover is a service that handle content management",
	Long:  "See https://github.com/fidellr/jastip_way/backend/rover for more information",
}
