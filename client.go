package ynabrt

import (
	"fmt"
	"time"

	"go.bmvs.io/ynab"
	"go.bmvs.io/ynab/api/account"
	"go.bmvs.io/ynab/api/payee"
	"go.bmvs.io/ynab/api/transaction"
)

type RecurringTransaction struct {
	Account     string
	Payee       string
	Occurrences []RecurringTransactionOccurrence
}

type RecurringTransactionOccurrence struct {
	Date   time.Time
	Amount int64 // millunits
}

type TransactionFetcher interface {
	GetTransactions(budgetID string, filter *transaction.Filter) ([]*transaction.Transaction, error)
}

type AccountFetcher interface {
	GetAccounts(budgetID string) ([]*account.Account, error)
}

type PayeeFetcher interface {
	GetPayees(budgetID string) ([]*payee.Payee, error)
}

type Client struct {
	transactionFetcher TransactionFetcher
	accountFetcher     AccountFetcher
	payeeFetcher       PayeeFetcher
}

func NewClient(svc ynab.ClientServicer) (*Client, error) {
	if svc == nil {
		return nil, fmt.Errorf("ynab client services must be provided")
	}

	return &Client{
		transactionFetcher: svc.Transaction(),
		accountFetcher:     svc.Account(),
		payeeFetcher:       svc.Payee(),
	}, nil
}
