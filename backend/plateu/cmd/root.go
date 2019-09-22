package cmd

import (
	"github.com/spf13/cobra"
)

var RootCMD = &cobra.Command{
	Use:   "plateu",
	Short: "Plateu is a service that handle picture compressing or thumbor",
	Long:  "See https://github.com/fidellr/jastip/backend/plateu for more information",
}
