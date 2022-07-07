package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "latency-monitor",
	Short: "--------------------------------------------------------------------\n" +
		"LMON\n" +
		"--------------------------------------------------------------------\n" +
		"Description: Latency monitoring tool.\n" +
		"Version: 1.0",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
