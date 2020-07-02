package commands

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/scnewma/ynabrt"
	"github.com/spf13/cobra"
)

func ListReccuringTransactions() *cobra.Command {
	var withinDays int
	var minRecurrenceDays int
	var maxRecurrenceDays int
	var budgetID string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List recurring transactions.",
		RunE: logError(func(cmd *cobra.Command, args []string) error {
			client, err := getClient(cmd)
			if err != nil {
				return err
			}

			// TODO: list all budgets
			if budgetID == "" {
				return fmt.Errorf("budget id must be provided")
			}

			dayDuration := func(n int) time.Duration { return time.Duration(n*24) * time.Hour }

			recurredTransactions, err := client.List(context.Background(), ynabrt.ListOptions{
				BudgetID:           budgetID,
				TransactionsWithin: dayDuration(withinDays),
				MinRecurrence:      dayDuration(minRecurrenceDays),
				MaxRecurrence:      dayDuration(maxRecurrenceDays),
			})
			if err != nil {
				return fmt.Errorf("could not list recurred transactions: %v", err)
			}

			return printRecurredTransactions(cmd.OutOrStdout(), recurredTransactions)
		}),
	}

	cmd.Flags().StringVar(&budgetID, "budget", "", "Only show recurring transactions for the budget with this id.")
	cmd.Flags().IntVar(&withinDays, "within", 366, "Only show transactions that occurred within the last X days.")
	cmd.Flags().IntVar(&minRecurrenceDays, "min-recurrence", 25, "Minimum amount of time (in days) between two transactions for the transactions to be considered recurring.")
	cmd.Flags().IntVar(&maxRecurrenceDays, "max-recurrence", 35, "Maximum amount of time (in days) between two transactions for the transactions to be considered recurring.")

	return cmd
}

func printRecurredTransactions(out io.Writer, ts []*ynabrt.RecurringTransaction) error {
	w := newTablewriter(out)
	w.SetHeader([]string{"ACCOUNT", "PAYEE", "DATE", "TOTAL"})
	for _, t := range ts {
		for _, o := range t.Occurrences {
			w.Append([]string{t.Account, t.Payee, o.Date.Format("2006-01-02"), currency(o.Amount)})
		}
	}
	w.Render()
	return nil
}

func currency(munits int64) string {
	return fmt.Sprintf("$%.2f", (float64(munits)/1000)*-1)
}
