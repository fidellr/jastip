package cmd

import "github.com/spf13/cobra"

var RootCMD = &cobra.Command{
	Use:   "uranus",
	Short: "Uranus is a service that handle users data",
	Long:  "See https://github.com/fidellr/jastip_way/backend/uranus for more information",
}
