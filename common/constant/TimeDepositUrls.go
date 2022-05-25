/*
 * @Author: Hugo
 * @Date: 2022-05-11 11:47:43
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-23 02:33:42
 */
package constant

// Domain Names
const (
	DomainName = "https://cbs-dev1.aladinbank.id"
)

const (
	GetTDAccountUrl        = DomainName + "/api/deposits/%v?detailsLevel=FULL"
	SearchTDAccountListUrl = DomainName + "/api/deposits:search?detailsLevel=FULL&offset=0&limit=500"

	UndoMaturityDateUrl  = DomainName + "/api/deposits/%v:undoMaturity"
	StartMaturityDateUrl = DomainName + "/api/deposits/%v:startMaturity"

	ApplyProfitUrl     = DomainName + "/api/deposits/%v:applyInterest"
	UpdateTDAccountUrl = DomainName + "/api/deposits/%v"

	CloseAccountUrl = DomainName + "/api/deposits/%v:changeState"
)

const (
	SearchTransactionUrl    = DomainName + "/api/deposits/transactions:search?paginationDetails=OFF&offset=0&limit=1"
	WithdrawTransactiontUrl = DomainName + "/api/deposits/%v/withdrawal-transactions"
	DepositTransactiontUrl  = DomainName + "/api/deposits/%v/deposit-transactions"
)

const (
	HolidayInfoUrl = DomainName + "/api/organization/holidays"
)
