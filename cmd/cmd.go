package cmd

import (
	"github.com/scnewma/ynabrt/commands"
	"github.com/spf13/cobra"
)

func NewYNABRTCommand() (*cobra.Command, error) {
	rootCmd := &cobra.Command{
		Use:   "ynabrt",
		Short: "CLI for determining recurring transactions in your YNAB budget.",
	}

	rootCmd.PersistentFlags().String("access-token", "", "YNAB personal access token.")

	rootCmd.AddCommand(commands.ListBudgets())
	rootCmd.AddCommand(commands.ListReccuringTransactions())

	return rootCmd, nil
}
