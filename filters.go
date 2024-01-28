package main

import (
	"github.com/spf13/cobra"

	"github.com/klev-dev/klev-api-go/filters"
	"github.com/klev-dev/klev-api-go/logs"
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
	cmd.AddCommand(filtersDelete())

	return cmd
}

func filtersList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list filters",
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
	}

	var in filters.CreateParams

	sourceID := cmd.Flags().String("source-id", "", "source log id of the filter")
	targetID := cmd.Flags().String("target-id", "", "target log id of the filter")
	cmd.Flags().StringVar(&in.Metadata, "metadata", "", "machine readable metadata")
	cmd.Flags().StringVar(&in.Expression, "expression", "", "expression to eval")

	cmd.MarkFlagRequired("source-id")
	cmd.MarkFlagRequired("target-id")
	cmd.MarkFlagRequired("expression")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		in.SourceID = logs.LogID(*sourceID)
		in.TargetID = logs.LogID(*targetID)

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
			out, err := klient.Filters.Get(cmd.Context(), filters.FilterID(args[0]))
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
			out, err := klient.Filters.Status(cmd.Context(), filters.FilterID(args[0]))
			return output(out, err)
		},
	}
}

func filtersDelete() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <filter-id>",
		Short: "delete a filter",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := klient.Filters.Delete(cmd.Context(), filters.FilterID(args[0]))
			return output(out, err)
		},
	}
}
