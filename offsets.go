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
	cmd.AddCommand(offsetsUpdate())
	cmd.AddCommand(offsetsDelete())

	return cmd
}

func offsetsList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list offsets",
		Args:  cobra.NoArgs,
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
		Args:  cobra.NoArgs,
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
			id, err := klev.ParseOffsetID(args[0])
			if err != nil {
				return err
			}

			out, err := klient.Offsets.Get(cmd.Context(), id)
			return output(out, err)
		},
	}
}

func offsetsUpdate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <offset-id>",
		Short: "update log offset",
		Args:  cobra.ExactArgs(1),
	}

	metadata := cmd.Flags().String("metadata", "", "machine readable metadata")
	value := cmd.Flags().Int64("value", 0, "value to set")
	valueMetadata := cmd.Flags().String("value-metadata", "", "machine readable metadata for the value")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		id, err := klev.ParseOffsetID(args[0])
		if err != nil {
			return err
		}

		var in klev.OffsetUpdateParams

		if cmd.Flags().Changed("metadata") {
			in.Metadata = metadata
		}
		if cmd.Flags().Changed("value") {
			in.Value = value
		}
		if cmd.Flags().Changed("value-metadata") {
			in.ValueMetadata = valueMetadata
		}

		out, err := klient.Offsets.UpdateRaw(cmd.Context(), id, in)
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
			id, err := klev.ParseOffsetID(args[0])
			if err != nil {
				return err
			}

			out, err := klient.Offsets.Delete(cmd.Context(), id)
			return output(out, err)
		},
	}
}
