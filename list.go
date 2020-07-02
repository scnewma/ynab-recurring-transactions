package ynabrt

import (
	"context"
	"fmt"
	"sort"
	"time"

	"go.bmvs.io/ynab/api"
	"go.bmvs.io/ynab/api/transaction"
)

type ListOptions struct {
	BudgetID           string
	TransactionsWithin time.Duration
	MinRecurrence      time.Duration
	MaxRecurrence      time.Duration
}

func (c *Client) List(ctx context.Context, opts ListOptions) ([]*RecurringTransaction, error) {
	filter := &transaction.Filter{}
	if opts.TransactionsWithin > 0 {
		filter.Since = &api.Date{Time: time.Now().Add(-opts.TransactionsWithin)}
	}
	transactions, err := c.transactionFetcher.GetTransactions(opts.BudgetID, filter)
	if err != nil {
		return nil, fmt.Errorf("could not fetch transactions: %w", err)
	}
	accounts, err := c.accountFetcher.GetAccounts(opts.BudgetID)
	if err != nil {
		return nil, fmt.Errorf("could not fetch accounts: %w", err)
	}
	payees, err := c.payeeFetcher.GetPayees(opts.BudgetID)
	if err != nil {
		return nil, fmt.Errorf("could not fetch payees: %w", err)
	}

	getAccountName := func(accountID string) string {
		for _, a := range accounts {
			if a.ID == accountID {
				return a.Name
			}
		}
		return "unknown"
	}

	getPayeeName := func(payeeID string) string {
		for _, p := range payees {
			if p.ID == payeeID {
				return p.Name
			}
		}
		return "unknown"
	}

	groupedTransactions := groupTransactions(transactions)
	return determineRecurring(groupedTransactions, getAccountName, getPayeeName, opts.MinRecurrence, opts.MaxRecurrence), nil
}

func determineRecurring(
	grouped map[string][]*transaction.Transaction,
	getAccountName, getPayeeName func(string) string,
	minRecurrence, maxRecurrence time.Duration) []*RecurringTransaction {
	rts := []*RecurringTransaction{}
	for _, ts := range grouped {
		lastDate := ts[0].Date.Time

		isRecurring := make([]bool, len(ts))
		for i, t := range ts {
			if i == 0 {
				continue
			}

			dateDiff := t.Date.Time.Sub(lastDate)
			if dateDiff >= minRecurrence && dateDiff <= maxRecurrence {
				isRecurring[i] = true
				isRecurring[i-1] = true
			}

			lastDate = t.Date.Time
		}

		rt := &RecurringTransaction{
			Account: getAccountName(ts[0].AccountID),
			Payee:   getPayeeName(getPayeeID(ts[0])),
		}
		recurred := false
		for i, t := range ts {
			if isRecurring[i] {
				recurred = true
				rt.Occurrences = append(rt.Occurrences, RecurringTransactionOccurrence{
					Date:   t.Date.Time,
					Amount: t.Amount,
				})
			}
		}
		if recurred {
			rts = append(rts, rt)
		}
	}
	return rts
}

func groupTransactions(ts []*transaction.Transaction) map[string][]*transaction.Transaction {
	grouped := map[string][]*transaction.Transaction{}

	for _, t := range ts {
		key := fmt.Sprintf("%s-%s", t.AccountID, getPayeeID(t))
		grouped[key] = append(grouped[key], t)
	}

	// sort the groups
	for _, ts := range grouped {
		sort.Slice(ts, func(i, j int) bool {
			return ts[i].Date.Before(ts[j].Date.Time)
		})
	}

	return grouped
}
