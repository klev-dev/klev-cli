package main

import (
	"fmt"
	"os"

	klev "github.com/klev-dev/klev-api-go"
	"github.com/spf13/cobra"
)

var klient *klev.Client

var rootCmd = &cobra.Command{
	Use:   "klev",
	Short: "cli to interact with klev",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cfg := klev.NewConfig(os.Getenv("KLEV_TOKEN"))
		klient = klev.New(cfg)
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
