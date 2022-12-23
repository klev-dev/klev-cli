package main

import (
	"github.com/spf13/cobra"

	api "github.com/klev-dev/klev-api-go"
)

func logs() *cobra.Command {
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
		if md := *metadata; md != "" {
			out, err := klient.LogsFind(cmd.Context(), md)
			return output(out, err)
		} else {
			out, err := klient.LogsList(cmd.Context())
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

	var in api.LogCreate

	cmd.Flags().StringVar(&in.Metadata, "metadata", "", "machine readable metadata")
	cmd.Flags().BoolVar(&in.Compacting, "compacting", false, "if the log is compacting")
	cmd.Flags().Int64Var(&in.TrimBytes, "trim-bytes", 0, "size of the log to trim")
	cmd.Flags().Int64Var(&in.TrimSeconds, "trim-seconds", 0, "age of the log to trim")
	cmd.Flags().Int64Var(&in.CompactSeconds, "compact-seconds", 0, "age of the log to compact")
	cmd.Flags().Int64Var(&in.ExpireSeconds, "expire-seconds", 0, "age of the log to expire")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		out, err := klient.LogCreate(cmd.Context(), in)
		return output(out, err)
	}

	return cmd
}

func logsGet() *cobra.Command {
	return &cobra.Command{
		Use:   "get log_id",
		Short: "get a log",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := klient.LogGet(cmd.Context(), api.LogID(args[0]))
			return output(out, err)
		},
	}
}

func logsDelete() *cobra.Command {
	return &cobra.Command{
		Use:   "delete log_id",
		Short: "delete a log",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := klient.LogDelete(cmd.Context(), api.LogID(args[0]))
			return output(out, err)
		},
	}
}
