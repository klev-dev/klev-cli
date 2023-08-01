package main

import (
	"github.com/spf13/cobra"

	api "github.com/klev-dev/klev-api-go"
)

func egressWebhooks() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "egress-webhooks",
		Short: "interact with egress webhooks",
	}

	cmd.AddCommand(egressWebhooksList())
	cmd.AddCommand(egressWebhooksCreate())
	cmd.AddCommand(egressWebhooksGet())
	cmd.AddCommand(egressWebhooksRotate())
	cmd.AddCommand(egressWebhooksStatus())
	cmd.AddCommand(egressWebhooksDelete())

	return cmd
}

func egressWebhooksList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list egress webhooks",
	}

	metadata := cmd.Flags().String("metadata", "", "webhook metadata")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if md := *metadata; md != "" {
			out, err := klient.EgressWebhooksFind(cmd.Context(), md)
			return output(out, err)
		} else {
			out, err := klient.EgressWebhooksList(cmd.Context())
			return output(out, err)
		}
	}

	return cmd
}

func egressWebhooksCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create new egress webhook",
	}

	var in api.EgressWebhookCreate

	logID := cmd.Flags().String("log-id", "", "log id that will store webhook data")
	cmd.Flags().StringVar(&in.Metadata, "metadata", "", "machine readable metadata")
	cmd.Flags().StringVar(&in.Destination, "destination", "", "where to deliver data")

	cmd.MarkFlagRequired("log-id")
	cmd.MarkFlagRequired("destination")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		in.LogID = api.LogID(*logID)

		out, err := klient.EgressWebhookCreate(cmd.Context(), in)
		return output(out, err)
	}

	return cmd
}

func egressWebhooksGet() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "get an egress webhook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := klient.EgressWebhookGet(cmd.Context(), api.EgressWebhookID(args[0]))
			return output(out, err)
		},
	}
}

func egressWebhooksRotate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rotate",
		Short: "rotate egress webhook secret",
		Args:  cobra.ExactArgs(1),
	}

	var in api.EgressWebhookRotate

	cmd.Flags().Int64Var(&in.ExpireSeconds, "expire-seconds", 0, "for how long the old secret is valid")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		out, err := klient.EgressWebhookRotateRaw(cmd.Context(), api.EgressWebhookID(args[0]), in)
		return output(out, err)
	}

	return cmd
}

func egressWebhooksStatus() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "status an egress webhook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := klient.EgressWebhookStatus(cmd.Context(), api.EgressWebhookID(args[0]))
			return output(out, err)
		},
	}
}

func egressWebhooksDelete() *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "delete an egress webhook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := klient.EgressWebhookDelete(cmd.Context(), api.EgressWebhookID(args[0]))
			return output(out, err)
		},
	}
}
