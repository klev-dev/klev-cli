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
	cmd.AddCommand(logsStats())
	cmd.AddCommand(logsUpdate())
	cmd.AddCommand(logsDelete())

	return cmd
}

func logsList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list logs",
		Args:  cobra.NoArgs,
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
		Args:  cobra.NoArgs,
	}

	var in klev.LogCreateParams

	cmd.Flags().StringVar(&in.Metadata, "metadata", "", "machine readable metadata")
	cmd.Flags().BoolVar(&in.Compacting, "compacting", false, "if the log is compacting")
	cmd.Flags().Int64Var(&in.TrimSeconds, "trim-seconds", 0, "age of the log to trim")
	cmd.Flags().Int64Var(&in.TrimSize, "trim-size", 0, "size of the log to trim")
	cmd.Flags().Int64Var(&in.TrimCount, "trim-count", 0, "count of message in the log to trim")
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
			id, err := klev.ParseLogID(args[0])
			if err != nil {
				return err
			}

			out, err := klient.Logs.Get(cmd.Context(), id)
			return output(out, err)
		},
	}
}

func logsStats() *cobra.Command {
	return &cobra.Command{
		Use:   "stats <log-id>",
		Short: "stats a log",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := klev.ParseLogID(args[0])
			if err != nil {
				return err
			}

			out, err := klient.Logs.Stats(cmd.Context(), id)
			return output(out, err)
		},
	}
}

func logsUpdate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <log-id>",
		Short: "update log",
		Args:  cobra.ExactArgs(1),
	}

	metadata := cmd.Flags().String("metadata", "", "machine readable metadata")
	trimSeconds := cmd.Flags().Int64("trim-seconds", 0, "age of the log to trim")
	trimSize := cmd.Flags().Int64("trim-size", 0, "size of the log to trim")
	trimCount := cmd.Flags().Int64("trim-count", 0, "count of message in the log to trim")
	compactSeconds := cmd.Flags().Int64("compact-seconds", 0, "age of the log to compact")
	expireSeconds := cmd.Flags().Int64("expire-seconds", 0, "age of the log to expire")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		id, err := klev.ParseLogID(args[0])
		if err != nil {
			return err
		}

		var in klev.LogUpdateParams

		if cmd.Flags().Changed("metadata") {
			in.Metadata = metadata
		}
		if cmd.Flags().Changed("trim-seconds") {
			in.TrimSeconds = trimSeconds
		}
		if cmd.Flags().Changed("trim-size") {
			in.TrimSize = trimSize
		}
		if cmd.Flags().Changed("trim-count") {
			in.TrimCount = trimCount
		}
		if cmd.Flags().Changed("compact-seconds") {
			in.CompactSeconds = compactSeconds
		}
		if cmd.Flags().Changed("expire-seconds") {
			in.ExpireSeconds = expireSeconds
		}

		out, err := klient.Logs.UpdateRaw(cmd.Context(), id, in)
		return output(out, err)
	}

	return cmd
}

func logsDelete() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <log-id>",
		Short: "delete a log",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := klev.ParseLogID(args[0])
			if err != nil {
				return err
			}

			out, err := klient.Logs.Delete(cmd.Context(), id)
			return output(out, err)
		},
	}
}
