package main

import (
	"github.com/spf13/cobra"

	"github.com/klev-dev/klev-api-go"
)

func filtersRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "filters",
		Short: "interact with filters",
	}

	cmd.AddCommand(filtersList())
	cmd.AddCommand(filtersCreate())
	cmd.AddCommand(filtersGet())
	cmd.AddCommand(filtersStatus())
	cmd.AddCommand(filtersUpdate())
	cmd.AddCommand(filtersDelete())

	return cmd
}

func filtersList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list filters",
		Args:  cobra.NoArgs,
	}

	metadata := cmd.Flags().String("metadata", "", "webhook metadata")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed("metadata") {
			out, err := klient.Filters.Find(cmd.Context(), *metadata)
			return output(out, err)
		} else {
			out, err := klient.Filters.List(cmd.Context())
			return output(out, err)
		}
	}

	return cmd
}

func filtersCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create new filter",
		Args:  cobra.NoArgs,
	}

	var in klev.FilterCreateParams

	sourceID := cmd.Flags().String("source-id", "", "source log id of the filter")
	targetID := cmd.Flags().String("target-id", "", "target log id of the filter")
	cmd.Flags().StringVar(&in.Metadata, "metadata", "", "machine readable metadata")
	cmd.Flags().StringVar(&in.Expression, "expression", "", "expression to eval")

	cmd.MarkFlagRequired("source-id")
	cmd.MarkFlagRequired("target-id")
	cmd.MarkFlagRequired("expression")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		var err error
		in.SourceID, err = klev.ParseLogID(*sourceID)
		if err != nil {
			return outputErr(err)
		}
		in.TargetID, err = klev.ParseLogID(*targetID)
		if err != nil {
			return outputErr(err)
		}

		out, err := klient.Filters.Create(cmd.Context(), in)
		return output(out, err)
	}

	return cmd
}

func filtersGet() *cobra.Command {
	return &cobra.Command{
		Use:   "get <filter-id>",
		Short: "get a filter",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := klev.ParseFilterID(args[0])
			if err != nil {
				return outputErr(err)
			}

			out, err := klient.Filters.Get(cmd.Context(), id)
			return output(out, err)
		},
	}
}

func filtersStatus() *cobra.Command {
	return &cobra.Command{
		Use:   "status <filter-id>",
		Short: "status of a filter",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := klev.ParseFilterID(args[0])
			if err != nil {
				return outputErr(err)
			}

			out, err := klient.Filters.Status(cmd.Context(), id)
			return output(out, err)
		},
	}
}

func filtersUpdate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <filter-id>",
		Short: "update filter",
		Args:  cobra.ExactArgs(1),
	}

	metadata := cmd.Flags().String("metadata", "", "machine readable metadata")
	expression := cmd.Flags().String("expression", "", "expression to eval")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		id, err := klev.ParseFilterID(args[0])
		if err != nil {
			return outputErr(err)
		}

		var in klev.FilterUpdateParams

		if cmd.Flags().Changed("metadata") {
			in.Metadata = metadata
		}
		if cmd.Flags().Changed("expression") {
			in.Expression = expression
		}

		out, err := klient.Filters.UpdateRaw(cmd.Context(), id, in)
		return output(out, err)
	}

	return cmd
}
func filtersDelete() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <filter-id>",
		Short: "delete a filter",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := klev.ParseFilterID(args[0])
			if err != nil {
				return outputErr(err)
			}

			out, err := klient.Filters.Delete(cmd.Context(), id)
			return output(out, err)
		},
	}
}
