package main

import (
	"github.com/spf13/cobra"

	"github.com/klev-dev/klev-api-go/ingress_webhooks"
	"github.com/klev-dev/klev-api-go/logs"
)

func ingressWebhooksRoot() *cobra.Command {
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
			out, err := klient.IngressWebhooks.Find(cmd.Context(), md)
			return output(out, err)
		} else {
			out, err := klient.IngressWebhooks.List(cmd.Context())
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

	var in ingress_webhooks.CreateParams

	logID := cmd.Flags().String("log-id", "", "log id that will store webhook data")
	cmd.Flags().StringVar(&in.Metadata, "metadata", "", "machine readable metadata")
	cmd.Flags().StringVar(&in.Type, "type", "", "the type of the webhook")
	cmd.Flags().StringVar(&in.Secret, "secret", "", "the secret to validate webhook messages")

	cmd.MarkFlagRequired("log-id")
	cmd.MarkFlagRequired("type")
	cmd.MarkFlagRequired("secret")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		in.LogID = logs.LogID(*logID)

		out, err := klient.IngressWebhooks.Create(cmd.Context(), in)
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
			out, err := klient.IngressWebhooks.Get(cmd.Context(), ingress_webhooks.IngressWebhookID(args[0]))
			return output(out, err)
		},
	}
}

func ingressWebhooksRotate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rotate",
		Short: "rotate ingress webhook secret",
	}

	var in ingress_webhooks.RotateParams

	cmd.Flags().StringVar(&in.Secret, "secret", "", "the secret to validate webhook messages")

	cmd.MarkFlagRequired("secret")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		out, err := klient.IngressWebhooks.RotateRaw(cmd.Context(), ingress_webhooks.IngressWebhookID(args[0]), in)
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
			out, err := klient.IngressWebhooks.Delete(cmd.Context(), ingress_webhooks.IngressWebhookID(args[0]))
			return output(out, err)
		},
	}
}
