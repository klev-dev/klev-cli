package main

import (
	"github.com/spf13/cobra"

	api "github.com/klev-dev/klev-api-go"
)

func webhooks() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "webhooks",
		Short: "interact with webhooks",
	}

	cmd.AddCommand(webhooksList())
	cmd.AddCommand(webhooksCreate())
	cmd.AddCommand(webhooksGet())
	cmd.AddCommand(webhooksDelete())

	return cmd
}

func webhooksList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list webhooks",
	}

	metadata := cmd.Flags().String("metadata", "", "webhook metadata")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if md := *metadata; md != "" {
			out, err := klient.WebhooksFind(cmd.Context(), md)
			return output(out, err)
		} else {
			out, err := klient.WebhooksList(cmd.Context())
			return output(out, err)
		}
	}

	return cmd
}

func webhooksCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create new webhook",
	}

	var in api.WebhookIn

	logID := cmd.Flags().String("log-id", "", "log id that will store webhook data")
	cmd.Flags().StringVar(&in.Metadata, "metadata", "", "metadata of the webhook")
	cmd.Flags().StringVar(&in.Type, "type", "", "the type of the webhook")
	cmd.Flags().StringVar(&in.Secret, "secret", "", "the secret to validate webhook messages")

	cmd.MarkFlagRequired("log-id")
	cmd.MarkFlagRequired("type")
	cmd.MarkFlagRequired("secret")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		in.LogID = api.LogID(*logID)

		out, err := klient.WebhookCreate(cmd.Context(), in)
		return output(out, err)
	}

	return cmd
}

func webhooksGet() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "get a webhook",
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := klient.WebhookGet(cmd.Context(), api.WebhookID(args[0]))
			return output(out, err)
		},
	}
}

func webhooksDelete() *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "delete a webhook",
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := klient.WebhookDelete(cmd.Context(), api.WebhookID(args[0]))
			return output(out, err)
		},
	}
}
