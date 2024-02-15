package main

import (
	"encoding/json"

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
	cmd.AddCommand(tokensUpdate())
	cmd.AddCommand(tokensDelete())

	return cmd
}

func tokensList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list tokens",
		Args:  cobra.NoArgs,
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
		Args:  cobra.NoArgs,
	}

	var in klev.TokenCreateParams

	cmd.Flags().StringVar(&in.Metadata, "metadata", "", "machine readable metadata")
	acl := cmd.Flags().StringArray("acl", nil, "token acl")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		for _, l := range *acl {
			var item klev.ACLItem
			if err := json.Unmarshal([]byte(l), &item); err != nil {
				return outputErr(err)
			}
			in.ACL = append(in.ACL, item)
		}
		out, bearer, err := klient.Tokens.Create(cmd.Context(), in)
		out.Bearer = bearer
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
			id, err := klev.ParseTokenID(args[0])
			if err != nil {
				return outputErr(err)
			}

			out, err := klient.Tokens.Get(cmd.Context(), id)
			return output(out, err)
		},
	}
}

func tokensUpdate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <token-id>",
		Short: "update token",
		Args:  cobra.ExactArgs(1),
	}

	metadata := cmd.Flags().String("metadata", "", "machine readable metadata")
	acl := cmd.Flags().StringArray("acl", nil, "token acl")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		id, err := klev.ParseTokenID(args[0])
		if err != nil {
			return outputErr(err)
		}

		var in klev.TokenUpdateParams

		if cmd.Flags().Changed("metadata") {
			in.Metadata = metadata
		}
		if cmd.Flags().Changed("acl") {
			var items []klev.ACLItem
			for _, l := range *acl {
				var item klev.ACLItem
				if err := json.Unmarshal([]byte(l), &item); err != nil {
					return outputErr(err)
				}
				items = append(items, item)
			}
			in.ACL = &items
		}

		out, err := klient.Tokens.UpdateRaw(cmd.Context(), id, in)
		return output(out, err)
	}

	return cmd
}

func tokensDelete() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <token-id>",
		Short: "delete a token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := klev.ParseTokenID(args[0])
			if err != nil {
				return outputErr(err)
			}

			out, err := klient.Tokens.Delete(cmd.Context(), id)
			return output(out, err)
		},
	}
}
