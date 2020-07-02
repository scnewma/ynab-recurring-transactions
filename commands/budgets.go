package commands

import (
	"io"

	"github.com/spf13/cobra"
	"go.bmvs.io/ynab/api/budget"
)

func ListBudgets() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "budgets",
		Short: "List budgets in your YNAB account.",
		RunE: logError(func(cmd *cobra.Command, args []string) error {
			ynabClient, err := getAPIClient(cmd)
			if err != nil {
				return err
			}

			summaries, err := ynabClient.Budget().GetBudgets()
			if err != nil {
				return err
			}

			return printSummaries(cmd.OutOrStdout(), summaries)
		}),
	}
	return cmd
}

func printSummaries(out io.Writer, summaries []*budget.Summary) error {
	w := newTablewriter(out)
	w.SetHeader([]string{"ID", "NAME"})
	for _, s := range summaries {
		w.Append([]string{s.ID, s.Name})
	}
	w.Render()
	return nil
}
