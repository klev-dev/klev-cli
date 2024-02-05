package main

import (
	"github.com/spf13/cobra"

	"github.com/klev-dev/klev-api-go"
)

func egressWebhooksRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "egress-webhooks",
		Short: "interact with egress webhooks",
	}

	cmd.AddCommand(egressWebhooksList())
	cmd.AddCommand(egressWebhooksCreate())
	cmd.AddCommand(egressWebhooksGet())
	cmd.AddCommand(egressWebhooksRotate())
	cmd.AddCommand(egressWebhooksStatus())
	cmd.AddCommand(egressWebhooksUpdate())
	cmd.AddCommand(egressWebhooksDelete())

	return cmd
}

func egressWebhooksList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list egress webhooks",
		Args:  cobra.NoArgs,
	}

	metadata := cmd.Flags().String("metadata", "", "webhook metadata")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed("metadata") {
			out, err := klient.EgressWebhooks.Find(cmd.Context(), *metadata)
			return output(out, err)
		} else {
			out, err := klient.EgressWebhooks.List(cmd.Context())
			return output(out, err)
		}
	}

	return cmd
}

func egressWebhooksCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create new egress webhook",
		Args:  cobra.NoArgs,
	}

	var in klev.EgressWebhookCreateParams

	logID := cmd.Flags().String("log-id", "", "log id that will store webhook data")
	cmd.Flags().StringVar(&in.Metadata, "metadata", "", "machine readable metadata")
	cmd.Flags().StringVar(&in.Destination, "destination", "", "where to deliver data")
	cmd.Flags().StringVar(&in.Payload, "payload", "message", "what payload to deliver")

	cmd.MarkFlagRequired("log-id")
	cmd.MarkFlagRequired("destination")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		in.LogID = klev.LogID(*logID)

		out, err := klient.EgressWebhooks.Create(cmd.Context(), in)
		return output(out, err)
	}

	return cmd
}

func egressWebhooksGet() *cobra.Command {
	return &cobra.Command{
		Use:   "get <egress-webhook-id>",
		Short: "get an egress webhook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := klev.ParseEgressWebhookID(args[0])
			if err != nil {
				return err
			}

			out, err := klient.EgressWebhooks.Get(cmd.Context(), id)
			return output(out, err)
		},
	}
}

func egressWebhooksRotate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rotate <egress-webhook-id>",
		Short: "rotate egress webhook secret",
		Args:  cobra.ExactArgs(1),
	}

	var in klev.EgressWebhookRotateParams

	cmd.Flags().Int64Var(&in.ExpireSeconds, "expire-seconds", 0, "for how long the old secret is valid")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		id, err := klev.ParseEgressWebhookID(args[0])
		if err != nil {
			return err
		}

		out, err := klient.EgressWebhooks.RotateRaw(cmd.Context(), id, in)
		return output(out, err)
	}

	return cmd
}

func egressWebhooksStatus() *cobra.Command {
	return &cobra.Command{
		Use:   "status <egress-webhook-id>",
		Short: "status an egress webhook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := klev.ParseEgressWebhookID(args[0])
			if err != nil {
				return err
			}

			out, err := klient.EgressWebhooks.Status(cmd.Context(), id)
			return output(out, err)
		},
	}
}

func egressWebhooksUpdate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <egress-webhook-id>",
		Short: "update egress webhook",
		Args:  cobra.ExactArgs(1),
	}

	metadata := cmd.Flags().String("metadata", "", "machine readable metadata")
	destination := cmd.Flags().String("destination", "", "where to deliver data")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		id, err := klev.ParseEgressWebhookID(args[0])
		if err != nil {
			return err
		}

		var in klev.EgressWebhookUpdateParams

		if cmd.Flags().Changed("metadata") {
			in.Metadata = metadata
		}
		if cmd.Flags().Changed("destination") {
			in.Destination = destination
		}

		out, err := klient.EgressWebhooks.UpdateRaw(cmd.Context(), id, in)
		return output(out, err)
	}

	return cmd
}

func egressWebhooksDelete() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <egress-webhook-id>",
		Short: "delete an egress webhook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := klev.ParseEgressWebhookID(args[0])
			if err != nil {
				return err
			}

			out, err := klient.EgressWebhooks.Delete(cmd.Context(), id)
			return output(out, err)
		},
	}
}
