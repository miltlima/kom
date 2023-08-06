package main

import (
	"fmt"
	"kom/cmd"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "kom"}
	rootCmd.AddCommand(cmd.PodsCmd, cmd.NodesCmd, cmd.LogsCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
