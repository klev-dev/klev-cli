package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	api "github.com/klev-dev/klev-api-go"
)

var klient *api.Client

var rootCmd = &cobra.Command{
	Use:   "klev",
	Short: "cli to interact with klev",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cfg := api.NewConfig(os.Getenv("KLEV_TOKEN"))
		klient = api.New(cfg)
	},
}

func main() {
	rootCmd.AddCommand(tokens())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func output(v any, err error) error {
	if err := api.GetError(err); err != nil {
		return outputValue(os.Stderr, err)
	} else if err != nil {
		return err
	}
	return outputValue(os.Stdout, v)
}

func outputValue(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
