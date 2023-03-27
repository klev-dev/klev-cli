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

func main() {
	rootCmd := root()
	rootCmd.AddCommand(paths())
	rootCmd.AddCommand(publish())
	rootCmd.AddCommand(consume())
	rootCmd.AddCommand(logs())
	rootCmd.AddCommand(offsets())
	rootCmd.AddCommand(tokens())
	rootCmd.AddCommand(webhooks())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func root() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "klev",
		Short: "cli to interact with klev",
	}

	authtoken := cmd.PersistentFlags().String("authtoken", "", "token to use for authorization")
	base := cmd.PersistentFlags().String("base-url", "", "base url to talk to")
	cmd.PersistentFlags().MarkHidden("base-url")

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		var auth string
		if token := *authtoken; token != "" {
			auth = token
		} else if token := os.Getenv("KLEV_TOKEN"); token != "" {
			auth = token
		} else {
			return fmt.Errorf("authtoken is missing. pass with with '--authtoken' or via KLEV_TOKEN env variable. get it from https://dash.klev.dev")
		}

		cfg := api.NewConfig(auth)
		if cmd.Flags().Changed("base-url") {
			cfg.BaseURL = *base
		} else if base := os.Getenv("KLEV_URL"); base != "" {
			cfg.BaseURL = base
		}
		klient = api.New(cfg)
		return nil
	}

	return cmd
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
