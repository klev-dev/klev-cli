package main

import (
	"github.com/klev-dev/klev-api-go"
	"github.com/spf13/cobra"
)

func logsRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "interact with logs",
	}

	cmd.AddCommand(logsList())
	cmd.AddCommand(logsCreate())
	cmd.AddCommand(logsGet())
	cmd.AddCommand(logsDelete())

	return cmd
}

func logsList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list logs",
	}

	metadata := cmd.Flags().String("metadata", "", "log metadata")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed("metadata") {
			out, err := klient.Logs.Find(cmd.Context(), *metadata)
			return output(out, err)
		} else {
			out, err := klient.Logs.List(cmd.Context())
			return output(out, err)
		}
	}

	return cmd
}

func logsCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create new log",
	}

	var in klev.LogCreateParams

	cmd.Flags().StringVar(&in.Metadata, "metadata", "", "machine readable metadata")
	cmd.Flags().BoolVar(&in.Compacting, "compacting", false, "if the log is compacting")
	cmd.Flags().Int64Var(&in.TrimBytes, "trim-bytes", 0, "size of the log to trim")
	cmd.Flags().Int64Var(&in.TrimSeconds, "trim-seconds", 0, "age of the log to trim")
	cmd.Flags().Int64Var(&in.CompactSeconds, "compact-seconds", 0, "age of the log to compact")
	cmd.Flags().Int64Var(&in.ExpireSeconds, "expire-seconds", 0, "age of the log to expire")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		out, err := klient.Logs.Create(cmd.Context(), in)
		return output(out, err)
	}

	return cmd
}

func logsGet() *cobra.Command {
	return &cobra.Command{
		Use:   "get <log-id>",
		Short: "get a log",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := klev.LogID(args[0])
			out, err := klient.Logs.Get(cmd.Context(), id)
			return output(out, err)
		},
	}
}

func logsDelete() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <log-id>",
		Short: "delete a log",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := klev.LogID(args[0])
			out, err := klient.Logs.Delete(cmd.Context(), id)
			return output(out, err)
		},
	}
}
