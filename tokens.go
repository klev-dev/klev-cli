package main

import (
	"github.com/spf13/cobra"

	api "github.com/klev-dev/klev-api-go"
)

func tokens() *cobra.Command {
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
		if md := *metadata; md != "" {
			out, err := klient.TokensFind(cmd.Context(), md)
			return output(out, err)
		} else {
			out, err := klient.TokensList(cmd.Context())
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

	var in api.TokenIn

	cmd.Flags().StringVar(&in.Metadata, "metadata", "", "machine readable metadata")
	cmd.Flags().StringArrayVar(&in.ACL, "acl", nil, "token acl")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		out, err := klient.TokenCreate(cmd.Context(), in)
		return output(out, err)
	}

	return cmd
}

func tokensGet() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "get a token",
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := klient.TokenGet(cmd.Context(), api.TokenID(args[0]))
			return output(out, err)
		},
	}
}

func tokensDelete() *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "delete a token",
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := klient.TokenDelete(cmd.Context(), api.TokenID(args[0]))
			return output(out, err)
		},
	}
}
