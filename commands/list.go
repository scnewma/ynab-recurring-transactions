package commands

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/scnewma/ynabrt"
	"github.com/spf13/cobra"
	"go.bmvs.io/ynab/api/budget"
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
			apiClient, err := getAPIClient(cmd)
			if err != nil {
				return err
			}
			client, err := getClient(cmd)
			if err != nil {
				return err
			}

			dayDuration := func(n int) time.Duration { return time.Duration(n*24) * time.Hour }

			budgetSummaries, err := apiClient.Budget().GetBudgets()
			if err != nil {
				return fmt.Errorf("could not list budgets: %v", err)
			}

			if budgetID != "" {
				var budget *budget.Summary
				for _, s := range budgetSummaries {
					if s.ID == budgetID {
						budget = s
					}
				}

				if budget == nil {
					return fmt.Errorf("budget with id %q not found", budgetID)
				}

				recurredTransactions, err := client.List(context.Background(), ynabrt.ListOptions{
					BudgetID:           budgetID,
					TransactionsWithin: dayDuration(withinDays),
					MinRecurrence:      dayDuration(minRecurrenceDays),
					MaxRecurrence:      dayDuration(maxRecurrenceDays),
				})
				if err != nil {
					return fmt.Errorf("could not list recurred transactions: %v", err)
				}

				printRecurredTransactions(cmd.OutOrStdout(), recurredTransactions)
				return nil
			}

			for _, budgetSummary := range budgetSummaries {
				fmt.Printf("\n%s:\n\n", strings.ToUpper(budgetSummary.Name))
				recurredTransactions, err := client.List(context.Background(), ynabrt.ListOptions{
					BudgetID:           budgetSummary.ID,
					TransactionsWithin: dayDuration(withinDays),
					MinRecurrence:      dayDuration(minRecurrenceDays),
					MaxRecurrence:      dayDuration(maxRecurrenceDays),
				})
				if err != nil {
					return fmt.Errorf("could not list recurred transactions: %v", err)
				}

				printRecurredTransactions(cmd.OutOrStdout(), recurredTransactions)
			}
			return nil
		}),
	}

	cmd.Flags().StringVar(&budgetID, "budget", "", "Only show recurring transactions for the budget with this id.")
	cmd.Flags().IntVar(&withinDays, "within", 366, "Only show transactions that occurred within the last X days.")
	cmd.Flags().IntVar(&minRecurrenceDays, "min-recurrence", 25, "Minimum amount of time (in days) between two transactions for the transactions to be considered recurring.")
	cmd.Flags().IntVar(&maxRecurrenceDays, "max-recurrence", 35, "Maximum amount of time (in days) between two transactions for the transactions to be considered recurring.")

	return cmd
}

func printRecurredTransactions(out io.Writer, ts []*ynabrt.RecurringTransaction) {
	if len(ts) == 0 {
		fmt.Print("no recurring transactions found\n")
		return
	}

	w := newTablewriter(out)

	w.SetHeader([]string{"ACCOUNT", "PAYEE", "DATE", "TOTAL"})
	for _, t := range ts {
		for _, o := range t.Occurrences {
			w.Append([]string{t.Account, t.Payee, o.Date.Format("2006-01-02"), currency(o.Amount)})
		}
	}
	w.Render()
}

func currency(munits int64) string {
	return fmt.Sprintf("$%.2f", (float64(munits)/1000)*-1)
}
