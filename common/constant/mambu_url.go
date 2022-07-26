// Package constant /*
package constant

import "gitlab.com/bns-engineering/td/common"

var (
	Accept = "application/vnd.mambu.v2+json"
)

const (
	TransactionSucceed = 1
	TransactionFailed  = 0
)

const (
	TransactionWithdraw = "WITHDRAWAL"
	TransactionDeposit  = "DEPOSIT"
)

var (
	GetTDAccountUrl        = "/api/deposits/%v?detailsLevel=FULL"
	SearchTDAccountListUrl = "/api/deposits:search?detailsLevel=FULL&offset=%d&limit=%d&paginationDetails=ON"

	UndoMaturityDateUrl  = "/api/deposits/%v:undoMaturity"
	StartMaturityDateUrl = "/api/deposits/%v:startMaturity"

	ApplyProfitUrl     = "/api/deposits/%v:applyInterest"
	UpdateTDAccountUrl = "/api/deposits/%v"

	CloseAccountUrl = "/api/deposits/%v:changeState"
)

var (
	SearchTransactionUrl   = "/api/deposits/transactions:search?paginationDetails=OFF&offset=0&limit=1"
	WithdrawTransactionUrl = "/api/deposits/%v/withdrawal-transactions"
	DepositTransactionUrl  = "/api/deposits/%v/deposit-transactions"
	AdjustTransactionUrl   = "/api/deposits/transactions/%s:adjust"
)

var (
	HolidayInfoUrl = "/api/organization/holidays"
)

func getDomainName() string {
	return common.C.Mambu.Host
}

func UrlOf(path string) string {
	return getDomainName() + path
}
