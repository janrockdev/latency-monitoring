package cmd

import (
	"github.com/spf13/cobra"
	"lmon/component"
)

var (
	cmdMinio = &cobra.Command{
		Use:   "minio",
		Short: "Run benchmark for Min.IO",
		RunE:  runMinio,
	}
)

func init() {
	rootCmd.AddCommand(cmdMinio)
}

func runMinio(*cobra.Command, []string) error {
	return component.RunMinio()
}
