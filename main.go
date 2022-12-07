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
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if token := os.Getenv("KLEV_TOKEN"); token != "" {
			cfg := api.NewConfig(token)
			klient = api.New(cfg)
			return nil
		}
		return fmt.Errorf("KLEV_TOKEN is required, get it from https://dash.klev.dev")
	},
}

func main() {
	rootCmd.AddCommand(paths())
	rootCmd.AddCommand(publish())
	rootCmd.AddCommand(consume())
	rootCmd.AddCommand(logs())
	rootCmd.AddCommand(tokens())
	rootCmd.AddCommand(webhooks())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func paths() *cobra.Command {
	return &cobra.Command{
		Use:   "paths",
		Short: "get paths in klev; validate token",
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := klient.Paths(cmd.Context())
			return output(out, err)
		},
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
