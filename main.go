package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/klev-dev/klev-api-go"
	"github.com/klev-dev/klev-api-go/clients"
)

var klient *clients.Clients

func main() {
	rootCmd := root()
	rootCmd.AddCommand(paths())
	rootCmd.AddCommand(publish())
	rootCmd.AddCommand(consume())
	rootCmd.AddCommand(receive())
	rootCmd.AddCommand(cleanup())
	rootCmd.AddCommand(logsRoot())
	rootCmd.AddCommand(offsetsRoot())
	rootCmd.AddCommand(tokensRoot())
	rootCmd.AddCommand(ingressWebhooksRoot())
	rootCmd.AddCommand(egressWebhooksRoot())
	rootCmd.AddCommand(filtersRoot())

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

		cfg := klev.NewConfig(auth)
		if cmd.Flags().Changed("base-url") {
			cfg.BaseURL = *base
		} else if base := os.Getenv("KLEV_URL"); base != "" {
			cfg.BaseURL = base
		}
		klient = clients.New(cfg)
		return nil
	}

	return cmd
}

func paths() *cobra.Command {
	return &cobra.Command{
		Use:   "paths",
		Short: "get paths in klev; validate token",
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := klient.Paths.Get(cmd.Context())
			return output(out, err)
		},
	}
}

func output(v any, err error) error {
	if err != nil {
		return outputErr(err)
	}
	return outputValue(os.Stdout, v)
}

func outputErr(err error) error {
	if err := klev.GetError(err); err != nil {
		if err := outputValue(os.Stderr, err); err != nil {
			return err
		}
		os.Exit(1)
	}
	return err
}

func outputValue(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
