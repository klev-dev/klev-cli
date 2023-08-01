package main

import (
	"github.com/spf13/cobra"

	api "github.com/klev-dev/klev-api-go"
)

func ingressWebhooks() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ingress-webhooks",
		Short: "interact with ingress webhooks",
	}

	cmd.AddCommand(ingressWebhooksList())
	cmd.AddCommand(ingressWebhooksCreate())
	cmd.AddCommand(ingressWebhooksGet())
	cmd.AddCommand(ingressWebhooksRotate())
	cmd.AddCommand(ingressWebhooksDelete())

	return cmd
}

func ingressWebhooksList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list ingress webhooks",
	}

	metadata := cmd.Flags().String("metadata", "", "webhook metadata")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if md := *metadata; md != "" {
			out, err := klient.IngressWebhooksFind(cmd.Context(), md)
			return output(out, err)
		} else {
			out, err := klient.IngressWebhooksList(cmd.Context())
			return output(out, err)
		}
	}

	return cmd
}

func ingressWebhooksCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create new ingress webhook",
	}

	var in api.IngressWebhookCreate

	logID := cmd.Flags().String("log-id", "", "log id that will store webhook data")
	cmd.Flags().StringVar(&in.Metadata, "metadata", "", "machine readable metadata")
	cmd.Flags().StringVar(&in.Type, "type", "", "the type of the webhook")
	cmd.Flags().StringVar(&in.Secret, "secret", "", "the secret to validate webhook messages")

	cmd.MarkFlagRequired("log-id")
	cmd.MarkFlagRequired("type")
	cmd.MarkFlagRequired("secret")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		in.LogID = api.LogID(*logID)

		out, err := klient.IngressWebhookCreate(cmd.Context(), in)
		return output(out, err)
	}

	return cmd
}

func ingressWebhooksGet() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "get an ingress webhook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := klient.IngressWebhookGet(cmd.Context(), api.IngressWebhookID(args[0]))
			return output(out, err)
		},
	}
}

func ingressWebhooksRotate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rotate",
		Short: "rotate ingress webhook secret",
	}

	var in api.IngressWebhookRotate

	cmd.Flags().StringVar(&in.Secret, "secret", "", "the secret to validate webhook messages")

	cmd.MarkFlagRequired("secret")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		out, err := klient.IngressWebhookRotateRaw(cmd.Context(), api.IngressWebhookID(args[0]), in)
		return output(out, err)
	}

	return cmd
}

func ingressWebhooksDelete() *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "delete an ingress webhook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := klient.IngressWebhookDelete(cmd.Context(), api.IngressWebhookID(args[0]))
			return output(out, err)
		},
	}
}
