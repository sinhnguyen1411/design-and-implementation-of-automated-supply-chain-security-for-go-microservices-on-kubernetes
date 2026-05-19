package report

type LineItem struct {
	Date        string  `json:"date"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
}

type Statement struct {
	AccountID      string     `json:"account_id"`
	Items          []LineItem `json:"items"`
	OpeningBalance float64    `json:"opening_balance"`
	ClosingBalance float64    `json:"closing_balance"`
}

func BuildStatement(accountID string, items []LineItem, opening float64) Statement {
	closing := opening
	for _, item := range items { closing += item.Amount }
	return Statement{AccountID: accountID, Items: items, OpeningBalance: opening, ClosingBalance: closing}
}
