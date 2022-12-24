package main

import (
	"github.com/spf13/cobra"

	api "github.com/klev-dev/klev-api-go"
)

func offsets() *cobra.Command {
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

	logID := cmd.Flags().String("log-id", "", "log id for this offset")
	metadata := cmd.Flags().String("metadata", "", "offset metadata")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed("log-id") || cmd.Flags().Changed("metadata") {
			out, err := klient.OffsetsFind(cmd.Context(), api.LogID(*logID), *metadata)
			return output(out, err)
		} else {
			out, err := klient.OffsetsList(cmd.Context())
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

	var in api.OffsetCreate

	logID := cmd.Flags().String("log-id", "", "log id for this offset")
	cmd.Flags().StringVar(&in.Metadata, "metadata", "", "machine readable metadata")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		in.LogID = api.LogID(*logID)
		out, err := klient.OffsetCreate(cmd.Context(), in)
		return output(out, err)
	}

	return cmd
}

func offsetsGet() *cobra.Command {
	return &cobra.Command{
		Use:   "get offset_id",
		Short: "get an offset",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := klient.OffsetGet(cmd.Context(), api.OffsetID(args[0]))
			return output(out, err)
		},
	}
}

func offsetsSet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set offset_id",
		Short: "set log offset",
		Args:  cobra.ExactArgs(1),
	}

	value := cmd.Flags().Int64("value", 0, "value to set")
	metadata := cmd.Flags().String("value-metadata", "", "machine readable metadata for the value")

	cmd.MarkFlagRequired("value")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		out, err := klient.OffsetSet(cmd.Context(), api.OffsetID(args[0]), *value, *metadata)
		return output(out, err)
	}

	return cmd
}

func offsetsDelete() *cobra.Command {
	return &cobra.Command{
		Use:   "delete offset_id",
		Short: "delete an offset",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := klient.OffsetDelete(cmd.Context(), api.OffsetID(args[0]))
			return output(out, err)
		},
	}
}
