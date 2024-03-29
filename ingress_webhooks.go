package main

import (
	"github.com/spf13/cobra"

	"github.com/klev-dev/klev-api-go"
)

func ingressWebhooksRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ingress-webhooks",
		Short: "interact with ingress webhooks",
	}

	cmd.AddCommand(ingressWebhooksList())
	cmd.AddCommand(ingressWebhooksCreate())
	cmd.AddCommand(ingressWebhooksGet())
	cmd.AddCommand(ingressWebhooksUpdate())
	cmd.AddCommand(ingressWebhooksDelete())

	return cmd
}

func ingressWebhooksList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list ingress webhooks",
		Args:  cobra.NoArgs,
	}

	metadata := cmd.Flags().String("metadata", "", "webhook metadata")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed("metadata") {
			out, err := klient.IngressWebhooks.Find(cmd.Context(), *metadata)
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
		Args:  cobra.NoArgs,
	}

	var in klev.IngressWebhookCreateParams

	logID := cmd.Flags().String("log-id", "", "log id that will store webhook data")
	cmd.Flags().StringVar(&in.Metadata, "metadata", "", "machine readable metadata")
	typ := cmd.Flags().String("type", "", "the type of the webhook")
	cmd.Flags().StringVar(&in.Secret, "secret", "", "the secret to validate webhook messages")

	cmd.MarkFlagRequired("log-id")
	cmd.MarkFlagRequired("type")
	cmd.MarkFlagRequired("secret")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		var err error

		in.LogID, err = klev.ParseLogID(*logID)
		if err != nil {
			return outputErr(err)
		}

		in.Type, err = klev.ParseIngressWebhookType(*typ)
		if err != nil {
			return outputErr(err)
		}

		out, err := klient.IngressWebhooks.Create(cmd.Context(), in)
		return output(out, err)
	}

	return cmd
}

func ingressWebhooksGet() *cobra.Command {
	return &cobra.Command{
		Use:   "get <ingress-webhook-id>",
		Short: "get an ingress webhook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := klev.ParseIngressWebhookID(args[0])
			if err != nil {
				return outputErr(err)
			}

			out, err := klient.IngressWebhooks.Get(cmd.Context(), id)
			return output(out, err)
		},
	}
}

func ingressWebhooksUpdate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <ingress-webhook-id>",
		Short: "update ingress webhook",
		Args:  cobra.ExactArgs(1),
	}

	metadata := cmd.Flags().String("metadata", "", "machine readable metadata")
	secret := cmd.Flags().String("secret", "", "the secret to validate webhook messages")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		id, err := klev.ParseIngressWebhookID(args[0])
		if err != nil {
			return outputErr(err)
		}

		var in klev.IngressWebhookUpdateParams

		if cmd.Flags().Changed("metadata") {
			in.Metadata = metadata
		}
		if cmd.Flags().Changed("secret") {
			in.Secret = secret
		}

		out, err := klient.IngressWebhooks.UpdateRaw(cmd.Context(), id, in)
		return output(out, err)
	}

	return cmd
}

func ingressWebhooksDelete() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <ingress-webhook-id>",
		Short: "delete an ingress webhook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := klev.ParseIngressWebhookID(args[0])
			if err != nil {
				return outputErr(err)
			}

			out, err := klient.IngressWebhooks.Delete(cmd.Context(), id)
			return output(out, err)
		},
	}
}
