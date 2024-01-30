package main

import (
	"github.com/klev-dev/klev-api-go"
	"github.com/spf13/cobra"
)

func tokensRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tokens",
		Short: "interact with tokens",
	}

	cmd.AddCommand(tokensList())
	cmd.AddCommand(tokensCreate())
	cmd.AddCommand(tokensGet())
	cmd.AddCommand(tokensDelete())

	return cmd
}

func tokensList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list tokens",
	}

	metadata := cmd.Flags().String("metadata", "", "token metadata")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed("metadata") {
			out, err := klient.Tokens.Find(cmd.Context(), *metadata)
			return output(out, err)
		} else {
			out, err := klient.Tokens.List(cmd.Context())
			return output(out, err)
		}
	}

	return cmd
}

func tokensCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create new token",
	}

	var in klev.TokenCreateParams

	cmd.Flags().StringVar(&in.Metadata, "metadata", "", "machine readable metadata")
	cmd.Flags().StringArrayVar(&in.ACL, "acl", nil, "token acl")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		out, err := klient.Tokens.Create(cmd.Context(), in)
		return output(out, err)
	}

	return cmd
}

func tokensGet() *cobra.Command {
	return &cobra.Command{
		Use:   "get <token-id>",
		Short: "get a token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := klev.TokenID(args[0])
			out, err := klient.Tokens.Get(cmd.Context(), id)
			return output(out, err)
		},
	}
}

func tokensDelete() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <token-id>",
		Short: "delete a token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := klev.TokenID(args[0])
			out, err := klient.Tokens.Delete(cmd.Context(), id)
			return output(out, err)
		},
	}
}
