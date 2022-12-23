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
	cmd.AddCommand(offsetsAck())
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
		if md := *metadata; md != "" {
			out, err := klient.OffsetsFind(cmd.Context(), md)
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

	cmd.Flags().StringVar(&in.Metadata, "metadata", "", "machine readable metadata")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
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
			out, err := klient.OffsetGetAll(cmd.Context(), api.OffsetID(args[0]))
			return output(out, err)
		},
	}
}

func offsetsAck() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ack offset_id log_id",
		Short: "ack log offset",
		Args:  cobra.ExactArgs(2),
	}

	offset := cmd.Flags().Int64("offset", 0, "offset to set for the log")
	metadata := cmd.Flags().String("metadata", "", "machine readable metadata")

	cmd.MarkFlagRequired("offset")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		out, err := klient.OffsetAck(cmd.Context(), api.OffsetID(args[0]), api.LogID(args[1]), *offset, *metadata)
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
