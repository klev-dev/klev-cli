package main

import (
	"github.com/klev-dev/klev-api-go"
	"github.com/spf13/cobra"
)

func offsetsRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "offsets",
		Short: "interact with offsets",
	}

	cmd.AddCommand(offsetsList())
	cmd.AddCommand(offsetsCreate())
	cmd.AddCommand(offsetsGet())
	cmd.AddCommand(offsetsSet())
	cmd.AddCommand(offsetsDelete())

	return cmd
}

func offsetsList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list offsets",
	}

	metadata := cmd.Flags().String("metadata", "", "offset metadata")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed("metadata") {
			out, err := klient.Offsets.Find(cmd.Context(), *metadata)
			return output(out, err)
		} else {
			out, err := klient.Offsets.List(cmd.Context())
			return output(out, err)
		}
	}

	return cmd
}

func offsetsCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create new offset",
	}

	var in klev.OffsetCreateParams

	logID := cmd.Flags().String("log-id", "", "log id for this offset")
	cmd.Flags().StringVar(&in.Metadata, "metadata", "", "machine readable metadata")

	cmd.MarkFlagRequired("log-id")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		in.LogID = klev.LogID(*logID)
		out, err := klient.Offsets.Create(cmd.Context(), in)
		return output(out, err)
	}

	return cmd
}

func offsetsGet() *cobra.Command {
	return &cobra.Command{
		Use:   "get <offset-id>",
		Short: "get an offset",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := klev.OffsetID(args[0])
			out, err := klient.Offsets.Get(cmd.Context(), id)
			return output(out, err)
		},
	}
}

func offsetsSet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set <offset-id>",
		Short: "set log offset",
		Args:  cobra.ExactArgs(1),
	}

	var in klev.OffsetSetParams
	cmd.Flags().Int64Var(&in.Value, "value", 0, "value to set")
	cmd.Flags().StringVar(&in.ValueMetadata, "value-metadata", "", "machine readable metadata for the value")

	cmd.MarkFlagRequired("value")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		id := klev.OffsetID(args[0])
		out, err := klient.Offsets.SetRaw(cmd.Context(), id, in)
		return output(out, err)
	}

	return cmd
}

func offsetsDelete() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <offset-id>",
		Short: "delete an offset",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := klev.OffsetID(args[0])
			out, err := klient.Offsets.Delete(cmd.Context(), id)
			return output(out, err)
		},
	}
}
