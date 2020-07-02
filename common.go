package ynabrt

import "go.bmvs.io/ynab/api/transaction"

func getPayeeID(t *transaction.Transaction) string {
	if t.PayeeID == nil {
		return "unknown"
	}
	return *t.PayeeID
}
