/*
 * @Author: Hugo
 * @Date: 2022-05-11 11:47:43
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-23 02:33:42
 */
package constant

import "gitlab.com/bns-engineering/td/common/config"

// Domain Names

var (
	GetTDAccountUrl        = "/api/deposits/%v?detailsLevel=FULL"
	SearchTDAccountListUrl = "/api/deposits:search?detailsLevel=FULL&offset=0&limit=500"

	UndoMaturityDateUrl  = "/api/deposits/%v:undoMaturity"
	StartMaturityDateUrl = "/api/deposits/%v:startMaturity"

	ApplyProfitUrl     = "/api/deposits/%v:applyInterest"
	UpdateTDAccountUrl = "/api/deposits/%v"

	CloseAccountUrl = "/api/deposits/%v:changeState"
)

var (
	SearchTransactionUrl    = "/api/deposits/transactions:search?paginationDetails=OFF&offset=0&limit=1"
	WithdrawTransactiontUrl = "/api/deposits/%v/withdrawal-transactions"
	DepositTransactiontUrl  = "/api/deposits/%v/deposit-transactions"
)

var (
	HolidayInfoUrl = "/api/organization/holidays"
)

func getDomainName() string {
	return config.TDConf.Mambu.Host
}

func UrlOf(path string) string {
	return getDomainName() + path
}
