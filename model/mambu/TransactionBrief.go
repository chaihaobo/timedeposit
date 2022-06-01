/*
 * @Author: Hugo
 * @Date: 2022-05-19 09:44:22
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-19 11:11:58
 */
package mambu

import "time"

type TransactionBrief struct {
	EncodedKey         string                             `json:"encodedKey"`
	ID                 string                             `json:"id"`
	CreationDate       time.Time                          `json:"creationDate"`
	ValueDate          time.Time                          `json:"valueDate"`
	Notes              string                             `json:"notes"`
	ParentAccountKey   string                             `json:"parentAccountKey"`
	Type               string                             `json:"type"`
	Amount             float64                            `json:"amount"`
	CurrencyCode       string                             `json:"currencyCode"`
	AccountBalances    TransactionBriefAccountBalances    `json:"accountBalances"`
	UserKey            string                             `json:"userKey"`
	BranchKey          string                             `json:"branchKey"`
	TransactionDetails TransactionBriefTransactionDetails `json:"transactionDetails"`
}

type TransactionBriefAccountBalances struct {
	TotalBalance float64 `json:"totalBalance"`
}

type TransactionBriefTransactionDetails struct {
}
