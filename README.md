# YNAB Recurring Transactions

CLI that searches through your YNAB transactions (using the [YNAB API](https://api.youneedabudget.com/)) and lists the transactions that look like they recur from month to month. This is useful for when you want to find automatic transactions that occur on a particular account so that you can change them.

## Usage

> Note: This CLI requires a [YNAB Personal Access Token](https://api.youneedabudget.com/#personal-access-tokens).

```
export YNAB_ACCESS_TOKEN="<access-token>"
ynabrt list 
```